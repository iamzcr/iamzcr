package handlers

import (
	"blog/models"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var admin models.Admin
	if err := models.DB.Where("username = ? AND status = 1", input.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	password := input.Password + admin.Salt
	hash := md5.Sum([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])

	if passwordHash != admin.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	admin.LoginNum++
	admin.LastLoginTime = int(time.Now().Unix())
	admin.LastLoginIP = c.ClientIP()
	models.DB.Save(&admin)

	var adminGroup models.AdminGroup
	models.DB.First(&adminGroup, admin.GroupID)

	token := hex.EncodeToString(hash[:]) + "-" + time.Now().Format("20060102150405")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"name":     admin.Name,
			"group":    adminGroup.Name,
			"token":    token,
		},
	})
}

func (h *AdminHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (h *AdminHandler) GetAdminInfo(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	var admin models.Admin
	if err := models.DB.First(&admin, 1).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户不存在"})
		return
	}

	var adminGroup models.AdminGroup
	models.DB.First(&adminGroup, admin.GroupID)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"name":     admin.Name,
			"group":    adminGroup.Name,
		},
	})
}
