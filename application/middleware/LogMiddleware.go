package middleware

import (
	"SunProject/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func LoggerToFile() gin.HandlerFunc {
	logger := config.Logger()
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		//             $response->header('Access-Control-Allow-Methods', 'GET,POST,OPTIONS,PUT,DELETED,PATCH');
		//            $response->header('Access-Control-Allow-Credentials', true);

		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETED,PATCH")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "3600")
		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusOK, "")
			c.Abort()
		}
		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		//日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}