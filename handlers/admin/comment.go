package admin

import (
	"iamzcr/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct{}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{}
}

func (h *CommentHandler) List(c *gin.Context) {
	var comments []models.Comment
	models.DB.Order("create_time DESC").Find(&comments)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    comments,
	})
}

func (h *CommentHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var comment models.Comment

	if err := models.DB.First(&comment, id).Error; err != nil {
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

	comment.CreateTime = int(time.Now().Unix())
	comment.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Create(&comment).Error; err != nil {
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
	id := c.Param("id")
	var comment models.Comment

	if err := models.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Comment not found"})
		return
	}

	var input struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	comment.Content = input.Content
	comment.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    comment,
	})
}

func (h *CommentHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := models.DB.Delete(&models.Comment{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
