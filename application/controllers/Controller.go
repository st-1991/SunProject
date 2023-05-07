package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Response struct {
	Code int `json:"code"`
	Msg string `json:"message"`
	Data interface{} `json:"data"`
}

func ApiResponse(c *gin.Context, response *Response)  {
	if response.Msg == "" {
		response.Msg = "success"
	}
	if  response.Data == nil {
		response.Data = make(map[string]string)
	}
	c.JSON(http.StatusOK, response)
}

func ApiError(c *gin.Context, response *Response, status ...int)  {
	if  response.Data == nil {
		response.Data = make(map[string]string)
	}
	httpStatus := http.StatusBadRequest
	if len(status) > 0 {
		httpStatus = status[0]
	}
	c.JSON(httpStatus, response)
}

func GetReaIp(c *gin.Context) string {
	realClientIp := c.GetHeader("X-Forwarded-For")
	if realClientIp != "" {
		ips := strings.Split(realClientIp, ", ")
		if len(ips) > 1 {
			realClientIp = ips[0]
		}
	} else {
		realClientIp = c.ClientIP()
	}
	return realClientIp
}