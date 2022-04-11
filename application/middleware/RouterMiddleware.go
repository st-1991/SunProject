package middleware

import (
	"SunProject/application/controllers"
	"SunProject/application/service"
	"github.com/gin-gonic/gin"
)

func KeepLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("keep-token")
		isLogin := false
		if token != "" {
			token := service.Token(token)
			if ok := token.Validate(); !ok {
				controllers.ApiResponse(c, &controllers.Response{
					Code: 4444,
					Msg: "token验证失败，请重新登录",
				})
				c.Abort()
			}
			if userData, err := token.GetUserInfo("user"); err == nil {
				c.Set("userId", userData.ID)
			}
		}
		isLogin = true
		c.Set("isLogin", isLogin)
		c.Header("Access-Control-Allow-Headers", "*")
		c.Next()
	}
}