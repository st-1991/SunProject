package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func ApiResponse(c *gin.Context, response *Response)  {
	if response.Msg == "" {
		response.Msg = "success"
	}
	c.JSON(http.StatusOK, response)
}