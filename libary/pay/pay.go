package pay

import (
	"SunProject/config"
	"SunProject/libary/request"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const Host = "https://www.chyzf.cn/"

const Key = "UN6d7DvSSf0Urm5m6DU63VNe7N3VVvu0"

//type ApiParam struct {
//	Pid int `json:"pid"`
//	Type string `json:"type"`
//	OutTradeNo string `json:"out_trade_no"`
//	Notify string `json:"notify"`
//	Name string `json:"name"`
//	ClientIp string `json:"clientip"`
//}{}

type ApiParam struct {
	Type string `json:"type"`
	Name string `json:"name"`
	OutTradeNo string `json:"out_trade_no"`
	Money string `json:"money"`
	ClientIp string `json:"clientip"`
	Device string `json:"device"`
}

func (p ApiParam) CreateOrder() (map[string]string, error) {
	notifyUrl := "https://api.aizj.top/api/keep/notify"
	params := map[string]interface{}{
		"pid": 1063,
		"type": p.Type,
		"out_trade_no": p.OutTradeNo,
		"notify_url": notifyUrl,
		"return_url": "https://www.aizj.top/result",
		"name": p.Name,
		"money": p.Money,
		"clientip": p.ClientIp,
		"device": "pc",
		"sign_type": "MD5",
	}
	params["sign"] = createSign(params)
	//paramB, _ := json.Marshal(params)
	//resp := request.JsonPost(Host + "mapi.php", paramB, nil)
	resp := request.FromUrlencodedPost(Host + "mapi.php", params)
	if resp.Err != nil {
		return nil, resp.Err
	}
	bodyB, _ := resp.Body()
	data := struct {
		Code int `json:"code"`
		Msg string `json:"msg"`
		TradeNo string `json:"trade_no"`
		QrCode string `json:"qrcode"`
		PayUrl string `json:"payurl"`
	}{}
	_ = json.Unmarshal(bodyB, &data)
	if data.Code != 1 {
		return nil, fmt.Errorf(data.Msg)
	}
	return map[string]string{
		"payurl": data.PayUrl,
		"qrcode": data.QrCode,
		"order_sn": p.OutTradeNo,
	}, nil
}

func createSign(params map[string]interface{}) string {
	delete(params, "sign")
	delete(params, "sign_type")
	var keys []string
	for key := range params {
		switch params[key].(type) {
		case string:
			if len(params[key].(string)) > 0 { // 空值不参加拼接串中
				keys = append(keys, key)
			}
		case int:
			if params[key].(int) > 0 {
				keys = append(keys, key)
			}
		}
	}
	sort.Strings(keys)
	var pList []string
	var v string
	for _, k := range keys {
		if k == "pid" {
			v = strconv.Itoa(params[k].(int))
		} else {
			v = params[k].(string)
		}
		pList = append(pList, k + "=" + v)
	}
	sortedParamsStr := strings.Join(pList, "&")
	config.Logger().Info("md5 前的值:", sortedParamsStr)
	hasher := md5.New()
	hasher.Write([]byte(sortedParamsStr + Key))
	signBytes := hasher.Sum(nil)
	sign := hex.EncodeToString(signBytes)
	return sign
}
