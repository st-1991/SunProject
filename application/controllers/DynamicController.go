package controllers

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"

	"SunProject/libary/uniapp"
)

type recommendDynamic struct {
	Avatar string `json:"avatar"`
	Nickname string `json:"nickname"`
	IsThumbUp bool `json:"is_thumb_up"`
	Images []string `json:"images"`
}

type AddDynamicParam struct {
	Text string `form:"text"`
	Tag string `form:"tag"`
	GoodId int `form:"good_id"`
	Images []*multipart.FileHeader `form:"images"`
}

func AddDynamic(c *gin.Context) {
	if !IsLogin(c) {
		return
	}
	//userId := c.GetInt("userId")
	var dParam AddDynamicParam
	if c.Bind(&dParam) != nil {
		ApiError(c, &Response{Code: -1, Msg: "参数绑定错误"})
		return
	}
	if dParam.Text == "" {
		ApiError(c, &Response{Code: -1, Msg: "请填写动态内容哦"})
		return
	}
	files := dParam.Images
	images := make([]string, 0)
	if len(files) > 0 {
		fileUrls := uniapp.UploadFiles(files)
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
	}
	ApiResponse(c, &Response{Code: 0, Msg: "发布动态成功～"})
}