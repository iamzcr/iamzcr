package admin

import (
	"iamzcr/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReadHandler struct{}

func NewReadHandler() *ReadHandler {
	return &ReadHandler{}
}

func (h *ReadHandler) List(c *gin.Context) {
	var reads []models.Read
	query := models.DB.Order("id DESC")

	aid := c.Query("aid")
	if aid != "" {
		query = query.Where("aid = ?", aid)
	}

	if err := query.Find(&reads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": reads})
}

func (h *ReadHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var read models.Read
	if err := models.DB.First(&read, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": read})
}

func (h *ReadHandler) Create(c *gin.Context) {
	var read models.Read
	if err := c.ShouldBindJSON(&read); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	read.CreateTime = int(time.Now().Unix())
	if err := models.DB.Create(&read).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": read})
}

func (h *ReadHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := models.DB.Delete(&models.Read{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
