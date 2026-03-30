package admin

import (
	"iamzcr/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DirectoryHandler struct{}

func NewDirectoryHandler() *DirectoryHandler {
	return &DirectoryHandler{}
}

func (h *DirectoryHandler) List(c *gin.Context) {
	var dirs []models.Directory
	cid := c.Query("cid")
	if cid != "" {
		models.DB.Where("cid = ?", cid).Order("weight DESC").Find(&dirs)
	} else {
		models.DB.Order("weight DESC").Find(&dirs)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    dirs,
	})
}

func (h *DirectoryHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var dir models.Directory

	if err := models.DB.First(&dir, id).Error; err != nil {
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

	dir.CreateTime = int(time.Now().Unix())
	dir.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Create(&dir).Error; err != nil {
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
	id := c.Param("id")
	var dir models.Directory

	if err := models.DB.First(&dir, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Directory not found"})
		return
	}

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

	dir.Cid = input.Cid
	dir.Type = input.Type
	dir.Parent = input.Parent
	dir.Mark = input.Mark
	dir.Author = input.Author
	dir.Name = input.Name
	dir.Weight = input.Weight
	dir.Status = input.Status
	dir.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&dir).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    dir,
	})
}

func (h *DirectoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := models.DB.Delete(&models.Directory{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
