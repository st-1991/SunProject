package controllers

import (
	"SunProject/config"
	"SunProject/models"
	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {
	res := config.Result{}
	res.Success(c)
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
	//cd := func(ch chan []models.ApiUser) chan []models.ApiUser {
	//	c := make(chan []models.ApiUser)
	//	go func() {
	//		select {
	//		case m := <-ch:
	//			ch <- m
	//		}
	//	}()
	//	return c
	//}(ch)
	//users := models.Users()
	res["users"] = models.Users()
	ApiResponse(c, &Response{Data: res })
}
