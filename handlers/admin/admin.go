package admin

import (
	"iamzcr/models"
	"net/http"
	"time"

	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *AdminHandler) ListAdmins(c *gin.Context) {
	var admins []models.Admin
	if err := models.DB.Order("id DESC").Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	for i := range admins {
		admins[i].Password = ""
		admins[i].Salt = ""
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": admins})
}

func (h *AdminHandler) GetAdmin(c *gin.Context) {
	id := c.Param("id")
	var admin models.Admin
	if err := models.DB.First(&admin, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	admin.Password = ""
	admin.Salt = ""
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": admin})
}

func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	salt := fmt.Sprintf("%d", time.Now().UnixNano())
	password := admin.Password + salt
	hashed := md5.Sum([]byte(password))
	admin.Salt = salt
	admin.Password = hex.EncodeToString(hashed[:])
	admin.CreateTime = int(time.Now().Unix())
	if err := models.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	admin.Password = ""
	admin.Salt = ""
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": admin})
}

func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	id := c.Param("id")
	var admin models.Admin
	if err := models.DB.First(&admin, id).Error; err != nil {
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
	delete(data, "password")
	delete(data, "salt")
	data["update_time"] = int(time.Now().Unix())
	if err := models.DB.Model(&admin).Updates(data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": admin})
}

func (h *AdminHandler) ChangeAdminPassword(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	var admin models.Admin
	if err := models.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
		return
	}

	salt := fmt.Sprintf("%d", time.Now().UnixNano())
	newPwd := input.NewPassword + salt
	newHash := md5.Sum([]byte(newPwd))

	admin.Salt = salt
	admin.Password = hex.EncodeToString(newHash[:])
	admin.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "密码修改成功"})
}

func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	id := c.Param("id")
	if err := models.DB.Delete(&models.Admin{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

func (h *AdminHandler) ChangePassword(c *gin.Context) {
	userID := c.GetInt("user_id")
	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	var admin models.Admin
	if err := models.DB.First(&admin, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
		return
	}

	oldPwd := input.OldPassword + admin.Salt
	oldHash := md5.Sum([]byte(oldPwd))
	if hex.EncodeToString(oldHash[:]) != admin.Password {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "原密码错误"})
		return
	}

	salt := fmt.Sprintf("%d", time.Now().UnixNano())
	newPwd := input.NewPassword + salt
	newHash := md5.Sum([]byte(newPwd))

	admin.Salt = salt
	admin.Password = hex.EncodeToString(newHash[:])
	admin.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "密码修改成功"})
}

type AdminGroupHandler struct{}

func NewAdminGroupHandler() *AdminGroupHandler {
	return &AdminGroupHandler{}
}

func (h *AdminGroupHandler) List(c *gin.Context) {
	var groups []models.AdminGroup
	if err := models.DB.Order("id DESC").Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": groups})
}

func (h *AdminGroupHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var group models.AdminGroup
	if err := models.DB.First(&group, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": group})
}

func (h *AdminGroupHandler) Create(c *gin.Context) {
	var group models.AdminGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	group.CreateTime = int(time.Now().Unix())
	if err := models.DB.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": group})
}

func (h *AdminGroupHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var group models.AdminGroup
	if err := models.DB.First(&group, id).Error; err != nil {
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
	if err := models.DB.Model(&group).Updates(data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": group})
}

func (h *AdminGroupHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := models.DB.Delete(&models.AdminGroup{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
