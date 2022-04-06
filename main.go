package main

import (
	"SunProject/config"
	"SunProject/router"

	"github.com/gin-gonic/gin"
)

func main() {
	defer config.Redis.Close()
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()

	//if !config.DB.Migrator().HasTable(&models.User{}) {
	//	if err := config.DB.Migrator().CreateTable(&models.User{}); err != nil {
	//		panic(err)
	//	}
	//}

	engine := gin.Default()
	engine.Use(config.LoggerToFile())
	route := router.Route{Engine: engine}
	route.Run()
	error := engine.Run(":6868")
	if error != nil {
		panic("服务启动失败")
	}
}
