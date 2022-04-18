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
