package controllers

import (
	"SunProject/application/models"
	"SunProject/libary/uniapp"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type Comment struct {
	DynamicId int `form:"dynamic_id"`
	CommentId int `form:"comment_id"`
	Text string `form:"text"`
	Images []*multipart.FileHeader `form:"images"`
}

func AddComment(c *gin.Context) {
	if !IsLogin(c) {
		return
	}
	var comment Comment
	if c.Bind(&comment) != nil {
		ApiError(c, &Response{Code: -1, Msg: "参数绑定错误"})
		return
	}
	if comment.DynamicId == 0 && comment.CommentId == 0{
		ApiError(c, &Response{Code: -1, Msg: "参数错误"})
		return
	}
	if comment.Text == "" {
		ApiError(c, &Response{Code: -1, Msg: "请填写评论内容～"})
		return
	}
	if len(comment.Images) > 1 {
		ApiError(c, &Response{Code: -1, Msg: "最多只能添加1张图片～"})
		return
	}

	parentComment := models.CommenterInfo{}
	if comment.CommentId != 0 {
		parentComment = models.GetCommenterInfoById(comment.CommentId)
		if parentComment.Id == 0 {
			ApiError(c, &Response{Code: -1, Msg: "评论不存在～"})
			return
		}
		if parentComment.Level == 1 {
			parentComment.TopCommentId = parentComment.Id
		}
	}

	images := make([]string, 0)
	if len(comment.Images) > 0 {
		fileUrls := uniapp.UploadFiles(comment.Images)
		if fileUrls[0].Err != nil {
			ApiError(c, &Response{Code: -1, Msg: "图片上传失败～"})
			return
		}
		images = append(images, fileUrls[0].FileUrl)
	}
	filesJson, _ := json.Marshal(images)

	var level int
	if comment.CommentId == 0 {
		level = 1
	} else {
		level = 2
	}
	cModel := &models.Comments{
		DynamicId: comment.DynamicId,
		ParentCommentId: comment.CommentId,
		ParentCommentUserId: parentComment.UserId,
		TopCommentId: parentComment.TopCommentId,
		Text: comment.Text,
		Images: string(filesJson),
		UserId: c.GetInt("userId"),
		Level: level,
	}
	if ok := cModel.Create(); !ok {
		ApiError(c, &Response{Code: -1, Msg: "评论失败～"})
		return
	}
	ApiResponse(c, &Response{})
}
