package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"

	"SunProject/application/models"
	"SunProject/libary/uniapp"
)

func AddDynamic(c *gin.Context) {
	if !IsLogin(c) {
		return
	}
	userId := c.GetInt("userId")
	text := c.PostForm("text")
	tag := c.PostForm("tag")
	//goodId := c.PostForm("goodId")
	if text == "" {
		ApiError(c, &Response{Code: -1, Msg: "请填写动态内容哦"})
		return
	}
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	fileUrls := uniapp.UploadFiles(files)
	var images []string
	urlMap := make(map[string]string)
	for _, url := range fileUrls {
		if url.Err == nil {
			urlMap[url.FileName] = url.FileUrl
		}
	}

	for _, file := range files {
		if fileUrl, ok := urlMap[file.Filename]; ok {
			images = append(images, fileUrl)
		}
	}
	imagesStr := ""
	if len(images) > 0 {
		filesJson, _ := json.Marshal(images)
		imagesStr = string(filesJson)
	}
	dynamic := models.Dynamic{Text: text, Tag: tag, Images: imagesStr, UserId: userId}
	if !dynamic.CreateDynamic() {
		ApiError(c, &Response{Code: -1, Msg: "发布动态失败，请重试～"})
		return
	}
	ApiResponse(c, &Response{Code: 0, Msg: "发布动态成功～"})
}
