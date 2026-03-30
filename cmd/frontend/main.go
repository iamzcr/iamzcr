package main

import (
	"iamzcr/config"
	"iamzcr/handlers/frontend"
	"iamzcr/middleware"
	"iamzcr/models"
	"log"
	"strconv"

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

		api.GET("/articles", func(c *gin.Context) {
			page := c.DefaultQuery("page", "1")
			pageSize := c.DefaultQuery("page_size", "10")
			p, _ := strconv.Atoi(page)
			ps, _ := strconv.Atoi(pageSize)

			var articles []models.Article
			var total int64
			models.DB.Model(&models.Article{}).Where("status = ?", 1).Count(&total)
			models.DB.Where("status = ?", 1).Offset((p - 1) * ps).Limit(ps).Order("create_time DESC").Find(&articles)

			type ArticleWithTags struct {
				models.Article
				Tags []models.Tags `json:"tags"`
			}

			result := make([]ArticleWithTags, len(articles))
			for i, article := range articles {
				result[i].Article = article
				var articleTags []models.ArticleTags
				models.DB.Where("aid = ?", article.ID).Find(&articleTags)
				var tagIds []int
				for _, at := range articleTags {
					tagIds = append(tagIds, at.Tid)
				}
				if len(tagIds) > 0 {
					var tags []models.Tags
					models.DB.Where("id IN ?", tagIds).Find(&tags)
					result[i].Tags = tags
				}
			}

			c.JSON(200, gin.H{
				"code":    0,
				"message": "success",
				"data": gin.H{
					"list":  result,
					"total": total,
				},
			})
		})
		api.GET("/articles/:id", frontendHandler.GetArticle)

		api.GET("/categories", frontendHandler.GetCategories)
		api.GET("/directories", frontendHandler.GetDirectories)
		api.GET("/tags", frontendHandler.GetTags)
		api.GET("/website", frontendHandler.GetWebsite)
	}

	log.Printf("Frontend API server starting on port %s...", cfg.FrontendPort)
	r.Run(":" + cfg.FrontendPort)
}
