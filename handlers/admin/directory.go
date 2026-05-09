package admin

import (
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DirectoryHandler struct {
	svc *svc.DirectoryService
}

func NewDirectoryHandler(s *svc.DirectoryService) *DirectoryHandler {
	return &DirectoryHandler{svc: s}
}

func (h *DirectoryHandler) List(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("page_size", "10"))
	cid := c.Query("cid")

	dirs, total, err := h.svc.List(page, pageSize, cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  dirs,
			"total": total,
		},
	})
}

func (h *DirectoryHandler) Get(c *gin.Context) {
	id := parseInt(c.Param("id"))

	dir, err := h.svc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Directory not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    dir,
	})
}

func (h *DirectoryHandler) Create(c *gin.Context) {
	var dir models.Directory
	if err := c.ShouldBindJSON(&dir); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if err := h.svc.Create(&dir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    dir,
	})
}

func (h *DirectoryHandler) Update(c *gin.Context) {
	id := parseInt(c.Param("id"))

	var input struct {
		Cid    int    `json:"cid"`
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

	dir := &models.Directory{
		Cid:    input.Cid,
		Type:   input.Type,
		Parent: input.Parent,
		Mark:   input.Mark,
		Author: input.Author,
		Name:   input.Name,
		Weight: input.Weight,
		Status: input.Status,
	}

	updated, err := h.svc.Update(id, dir)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Directory not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    updated,
	})
}

func (h *DirectoryHandler) Delete(c *gin.Context) {
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
