package main

import (
	"iamzcr/config"
	"iamzcr/handlers/frontend"
	"iamzcr/middleware"
	"iamzcr/models"
	svc "iamzcr/services/frontend"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	models.InitDB(cfg)

	articleSvc := svc.NewArticleService()
	categorySvc := svc.NewCategoryService()
	directorySvc := svc.NewDirectoryService()
	tagsSvc := svc.NewTagsService()
	websiteSvc := svc.NewWebsiteService()

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())

	api := r.Group("/api")
	{
		frontendHandler := frontend.NewFrontendHandler(articleSvc, categorySvc, directorySvc, tagsSvc, websiteSvc)

		api.GET("/articles", frontendHandler.ListArticles)
		api.GET("/articles/:id", frontendHandler.GetArticle)

		api.GET("/categories", frontendHandler.GetCategories)
		api.GET("/directories", frontendHandler.GetDirectories)
		api.GET("/tags", frontendHandler.GetTags)
		api.GET("/website", frontendHandler.GetWebsite)
	}

	r.Static("/cdn/asset", cfg.AssetDir())

	log.Printf("Frontend API server starting on port %s...", cfg.FrontendPort)
	r.Run(":" + cfg.FrontendPort)
}
