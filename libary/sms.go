package libary

import (
	"SunProject/libary/request"
	"encoding/json"
	"net/http"
)

var SmsHost = "https://uni.mvc.cc/functionName/send/sms"

var templateId = "14113"

type SmsData struct {
	Code      string `json:"code"`
	ExpMinute int `json:"expMinute"`
}

// SendSms 发送验证码
func SendSms(phone string, code string) bool {
	smsData, _:= json.Marshal(&SmsData{
		Code: code,
		ExpMinute: 5,
	})

	params := map[string][]byte{
		"phone": []byte(phone),
		"templateId": []byte(templateId),
		"data": smsData,
	}

	res := request.FromDataPost(SmsHost, params, nil, nil)
	if res.Err != nil {
		return false
	}
	data, err := res.Body()
	if err != nil {
		return false
	}
	r := struct {
		Code int `json:"code"`
		Success bool `json:"success"`
		Phone string `json:"phone"`
	}{}
	_ = json.Unmarshal(data, &r)
	if 	res.StatusCode() != http.StatusOK || r.Code != 0 || r.Success != true {
		return false
	}
	return true
}

func RequestPost()  {
	
}