package admin

import (
	"iamzcr/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogHandler struct{}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (h *LogHandler) List(c *gin.Context) {
	var logs []models.Log
	query := models.DB.Order("id DESC")

	username := c.Query("username")
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	if err := query.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": logs})
}

func (h *LogHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var log models.Log
	if err := models.DB.First(&log, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": log})
}

func (h *LogHandler) Create(c *gin.Context) {
	var log models.Log
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	log.CreateTime = int(time.Now().Unix())
	if err := models.DB.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": log})
}

func (h *LogHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := models.DB.Delete(&models.Log{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
