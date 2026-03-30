package admin

import (
	"iamzcr/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LangHandler struct{}

func NewLangHandler() *LangHandler {
	return &LangHandler{}
}

func (h *LangHandler) List(c *gin.Context) {
	var langs []models.Lang
	var total int64

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	models.DB.Model(&models.Lang{}).Count(&total)
	models.DB.Order("weight DESC").Limit(parseInt(pageSize)).Offset((parseInt(page) - 1) * parseInt(pageSize)).Find(&langs)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{"list": langs, "total": total}})
}

func (h *LangHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var lang models.Lang
	if err := models.DB.First(&lang, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": lang})
}

func (h *LangHandler) Create(c *gin.Context) {
	var lang models.Lang
	if err := c.ShouldBindJSON(&lang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	lang.CreateTime = int(time.Now().Unix())
	if err := models.DB.Create(&lang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": lang})
}

func (h *LangHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var lang models.Lang
	if err := models.DB.First(&lang, id).Error; err != nil {
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
	if err := models.DB.Model(&lang).Updates(data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": lang})
}

func (h *LangHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := models.DB.Delete(&models.Lang{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
