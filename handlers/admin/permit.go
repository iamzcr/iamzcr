package admin

import (
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermitHandler struct {
	svc *svc.PermitService
}

func NewPermitHandler(s *svc.PermitService) *PermitHandler {
	return &PermitHandler{svc: s}
}

func (h *PermitHandler) List(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("page_size", "10"))

	permits, total, err := h.svc.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{"list": permits, "total": total}})
}

func (h *PermitHandler) Get(c *gin.Context) {
	id := parseInt(c.Param("id"))

	permit, err := h.svc.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": permit})
}

func (h *PermitHandler) Create(c *gin.Context) {
	var permit models.Permit
	if err := c.ShouldBindJSON(&permit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	if err := h.svc.Create(&permit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": permit})
}

func (h *PermitHandler) Update(c *gin.Context) {
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

func (h *PermitHandler) Delete(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
