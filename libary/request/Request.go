package request

import (
	"SunProject/config"
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type Result struct {
	Response *http.Response
	Err error
}

type FileParam struct {
	FiledName string
	Name string
	FormFile multipart.File
}

// FromDataPost Content-Type: multipart/form-data
func FromDataPost(url string, params map[string][]byte, files []FileParam, headers map[string]string) Result {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	for k, v := range params {
		if 	fw, err := w.CreateFormField(k); err == nil {
			fw.Write(v)
		}
	}

	for _, file := range files {
		if createFormFile, err := w.CreateFormFile(file.FiledName, file.Name); err == nil {
			readAll, _ := ioutil.ReadAll(file.FormFile)
			createFormFile.Write(readAll)
		}
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		config.Logger().Error(fmt.Sprintf("建立链接失败：%s", err))
		return Result{nil, err}
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		config.Logger().Error(fmt.Sprintf("发送请求失败：%s", err))
		return Result{nil, err}
	}
	config.Logger().Info(fmt.Sprintf("%d", res.StatusCode))
	return Result{res, nil}
}

func JsonPost(url string, params []byte, headers map[string]string) Result {
	reader := bytes.NewReader(params)
	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		config.Logger().Error(fmt.Sprintf("建立链接失败：%s", err))
		return Result{nil, err}
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		config.Logger().Error(fmt.Sprintf("发送请求失败：%s", err))
		return Result{nil, err}
	}
	return Result{res, nil}
}



// www-from-urlencoded
func FromUrlencodedPost() {
	// TODO: 待实现
	//resp, err := http.PostForm(SmsHost, url.Values{
	//	"phone": {phone},
	//	"templateId": {templateId},
	//	"data": {string(smsData)},
	//})
}

func (r Result) StatusCode() int {
	return r.Response.StatusCode
}

func (r Result) Body() ([]byte, error){
	body, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		config.Logger().Error(fmt.Sprintf("获取Body失败：%s", err))
		return nil, err
	}
	defer r.Response.Body.Close()
	config.Logger().Info(fmt.Sprintf("请求结果：%s", body))
	return body, nil
}

func (r Result) Header() http.Header{
	return r.Response.Header
}