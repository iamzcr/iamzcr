package admin

import (
	"iamzcr/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebsiteHandler struct{}

func NewWebsiteHandler() *WebsiteHandler {
	return &WebsiteHandler{}
}

func (h *WebsiteHandler) List(c *gin.Context) {
	var websites []models.Website
	models.DB.Order("id DESC").Find(&websites)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    websites,
	})
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

func (h *WebsiteHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	if err := models.DB.Delete(&models.Website{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
