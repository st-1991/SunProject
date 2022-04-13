package controllers

import (
	"SunProject/application/models"
	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {

	if isLogin, ok := c.Get("isLogin"); !ok || isLogin.(bool) == false {
		ApiResponse(c, &Response{Code: 4444, Msg: "请先登录"})
		return
	}
	userId := c.MustGet("userId").(int)

	user, ok := models.GetUser("", userId)
	if !ok {
		ApiResponse(c, &Response{Code: -1, Msg: "用户不存在"})
		return
	}
	ApiResponse(c, &Response{Code: 0, Msg: "success", Data: map[string]models.User{
		"user": user,
	}})
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
	res["users"] = models.Users()
	ApiResponse(c, &Response{Data: res })
}
