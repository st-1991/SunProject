package middleware

import "github.com/gin-gonic/gin"

func KeepLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("keep-token")
		c.Set("token", token)
		c.Next()
	}
}