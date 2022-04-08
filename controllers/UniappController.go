package controllers

import (
	"SunProject/libary/uniapp"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context)  {
	UniApp := uniapp.UniApp{}
	token, err := UniApp.InitConfig().PreUploadFile("a.jpg")
	if err != nil {
		ApiResponse(c, &Response{Code: -1, Msg: "预加载失败"})
	}
	ApiResponse(c, &Response{Code: 200, Msg: "上传成功", Data: token})
}
