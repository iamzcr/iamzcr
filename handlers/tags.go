package handlers

import (
	"blog/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TagsHandler struct{}

func NewTagsHandler() *TagsHandler {
	return &TagsHandler{}
}

func (h *TagsHandler) List(c *gin.Context) {
	var tags []models.Tags
	models.DB.Order("weight DESC").Find(&tags)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tags,
	})
}

func (h *TagsHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var tag models.Tags

	if err := models.DB.First(&tag, id).Error; err != nil {
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

	tag.CreateTime = int(time.Now().Unix())
	tag.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Create(&tag).Error; err != nil {
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
	id := c.Param("id")
	var tag models.Tags

	if err := models.DB.First(&tag, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Tag not found"})
		return
	}

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

	tag.Type = input.Type
	tag.Mark = input.Mark
	tag.Author = input.Author
	tag.Name = input.Name
	tag.Weight = input.Weight
	tag.Status = input.Status
	tag.IsHot = input.IsHot
	tag.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tag,
	})
}

func (h *TagsHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := models.DB.Delete(&models.Tags{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
