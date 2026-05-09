package admin

import (
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	svc *svc.CategoryService
}

func NewCategoryHandler(s *svc.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: s}
}

func (h *CategoryHandler) List(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("page_size", "10"))

	categories, total, err := h.svc.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  categories,
			"total": total,
		},
	})
}

func (h *CategoryHandler) Get(c *gin.Context) {
	id := parseInt(c.Param("id"))

	category, err := h.svc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    category,
	})
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if err := h.svc.Create(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    category,
	})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := parseInt(c.Param("id"))

	var input struct {
		Type   string `json:"type"`
		Parent string `json:"parent"`
		Mark   string `json:"mark"`
		Author string `json:"author"`
		Name   string `json:"name"`
		Weight int    `json:"weight"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	category := &models.Category{
		Type:   input.Type,
		Parent: input.Parent,
		Mark:   input.Mark,
		Author: input.Author,
		Name:   input.Name,
		Weight: input.Weight,
		Status: input.Status,
	}

	updated, err := h.svc.Update(id, category)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    updated,
	})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
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
