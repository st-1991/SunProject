package controllers

import (
	"SunProject/libary/uniapp"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	var fileUrls []string
	ch := make(chan string, len(files))
	for _, file := range files {
		//UniApp := uniapp.UniApp{}
		//url, err := UniApp.InitConfig().CompleteUploadFile(file)
		//if err != nil {
		//	errs = append(errs, fmt.Errorf("预加载失败：%s", err))
		//}
		//fileUrls = append(fileUrls, url)
		uniapp.UploadWorker(file, ch)
	}
	for {
		if len(fileUrls)  == len(files) {
			break
		}
		f := <-ch
		fileUrls = append(fileUrls, f)
	}
	ApiResponse(c, &Response{Code: 200, Msg: "上传成功", Data: map[string]interface{}{
		"file_urls": fileUrls,
	}})
}
