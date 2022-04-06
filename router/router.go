package router

import (
	"SunProject/controllers"
	"SunProject/middleware"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Engine *gin.Engine
}

//路由路口
func (r *Route) Run() {
	api := r.Engine.Group("/api").Use(middleware.KeepLogin())
	{
		api.GET("/send_sms", controllers.SendSms)
		api.POST("/login", controllers.Login)
		api.GET("/user_info", controllers.UserInfo)
		api.GET("/users", controllers.UserList)
	}
}