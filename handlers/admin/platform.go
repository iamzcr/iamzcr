package admin

import (
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlatformHandler struct {
	svc *svc.PlatformService
}

func NewPlatformHandler(s *svc.PlatformService) *PlatformHandler {
	return &PlatformHandler{svc: s}
}

func (h *PlatformHandler) List(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("page_size", "10"))

	platforms, total, err := h.svc.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  platforms,
			"total": total,
		},
	})
}

func (h *PlatformHandler) Get(c *gin.Context) {
	id := parseInt(c.Param("id"))

	platform, err := h.svc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "platform not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    platform,
	})
}

func (h *PlatformHandler) Create(c *gin.Context) {
	var platform models.Platform
	if err := c.ShouldBindJSON(&platform); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if err := h.svc.Create(&platform); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    platform,
	})
}

func (h *PlatformHandler) Update(c *gin.Context) {
	id := parseInt(c.Param("id"))

	var input struct {
		Mark   string `json:"mark"`
		Name   string `json:"name"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	platform := &models.Platform{
		Mark:   input.Mark,
		Name:   input.Name,
		Status: input.Status,
	}

	updated, err := h.svc.Update(id, platform)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "platform not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    updated,
	})
}

func (h *PlatformHandler) Delete(c *gin.Context) {
	id := parseInt(c.Param("id"))

	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
