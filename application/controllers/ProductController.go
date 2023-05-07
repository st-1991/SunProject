package controllers

import (
	"SunProject/application/models"
	"SunProject/config"
	"github.com/gin-gonic/gin"
)

func Products(c *gin.Context) {
	var p models.Product
	var result []struct {
		Id       int    `json:"id"`
		Title    string `json:"title"`
		Amount   string `json:"amount"`
		Integral int    `json:"integral"`
	}
	config.DB.Model(&p).Select("id", "title", "amount", "integral").Find(&result)
	ApiResponse(c, &Response{Data: result})
}
