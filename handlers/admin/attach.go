package admin

import (
	"iamzcr/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttachHandler struct{}

func NewAttachHandler() *AttachHandler {
	return &AttachHandler{}
}

func (h *AttachHandler) List(c *gin.Context) {
	var attaches []models.Attach
	if err := models.DB.Order("id DESC").Find(&attaches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": attaches})
}

func (h *AttachHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var attach models.Attach
	if err := models.DB.First(&attach, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": attach})
}

func (h *AttachHandler) Create(c *gin.Context) {
	var attach models.Attach
	if err := c.ShouldBindJSON(&attach); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	attach.CreateTime = int(time.Now().Unix())
	if err := models.DB.Create(&attach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": attach})
}

func (h *AttachHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var attach models.Attach
	if err := models.DB.First(&attach, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
		return
	}
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	delete(data, "id")
	delete(data, "create_time")
	data["update_time"] = int(time.Now().Unix())
	if err := models.DB.Model(&attach).Updates(data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": attach})
}

func (h *AttachHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := models.DB.Delete(&models.Attach{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
