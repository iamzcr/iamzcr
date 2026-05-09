package admin

import (
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	svc *svc.CommentService
}

func NewCommentHandler(s *svc.CommentService) *CommentHandler {
	return &CommentHandler{svc: s}
}

func (h *CommentHandler) List(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("page_size", "10"))

	comments, total, err := h.svc.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  comments,
			"total": total,
		},
	})
}

func (h *CommentHandler) Get(c *gin.Context) {
	id := parseInt(c.Param("id"))

	comment, err := h.svc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    comment,
	})
}

func (h *CommentHandler) Create(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if err := h.svc.Create(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    comment,
	})
}

func (h *CommentHandler) Update(c *gin.Context) {
	id := parseInt(c.Param("id"))

	var input struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	updated, err := h.svc.Update(id, input.Content)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    updated,
	})
}

func (h *CommentHandler) Delete(c *gin.Context) {
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
