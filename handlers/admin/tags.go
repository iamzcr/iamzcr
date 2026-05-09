package admin

import (
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagsHandler struct {
	svc *svc.TagsService
}

func NewTagsHandler(s *svc.TagsService) *TagsHandler {
	return &TagsHandler{svc: s}
}

func (h *TagsHandler) List(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("page_size", "10"))

	tags, total, err := h.svc.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  tags,
			"total": total,
		},
	})
}

func (h *TagsHandler) Get(c *gin.Context) {
	id := parseInt(c.Param("id"))

	tag, err := h.svc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Tag not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tag,
	})
}

func (h *TagsHandler) Create(c *gin.Context) {
	var tag models.Tags
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if err := h.svc.Create(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tag,
	})
}

func (h *TagsHandler) Update(c *gin.Context) {
	id := parseInt(c.Param("id"))

	var input struct {
		Type   string `json:"type"`
		Mark   string `json:"mark"`
		Author string `json:"author"`
		Name   string `json:"name"`
		Weight int    `json:"weight"`
		Status int    `json:"status"`
		IsHot  int    `json:"is_hot"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	tag := &models.Tags{
		Type:   input.Type,
		Mark:   input.Mark,
		Author: input.Author,
		Name:   input.Name,
		Weight: input.Weight,
		Status: input.Status,
		IsHot:  input.IsHot,
	}

	updated, err := h.svc.Update(id, tag)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Tag not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    updated,
	})
}

func (h *TagsHandler) Delete(c *gin.Context) {
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
