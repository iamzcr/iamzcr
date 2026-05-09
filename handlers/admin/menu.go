package admin

import (
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	svc *svc.MenuService
}

func NewMenuHandler(s *svc.MenuService) *MenuHandler {
	return &MenuHandler{svc: s}
}

func (h *MenuHandler) List(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("page_size", "10"))

	menus, total, err := h.svc.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  menus,
			"total": total,
		},
	})
}

func (h *MenuHandler) Get(c *gin.Context) {
	id := parseInt(c.Param("id"))

	menu, err := h.svc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Menu not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    menu,
	})
}

func (h *MenuHandler) Create(c *gin.Context) {
	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if err := h.svc.Create(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    menu,
	})
}

func (h *MenuHandler) Update(c *gin.Context) {
	id := parseInt(c.Param("id"))

	var input struct {
		Type   int    `json:"type"`
		Mark   string `json:"mark"`
		Author string `json:"author"`
		Name   string `json:"name"`
		Url    string `json:"url"`
		Parent int    `json:"parent"`
		Icon   string `json:"icon"`
		Weight int    `json:"weight"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	menu := &models.Menu{
		Type:   input.Type,
		Mark:   input.Mark,
		Author: input.Author,
		Name:   input.Name,
		Url:    input.Url,
		Parent: input.Parent,
		Icon:   input.Icon,
		Weight: input.Weight,
		Status: input.Status,
	}

	updated, err := h.svc.Update(id, menu)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Menu not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    updated,
	})
}

func (h *MenuHandler) Delete(c *gin.Context) {
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
