package main

import (
	"blog/config"
	"blog/handlers"
	"blog/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	models.InitDB(cfg)

	r := gin.Default()
	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		adminHandler := handlers.NewAdminHandler()
		api.POST("/login", adminHandler.Login)
		api.POST("/logout", adminHandler.Logout)
		api.GET("/admin/info", adminHandler.GetAdminInfo)

		articleHandler := handlers.NewArticleHandler()
		api.GET("/articles", articleHandler.List)
		api.GET("/articles/:id", articleHandler.Get)
		api.POST("/articles", articleHandler.Create)
		api.PUT("/articles/:id", articleHandler.Update)
		api.DELETE("/articles/:id", articleHandler.Delete)

		categoryHandler := handlers.NewCategoryHandler()
		api.GET("/categories", categoryHandler.List)
		api.GET("/categories/:id", categoryHandler.Get)
		api.POST("/categories", categoryHandler.Create)
		api.PUT("/categories/:id", categoryHandler.Update)
		api.DELETE("/categories/:id", categoryHandler.Delete)

		commentHandler := handlers.NewCommentHandler()
		api.GET("/comments", commentHandler.List)
		api.GET("/comments/:id", commentHandler.Get)
		api.POST("/comments", commentHandler.Create)
		api.PUT("/comments/:id", commentHandler.Update)
		api.DELETE("/comments/:id", commentHandler.Delete)

		menuHandler := handlers.NewMenuHandler()
		api.GET("/menus", menuHandler.List)
		api.GET("/menus/:id", menuHandler.Get)
		api.POST("/menus", menuHandler.Create)
		api.PUT("/menus/:id", menuHandler.Update)
		api.DELETE("/menus/:id", menuHandler.Delete)

		tagsHandler := handlers.NewTagsHandler()
		api.GET("/tags", tagsHandler.List)
		api.GET("/tags/:id", tagsHandler.Get)
		api.POST("/tags", tagsHandler.Create)
		api.PUT("/tags/:id", tagsHandler.Update)
		api.DELETE("/tags/:id", tagsHandler.Delete)

		directoryHandler := handlers.NewDirectoryHandler()
		api.GET("/directories", directoryHandler.List)
		api.GET("/directories/:id", directoryHandler.Get)
		api.POST("/directories", directoryHandler.Create)
		api.PUT("/directories/:id", directoryHandler.Update)
		api.DELETE("/directories/:id", directoryHandler.Delete)

		websiteHandler := handlers.NewWebsiteHandler()
		api.GET("/website", websiteHandler.Get)
		api.PUT("/website", websiteHandler.Update)
	}

	log.Printf("Admin API server starting on port %s...", cfg.AdminPort)
	r.Run(":" + cfg.AdminPort)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
