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
	{
		api.GET("/send_sms", controllers.SendSms)
		api.POST("/login", controllers.Login)
		api.GET("/tabs", controllers.Tabs)
	}

	apiNeedToken := r.Engine.Group("/api").Use(middleware.KeepLogin())
	{
		apiNeedToken.GET("/user/info", controllers.UserInfo)
		apiNeedToken.POST("/user/edit", controllers.EditUser)
		apiNeedToken.GET("/users", controllers.UserList)
		apiNeedToken.POST("/upload_file", controllers.UploadFile)

		apiNeedToken.POST("/dynamic/add", controllers.AddDynamic)
		apiNeedToken.GET("/dynamic/recommend", controllers.RecommendDynamics)
		apiNeedToken.POST("/dynamic/thumb", controllers.DynamicThumbUp)
	}
}