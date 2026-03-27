package handlers

import (
	"blog/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebsiteHandler struct{}

func NewWebsiteHandler() *WebsiteHandler {
	return &WebsiteHandler{}
}

func (h *WebsiteHandler) Get(c *gin.Context) {
	var settings []models.Website
	models.DB.Find(&settings)

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    result,
	})
}

func (h *WebsiteHandler) Update(c *gin.Context) {
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	for key, value := range input {
		var site models.Website
		if err := models.DB.Where("`key` = ?", key).First(&site).Error; err == nil {
			site.Value = value
			models.DB.Save(&site)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
