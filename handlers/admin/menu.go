package admin

import (
	"iamzcr/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct{}

func NewMenuHandler() *MenuHandler {
	return &MenuHandler{}
}

func (h *MenuHandler) List(c *gin.Context) {
	var menus []models.Menu
	models.DB.Order("weight DESC").Find(&menus)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    menus,
	})
}

func (h *MenuHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var menu models.Menu

	if err := models.DB.First(&menu, id).Error; err != nil {
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

	menu.CreateTime = int(time.Now().Unix())
	menu.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Create(&menu).Error; err != nil {
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
	id := c.Param("id")
	var menu models.Menu

	if err := models.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Menu not found"})
		return
	}

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

	menu.Type = input.Type
	menu.Mark = input.Mark
	menu.Author = input.Author
	menu.Name = input.Name
	menu.Url = input.Url
	menu.Parent = input.Parent
	menu.Icon = input.Icon
	menu.Weight = input.Weight
	menu.Status = input.Status
	menu.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    menu,
	})
}

func (h *MenuHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := models.DB.Delete(&models.Menu{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
