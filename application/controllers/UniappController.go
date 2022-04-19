package controllers

import (
	"SunProject/libary/uniapp"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	fileUrls := uniapp.UploadFiles(files)
	ApiResponse(c, &Response{Code: 200, Msg: "上传成功", Data: map[string]interface{}{
		"file_urls": fileUrls,
	}})
}
