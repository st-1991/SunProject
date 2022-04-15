package main

import (
	"SunProject/application/middleware"
	"SunProject/application/models"
	"SunProject/config"
	"SunProject/router"

	"github.com/gin-gonic/gin"
)

func main() {
	defer config.Redis.Close()
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()
	gin.SetMode(gin.DebugMode)


	tables := models.Tables()
	for _, table := range tables {
		if !config.DB.Migrator().HasTable(table) {
			if err := config.DB.Migrator().CreateTable(table); err != nil {
				panic(err)
			}
		}
	}

	engine := gin.Default()
	engine.Use(middleware.LoggerToFile())
	route := router.Route{Engine: engine}
	route.Run()
	err := engine.Run(":6868")
	if err != nil {
		panic("服务启动失败")
	}
}
