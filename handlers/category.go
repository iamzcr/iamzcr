package handlers

import (
	"blog/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct{}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

func (h *CategoryHandler) List(c *gin.Context) {
	var categories []models.Category
	models.DB.Order("weight DESC").Find(&categories)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    categories,
	})
}

func (h *CategoryHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := models.DB.First(&category, id).Error; err != nil {
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

	category.CreateTime = int(time.Now().Unix())
	category.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Create(&category).Error; err != nil {
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
	id := c.Param("id")
	var category models.Category

	if err := models.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Category not found"})
		return
	}

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

	category.Type = input.Type
	category.Parent = input.Parent
	category.Mark = input.Mark
	category.Author = input.Author
	category.Name = input.Name
	category.Weight = input.Weight
	category.Status = input.Status
	category.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    category,
	})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := models.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
