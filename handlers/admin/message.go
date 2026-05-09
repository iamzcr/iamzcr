package admin

import (
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageHandler struct {
	svc *svc.MessageService
}

func NewMessageHandler(s *svc.MessageService) *MessageHandler {
	return &MessageHandler{svc: s}
}

func (h *MessageHandler) List(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("page_size", "10"))

	messages, total, err := h.svc.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{"list": messages, "total": total}})
}

func (h *MessageHandler) Get(c *gin.Context) {
	id := parseInt(c.Param("id"))

	message, err := h.svc.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": message})
}

func (h *MessageHandler) Create(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	if err := h.svc.Create(&message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": message})
}

func (h *MessageHandler) Update(c *gin.Context) {
	id := parseInt(c.Param("id"))
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	if err := h.svc.Update(id, data); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

func (h *MessageHandler) Delete(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
