package frontend

import (
	"iamzcr/models"
	"iamzcr/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FrontendHandler struct {
	articleService *services.ArticleService
}

func NewFrontendHandler() *FrontendHandler {
	return &FrontendHandler{
		articleService: services.NewArticleService(),
	}
}

func (h *FrontendHandler) GetArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	article := h.articleService.GetByID(id)

	if article == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    article,
	})
}

func (h *FrontendHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	models.DB.Where("status = ?", 1).Find(&categories)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    categories,
	})
}

func (h *FrontendHandler) GetDirectories(c *gin.Context) {
	var directories []models.Directory
	models.DB.Where("status = ?", 1).Find(&directories)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    directories,
	})
}

func (h *FrontendHandler) GetTags(c *gin.Context) {
	var tags []models.Tags
	models.DB.Where("status = ?", 1).Find(&tags)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tags,
	})
}

func (h *FrontendHandler) GetWebsite(c *gin.Context) {
	var websites []models.Website
	models.DB.Find(&websites)

	data := make(map[string]string)
	for _, w := range websites {
		data[w.Key] = w.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}
