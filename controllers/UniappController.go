package controllers

import (
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context)  {
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	filenames := make([]string, len(files))
	for _, file := range files {
		filename := file.Filename
		if filename != "" {
			filenames = append(filenames, file.Filename)
		}
	}
	//UniApp := uniapp.UniApp{}
	//token, err := UniApp.InitConfig().PreUploadFile("a.jpg")
	//if err != nil {
	//	ApiResponse(c, &Response{Code: -1, Msg: "预加载失败"})
	//}
	ApiResponse(c, &Response{Code: 200, Msg: "上传成功", Data: filenames})
}
