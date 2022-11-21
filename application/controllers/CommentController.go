package controllers

import (
	"SunProject/application/models"
	"SunProject/application/service"
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

// CommentsParam 获取评论列表参数
type CommentsParam struct {
	Id int `form:"id"`
	Level int `form:"level"`
	Page int `form:"page"`
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
	// 评论数加1
	go service.AddCommentNum(cModel.DynamicId, cModel.ParentCommentId, level)
	ApiResponse(c, &Response{})
}


type CommentInfo struct {
	models.Comments
	Avatar string `json:"avatar"`
	Nickname string `json:"nickname"`
	ParentAvatar string `json:"parent_avatar"`
	ParentNickname string `json:"parent_nickname"`
}

func Comments(c *gin.Context) {
	var CParam CommentsParam
	if c.BindQuery(&CParam) != nil {
		ApiError(c, &Response{Code: -1, Msg: "参数绑定错误"})
		return
	}
	// 获取评论列表
	comments := models.GetCommentsById(CParam.Id, CParam.Level, CParam.Page, 10)
	var userIds []int
	usersMap := make(map[int]models.UserBase)
	// 获取用户id
	for _, comment := range comments {
		if _, ok := usersMap[comment.UserId]; !ok {
			userIds = append(userIds, comment.UserId)
			usersMap[comment.UserId] = models.UserBase{}
		}
		if _, ok := usersMap[comment.ParentCommentUserId]; !ok {
			userIds = append(userIds, comment.ParentCommentUserId)
			usersMap[comment.ParentCommentUserId] = models.UserBase{}
		}
	}
	// 获取用户信息
	u := models.User{}
	users := u.GetUsersByIds(userIds)
	for _, user := range users{
		usersMap[user.Id] = user
	}

	// 数据拼接
	var apiComments []CommentInfo
	for _, comment := range comments {
		cInfo := CommentInfo{Comments: comment}
		cInfo.Avatar = usersMap[comment.UserId].Avatar
		cInfo.Nickname = usersMap[comment.UserId].Nickname
		if parentUser, ok := usersMap[comment.ParentCommentUserId]; ok {
			cInfo.ParentAvatar = parentUser.Avatar
			cInfo.ParentNickname = parentUser.Nickname
		}
		apiComments = append(apiComments, cInfo)
	}
	ApiResponse(c, &Response{Data: apiComments})
}
