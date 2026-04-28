package main

import (
	"iamzcr/config"
	"iamzcr/handlers/frontend"
	"iamzcr/middleware"
	"iamzcr/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	models.InitDB(cfg)

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())

	api := r.Group("/api")
	{
		frontendHandler := frontend.NewFrontendHandler()

		api.GET("/articles", frontendHandler.ListArticles)
		api.GET("/articles/:id", frontendHandler.GetArticle)

		api.GET("/categories", frontendHandler.GetCategories)
		api.GET("/directories", frontendHandler.GetDirectories)
		api.GET("/tags", frontendHandler.GetTags)
		api.GET("/website", frontendHandler.GetWebsite)
	}

	r.Static("/asset", cfg.AssetDir())

	log.Printf("Frontend API server starting on port %s...", cfg.FrontendPort)
	r.Run(":" + cfg.FrontendPort)
}
