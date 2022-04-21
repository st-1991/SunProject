package controllers

import (
	"SunProject/application/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsLogin(c *gin.Context) bool {
	if isLogin, ok := c.Get("isLogin"); !ok || isLogin.(bool) == false {
		ApiError(c, &Response{Code: 4444, Msg: "请先登录"}, http.StatusUnauthorized)
		return false
	}
	return true
}

// UserInfo 用户详情
func UserInfo(c *gin.Context) {
	if !IsLogin(c) {
		return
	}
	userId := c.MustGet("userId").(int)

	user, ok := models.GetUser("", userId)
	if !ok {
		ApiError(c, &Response{Code: -1, Msg: "用户不存在"})
		return
	}
	ApiResponse(c, &Response{Code: 0, Msg: "success", Data: user})
}

// EditUser 编辑用户信息
func EditUser(c *gin.Context)  {
	if !IsLogin(c) {
		return
	}
	userId := c.MustGet("userId").(int)
	avatar := c.PostForm("avatar")
	nickname := c.PostForm("nickname")
	sex := c.PostForm("sex")
	birthday := c.PostForm("birthday")
	profile := c.PostForm("profile")

	user := models.User{
		Id: userId,
		Avatar: avatar,
		Nickname: nickname,
		Sex: sex,
		Birthday: birthday,
		Profile: profile,
	}
	res := user.EditUser()
	if !res {
		ApiError(c, &Response{Code: -1, Msg: "编辑资料失败，请重试！"})
		return
	}
	ApiResponse(c, &Response{})
}

func UserList(c *gin.Context) {

	//userList := new([]map[string]interface{})
	//config.DB.Table("bm_users").Take(userList)

	//var ch chan []models.ApiUser
	//ch := make(chan []models.ApiUser, 3)
	//go func(ch chan []models.ApiUser) {
	//	 ch <- models.Users()
	//	 ch <- models.Users()
	//	 ch <- models.Users()
	//}(ch)

	//select {
	//case m := <-ch:
	//	ch <- m
	//}
	res := make(map[string]interface{})
	//users := models.Users()
	//res["users"] = models.Users()
	ApiResponse(c, &Response{Data: res })
}
