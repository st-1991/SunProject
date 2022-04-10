package uniapp

import (
	"SunProject/config"
	"SunProject/libary/request"
	"encoding/json"
	"fmt"
	"mime/multipart"
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
}

type AliYun struct {
	ID string `json:"id"`
	CdnDomain string `json:"cdnDomain"`
	Signature string `json:"signature"`
	Policy string `json:"policy"`
	AccessKeyId string `json:"accessKeyId"`
	OssPath string `json:"ossPath"`
	Host string `json:"host"`
}

func (u *UniApp) InitConfig() *UniApp {
	u.SpaceId = "6dd07f54-c3ab-447b-87ff-51ff48538c90"
	u.ClientSecret = "fWQ7szy9J4u0gKz0ruZoGA=="
	u.Url = "https://api.bspapp.com/client"
	return u
}

func (u UniApp) getSign(params map[string]string) string {
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

func (u UniApp) getToken() (string, error) {
	params := map[string]string{
		"method":    "serverless.auth.user.anonymousAuthorize",
		"params":    "{}",
		"spaceId":   u.SpaceId,
		"timestamp": strconv.FormatInt(time.Now().UnixNano() / 1e6, 10),
	}
	headers := CommonHeaders()
	headers["x-serverless-sign"] = u.getSign(params)
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

func (u *UniApp) preUploadFile(filename string) (AliYun, error) {
	paramsStruct := struct {
		FileName string `json:"filename"`
		Env string `json:"env"`
	}{}
	paramsStruct.FileName = filename
	paramsStruct.Env = "public"
	params, _ := json.Marshal(paramsStruct)

	token, err := u.getToken()
	if err != nil {
		return AliYun{}, err
	}
	requestParams := map[string]string{
		"method":    "serverless.file.resource.generateProximalSign",
		"params": string(params),
		"spaceId":   u.SpaceId,
		"timestamp": strconv.FormatInt(time.Now().UnixNano() / 1e6, 10),
		"token" : token,
	}
	headers := CommonHeaders()
	headers["x-serverless-sign"] = u.getSign(requestParams)
	headers["x-basement-token"] = token
	requestParamsJson, _ := json.Marshal(requestParams)
	result := request.JsonPost(u.Url, requestParamsJson, headers)
	data, err := result.Body()
	if err != nil {
		return AliYun{}, err
	}
	dataJson := struct {
		Success bool `json:"success"`
		Data AliYun `json:"data"`
	}{}
	_ = json.Unmarshal(data, &dataJson)
	if 	result.StatusCode() != http.StatusOK || dataJson.Success != true {
		return AliYun{}, fmt.Errorf("请求失败")
	}
	u.Token = token
	return dataJson.Data, nil
}

func (a AliYun) upload(file *multipart.FileHeader) (string, error){
	params := map[string][]byte{
		"Cache-Control": []byte("max-age=2592000"),
		"Content-Disposition": []byte("attachment"),
		"OSSAccessKeyId": []byte(a.AccessKeyId),
		"Signature": []byte(a.Signature),
		"host": []byte(a.Host),
		"id": []byte(a.ID),
		"key": []byte(a.OssPath),
		"policy": []byte(a.Policy),
		"success_action_status": []byte("200"),
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	fileParam := []request.FileParam{
		{
			FiledName: "file",
			Name: file.Filename,
			FormFile: src,
		},
	}
	res := request.FromDataPost("https://" + a.Host + "/", params, fileParam, map[string]string{"X-OSS-server-side-encrpytion":"AES256"})
	if res.Err != nil {
		return "", res.Err
	}
	data, err := res.Body()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (u *UniApp) CompleteUploadFile(file *multipart.FileHeader) (string, error) {
	// 图片预加载
	AliYun, err := u.preUploadFile(file.Filename)
	if err != nil {
		config.Logger().Error(fmt.Sprintf("图片上传，预加载失败：%s", err))
		return "", err
	}
	// 图片上传
	_, err = AliYun.upload(file)
	if err != nil {
		config.Logger().Error(fmt.Sprintf("图片上传失败：%s", err))
		return "", err
	}

	paramsStruct := struct {
		ID string `json:"id"`
	}{}
	paramsStruct.ID = AliYun.ID
	params, _ := json.Marshal(paramsStruct)

	// token不存在重新加载
	if len(u.Token) == 0 {
		if token, err := u.getToken(); err == nil {
			u.Token = token
		} else {
			return "", err
		}
	}

	requestParams := map[string]string{
		"method":    "serverless.file.resource.report",
		"params": string(params),
		"spaceId":   u.SpaceId,
		"timestamp": strconv.FormatInt(time.Now().UnixNano() / 1e6, 10),
		"token" : u.Token,
	}
	headers := CommonHeaders()
	headers["x-serverless-sign"] = u.getSign(requestParams)
	headers["x-basement-token"] = u.Token
	requestParamsJson, _ := json.Marshal(requestParams)
	result := request.JsonPost(u.Url, requestParamsJson, headers)
	data, err := result.Body()
	if err != nil {
		return "", err
	}
	dataJson := struct {
		Success bool `json:"success"`
	}{}
	_ = json.Unmarshal(data, &dataJson)
	if 	result.StatusCode() != http.StatusOK || dataJson.Success != true {
		return "", err
	}
	return "https://" + AliYun.Host + "/" + AliYun.OssPath, nil
}

func UploadWorker(file *multipart.FileHeader, c chan string)  {
	go func(file *multipart.FileHeader, c chan string) {
		uniApp := UniApp{}
		url, err := uniApp.InitConfig().CompleteUploadFile(file)
		if err != nil {
			c <- fmt.Sprintf("预加载失败：%s", err)
		}
		c <- url
		return
	}(file, c)
}