package controllers

import (
	"SunProject/application/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetIntegralLogs(c *gin.Context) {
	if !IsLogin(c) {
		return
	}
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	userId := c.MustGet("userId").(int)
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	logs := models.GetIntegralLogsByUserId(userId, pageInt, pageSizeInt)
	ApiResponse(c, &Response{Data: logs})
}