package main

import (
	"iamzcr/config"
	"iamzcr/middleware"
	"iamzcr/models"
	"iamzcr/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	models.InitDB(cfg)

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())

	router.SetupAdminRoutes(r, cfg)

	log.Printf("Admin API server starting on port %s...", cfg.AdminPort)
	r.Run(":" + cfg.AdminPort)
}
