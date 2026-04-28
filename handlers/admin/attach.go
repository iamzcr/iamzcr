package admin

import (
	"fmt"
	"iamzcr/config"
	"iamzcr/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttachHandler struct{}

func NewAttachHandler() *AttachHandler {
	return &AttachHandler{}
}

func (h *AttachHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请选择文件"})
		return
	}

	uploadDir := config.Cfg.AssetDir()
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", time.Now().Unix(), time.Now().UnixNano()%100000, ext)
	savePath := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "文件读取失败"})
		return
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "文件保存失败"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "文件写入失败"})
		return
	}

	link := "/asset/" + filename

	var cdn models.Website
	if err := models.DB.Where("`key` = ?", "cdn_url").First(&cdn).Error; err == nil && cdn.Value != "" {
		link = cdn.Value + "/asset/" + filename
	}

	attach := models.Attach{
		Name:       file.Filename,
		Link:       link,
		Path:       savePath,
		Status:     1,
		Type:       1,
		CreateTime: int(time.Now().Unix()),
	}

	if err := models.DB.Create(&attach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": attach})
}

func (h *AttachHandler) List(c *gin.Context) {
	var attaches []models.Attach
	var total int64

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	models.DB.Model(&models.Attach{}).Count(&total)
	models.DB.Order("id DESC").Limit(parseInt(pageSize)).Offset((parseInt(page) - 1) * parseInt(pageSize)).Find(&attaches)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{"list": attaches, "total": total}})
}

func (h *AttachHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var attach models.Attach
	if err := models.DB.First(&attach, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": attach})
}

func (h *AttachHandler) Create(c *gin.Context) {
	var attach models.Attach
	if err := c.ShouldBindJSON(&attach); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	attach.CreateTime = int(time.Now().Unix())
	if err := models.DB.Create(&attach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": attach})
}

func (h *AttachHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var attach models.Attach
	if err := models.DB.First(&attach, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
		return
	}
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	delete(data, "id")
	delete(data, "create_time")
	data["update_time"] = int(time.Now().Unix())
	if err := models.DB.Model(&attach).Updates(data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": attach})
}

func (h *AttachHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := models.DB.Delete(&models.Attach{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
