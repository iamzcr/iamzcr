package main

import (
	"iamzcr/config"
	"iamzcr/middleware"
	"iamzcr/models"
	"iamzcr/router"
	"log"

	"github.com/gin-gonic/gin"
)

// @title           堆栈人生 Frontend API
// @version         1.0
// @description     堆栈人生博客系统 - 前台公开 API，供 Web 前端和微信小程序使用
// @host            localhost:8082
// @BasePath        /api

func main() {
	cfg := config.Load()
	models.InitDB(cfg)

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())

	router.SetupFrontendRoutes(r, cfg)

	log.Printf("Frontend API server starting on port %s...", cfg.FrontendPort)
	r.Run(":" + cfg.FrontendPort)
}
