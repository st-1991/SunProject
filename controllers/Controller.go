package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func ApiResponse(c *gin.Context, response *Response)  {
	if response.Message == "" {
		response.Message = "success"
	}
	c.JSON(http.StatusOK, response)
}