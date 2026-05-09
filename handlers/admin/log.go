package admin

import (
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogHandler struct {
	svc *svc.LogService
}

func NewLogHandler(s *svc.LogService) *LogHandler {
	return &LogHandler{svc: s}
}

func (h *LogHandler) List(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("page_size", "10"))
	username := c.Query("username")

	logs, total, err := h.svc.List(page, pageSize, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{"list": logs, "total": total}})
}

func (h *LogHandler) Get(c *gin.Context) {
	id := parseInt(c.Param("id"))

	log, err := h.svc.Get(id)
	if err != nil {
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
	if err := h.svc.Create(&log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": log})
}

func (h *LogHandler) Delete(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
