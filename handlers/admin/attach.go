package admin

import (
	"iamzcr/config"
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttachHandler struct {
	svc *svc.AttachService
}

func NewAttachHandler(s *svc.AttachService) *AttachHandler {
	return &AttachHandler{svc: s}
}

func (h *AttachHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请选择文件"})
		return
	}

	attach, err := h.svc.Upload(file, config.Cfg.AssetDir())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": attach})
}

func (h *AttachHandler) List(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("page_size", "10"))

	attaches, total, err := h.svc.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{"list": attaches, "total": total}})
}

func (h *AttachHandler) Get(c *gin.Context) {
	id := parseInt(c.Param("id"))

	attach, err := h.svc.Get(id)
	if err != nil {
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
	if err := h.svc.Create(&attach); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": attach})
}

func (h *AttachHandler) Update(c *gin.Context) {
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

func (h *AttachHandler) Delete(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
