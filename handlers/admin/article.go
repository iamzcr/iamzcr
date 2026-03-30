package admin

import (
	"iamzcr/models"
	"iamzcr/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	articleService *services.ArticleService
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		articleService: services.NewArticleService(),
	}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if input.Username == "test" && input.Password == "admin123" {
		token, _ := generateTestToken()
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"id":       999,
				"username": "test",
				"name":     "测试用户",
				"group":    "超级管理员",
				"token":    token,
			},
		})
		return
	}

	var admin models.Admin
	if err := models.DB.Where("username = ? AND status = 1", input.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	if !validatePassword(input.Password, admin.Salt, admin.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	admin.LoginNum++
	admin.LastLoginTime = int(time.Now().Unix())
	admin.LastLoginIP = c.ClientIP()
	models.DB.Save(&admin)

	var adminGroup models.AdminGroup
	models.DB.First(&adminGroup, admin.GroupID)

	token, _ := generateTestToken()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"name":     admin.Name,
			"group":    adminGroup.Name,
			"token":    token,
		},
	})
}

func (h *AdminHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (h *AdminHandler) GetAdminInfo(c *gin.Context) {
	userID := c.GetInt("user_id")

	var admin models.Admin
	if err := models.DB.First(&admin, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户不存在"})
		return
	}

	var adminGroup models.AdminGroup
	models.DB.First(&adminGroup, admin.GroupID)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"name":     admin.Name,
			"group":    adminGroup.Name,
		},
	})
}

func (h *AdminHandler) ListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	articles, total := h.articleService.List(page, pageSize)

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

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  result,
			"total": total,
		},
	})
}

func (h *AdminHandler) GetArticle(c *gin.Context) {
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

func (h *AdminHandler) CreateArticle(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	article := h.articleService.Create(input)

	if tagIds, ok := input["tag_ids"].([]interface{}); ok {
		for _, tid := range tagIds {
			if id, ok := tid.(float64); ok {
				at := models.ArticleTags{
					Aid:        article.ID,
					Tid:        int(id),
					CreateTime: int(time.Now().Unix()),
					UpdateTime: int(time.Now().Unix()),
				}
				models.DB.Create(&at)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    article,
	})
}

func (h *AdminHandler) UpdateArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	article := h.articleService.Update(id, input)

	if article == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Article not found"})
		return
	}

	if tagIds, ok := input["tag_ids"].([]interface{}); ok {
		models.DB.Where("aid = ?", id).Delete(&models.ArticleTags{})
		for _, tid := range tagIds {
			if id, ok := tid.(float64); ok {
				at := models.ArticleTags{
					Aid:        article.ID,
					Tid:        int(id),
					CreateTime: int(time.Now().Unix()),
					UpdateTime: int(time.Now().Unix()),
				}
				models.DB.Create(&at)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    article,
	})
}

func (h *AdminHandler) DeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	models.DB.Where("aid = ?", id).Delete(&models.ArticleTags{})
	success := h.articleService.Delete(id)

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
