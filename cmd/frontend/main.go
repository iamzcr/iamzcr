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
		articleHandler := handlers.NewArticleHandler()
		api.GET("/articles", articleHandler.List)
		api.GET("/articles/:id", articleHandler.Get)

		categoryHandler := handlers.NewCategoryHandler()
		api.GET("/categories", categoryHandler.List)

		directoryHandler := handlers.NewDirectoryHandler()
		api.GET("/directories", directoryHandler.List)

		tagsHandler := handlers.NewTagsHandler()
		api.GET("/tags", tagsHandler.List)
	}

	log.Printf("Frontend API server starting on port %s...", cfg.FrontendPort)
	r.Run(":" + cfg.FrontendPort)
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
