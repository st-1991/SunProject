package uniapp

import (
	"SunProject/config"
	"SunProject/libary/request"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func CommonHeaders() map[string]string {
	return map[string]string{
		"accept": "*/*",
		"accept-encoding": "gzip, deflate, sdch",
		"accept-language": "zh-CN,zh;q=0.8",
		"cache-control": "no-cache",
		"pragma": "no-cache",
		"connection": "close",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	}
}

type UniApp struct {
	SpaceId string
	ClientSecret string
	Token string
	Url string
	Sign string
}

func (u UniApp) InitConfig() UniApp {
	u.SpaceId = "6dd07f54-c3ab-447b-87ff-51ff48538c90"
	u.ClientSecret = "fWQ7szy9J4u0gKz0ruZoGA=="
	u.Url = "https://api.bspapp.com/client"
	return u
}

func (u UniApp) GetSign(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var signStr string
	for _, k := range keys {
		signStr += k + "=" + params[k] + "&"
	}
	signStr = signStr[:len(signStr)-1]
	return config.GenHMACMd5([]byte(signStr), []byte(u.ClientSecret))
}

func (u UniApp) GetToken() (string, error) {
	params := map[string]string{
		"method":    "serverless.auth.user.anonymousAuthorize",
		"params":    "{}",
		"spaceId":   u.SpaceId,
		"timestamp": strconv.FormatInt(time.Now().UnixNano() / 1e6, 10),
	}
	headers := CommonHeaders()
	headers["x-serverless-sign"] = u.GetSign(params)
	paramsJson, _ := json.Marshal(params)
	result := request.JsonPost(u.Url, paramsJson, headers)
	data, err := result.Body()
	if err != nil {
		return "", err
	}
	r := struct {
		Success bool `json:"success"`
		Data struct{
			AccessToken string `json:"accessToken"`
		} `json:"data"`
	}{}
	_ = json.Unmarshal(data, &r)
	if 	result.StatusCode() != http.StatusOK || r.Success != true {
		return "", nil
	}
	return r.Data.AccessToken, nil
}

func (u UniApp) PreUploadFile(filename string) (string, error) {
	paramsStruct := struct {
		FileName string `json:"filename"`
		Env string `json:"env"`
	}{}
	paramsStruct.FileName = filename
	paramsStruct.Env = "public"
	params, _ := json.Marshal(paramsStruct)

	token, err := u.GetToken()
	if err != nil {
		return "", err
	}
	requestParams := map[string]string{
		"method":    "serverless.file.resource.generateProximalSign",
		"params": string(params),
		"spaceId":   u.SpaceId,
		"timestamp": strconv.FormatInt(time.Now().UnixNano() / 1e6, 10),
		"token" : token,
	}
	headers := CommonHeaders()
	headers["x-serverless-sign"] = u.GetSign(requestParams)
	headers["x-basement-token"] = token
	requestParamsJson, _ := json.Marshal(requestParams)
	result := request.JsonPost(u.Url, requestParamsJson, headers)
	data, err := result.Body()
	if err != nil {
		return "", err
	}
	return string(data), nil
}