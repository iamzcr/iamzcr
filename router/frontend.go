package router

import (
	"iamzcr/config"
	frontH "iamzcr/handlers/frontend"
	frontS "iamzcr/services/frontend"

	"github.com/gin-gonic/gin"
)

func SetupFrontendRoutes(r *gin.Engine, cfg *config.Config) {
	articleSvc := frontS.NewArticleService()
	categorySvc := frontS.NewCategoryService()
	directorySvc := frontS.NewDirectoryService()
	tagsSvc := frontS.NewTagsService()
	websiteSvc := frontS.NewWebsiteService()
	messageSvc := frontS.NewMessageService()

	api := r.Group("/api")
	{
		frontendHandler := frontH.NewFrontendHandler(articleSvc, categorySvc, directorySvc, tagsSvc, websiteSvc, messageSvc)

		api.GET("/articles", frontendHandler.ListArticles)
		api.GET("/articles/:id", frontendHandler.GetArticle)

		api.GET("/categories", frontendHandler.GetCategories)
		api.GET("/directories", frontendHandler.GetDirectories)
		api.GET("/tags", frontendHandler.GetTags)
		api.GET("/website", frontendHandler.GetWebsite)

		api.GET("/messages", frontendHandler.GetMessages)
		api.POST("/messages", frontendHandler.CreateMessage)
	}

	r.Static("/api/docs", "./docs")
	r.Static("/asset", cfg.AssetDir())
}
