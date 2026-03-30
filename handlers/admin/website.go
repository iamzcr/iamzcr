package admin

import (
	"iamzcr/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WebsiteHandler struct{}

func NewWebsiteHandler() *WebsiteHandler {
	return &WebsiteHandler{}
}

func (h *WebsiteHandler) Get(c *gin.Context) {
	var websites []models.Website
	models.DB.Find(&websites)

	data := make(map[string]interface{})
	for _, w := range websites {
		data[w.Key] = w.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

func (h *WebsiteHandler) Update(c *gin.Context) {
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	for key, value := range input {
		var website models.Website
		if err := models.DB.Where("`key` = ?", key).First(&website).Error; err == nil {
			website.Value = value
			website.UpdateTime = int(time.Now().Unix())
			models.DB.Save(&website)
		} else {
			website := models.Website{
				Key:        key,
				Value:      value,
				Staus:      1,
				CreateTime: int(time.Now().Unix()),
				UpdateTime: int(time.Now().Unix()),
			}
			models.DB.Create(&website)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
