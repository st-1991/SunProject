package controllers

import (
	"SunProject/application/models"
	"github.com/gin-gonic/gin"
)

func Tabs(c *gin.Context) {
	var tab models.Tabs
	tabs := tab.GetTabs()
	ApiResponse(c, &Response{Data: tabs})
}
