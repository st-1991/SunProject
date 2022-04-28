package controllers

import (
	"encoding/json"
	"mime/multipart"
	"strconv"

	"github.com/gin-gonic/gin"

	"SunProject/application/models"
	"SunProject/application/service"
	"SunProject/libary/uniapp"
)

type recommendDynamic struct {
	models.ApiDynamic
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
	userId := c.GetInt("userId")
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
	filesJson, _ := json.Marshal(images)
	dynamic := models.Dynamic{Text: dParam.Text, Tag: dParam.Tag, Images: string(filesJson), UserId: userId}
	if !dynamic.CreateDynamic() {
		ApiError(c, &Response{Code: -1, Msg: "发布动态失败，请重试～"})
		return
	}
	ApiResponse(c, &Response{Code: 0, Msg: "发布动态成功～"})
}

// RecommendDynamics 推荐动态
func RecommendDynamics(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	var recommendDynamics []recommendDynamic
	var d models.Dynamic
	var userIds []int

	dynamics := d.GetDynamics(page, pageSize)
	for _, dynamic := range dynamics {
		userIds = append(userIds, dynamic.UserId)
	}
	u := models.User{}
	users := u.GetUsersByIds(userIds)
	usersMap := make(map[int]models.UserBase)
	for _, user := range users{
		usersMap[user.Id] = user
	}
	for _, dynamic := range dynamics {
		var rd recommendDynamic
		ud := service.UserDynamic{
			Id: dynamic.Id,
			UserId: c.GetInt("userId"),
		}
		_ = json.Unmarshal([]byte(dynamic.Images), &rd.Images)
		rd.ApiDynamic = dynamic
		rd.Avatar = usersMap[dynamic.UserId].Avatar
		rd.Nickname = usersMap[dynamic.UserId].Nickname
		rd.IsThumbUp = ud.IsThumbUp()
		recommendDynamics = append(recommendDynamics, rd)
	}
	ApiResponse(c, &Response{Code: 0, Data: recommendDynamics})
}

func DynamicThumbUp(c *gin.Context) {
	if !IsLogin(c) {
		return
	}
	userId := c.GetInt("userId")
	dIdStr := c.Param("dId")
	dId, err := strconv.Atoi(dIdStr)
	if err != nil {
		ApiError(c, &Response{Code: -1, Msg: "参数错误：d_id is not int"})
		return
	}
	go service.DynamicThumbUp(service.UserDynamic{
		UserId: userId,
		Id: dId,
	})
	ApiResponse(c, &Response{Code: 0, Msg: "点赞完成～"})
}