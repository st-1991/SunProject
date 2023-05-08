package router

import (
	"SunProject/application/controllers"
	"SunProject/application/middleware"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Engine *gin.Engine
}

// Run 路由路口
func (r *Route) Run() {

	api := r.Engine.Group("/api")
	api.GET("/keep/notify", controllers.Notify)
	api.POST("/keep/notify", controllers.Notify)

	api.Use(middleware.VerifySign()) //签名验证

	api.GET("/send_sms", controllers.SendSms)
	api.POST("/login", controllers.Login)
	api.GET("/product", controllers.Products)

	//api.GET("/tabs", controllers.Tabs)
	apiNeedToken := api.Group("").Use(middleware.KeepLogin())
	{
		apiNeedToken.GET("/user/info", controllers.UserInfo)
		apiNeedToken.POST("/user/edit", controllers.EditUser)
		apiNeedToken.GET("/users", controllers.UserList)

		apiNeedToken.POST("/completions", controllers.Completions)
		apiNeedToken.POST("/images/generations", controllers.CreateImages)

		apiNeedToken.POST("/prepay", controllers.PrePay)
		apiNeedToken.GET("/order/status", controllers.OrderStatus)

		apiNeedToken.GET("/integral_logs", controllers.GetIntegralLogs)
	}
}