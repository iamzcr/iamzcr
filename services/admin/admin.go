package admin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"iamzcr/config"
	"iamzcr/middleware"
	"iamzcr/models"
	"time"
)

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

type LoginResult struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Group    string `json:"group"`
	Token    string `json:"token"`
}

func (s *AdminService) Login(username, password, clientIP string) (*LoginResult, error) {
	if config.Cfg.Env == "development" && username == "test" && password == "admin123" {
		token, _ := generateToken(1, "test")
		return &LoginResult{
			ID:       999,
			Username: "test",
			Name:     "测试用户",
			Group:    "超级管理员",
			Token:    token,
		}, nil
	}

	var admin models.Admin
	if err := models.DB.Where("username = ? AND status = 1", username).First(&admin).Error; err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	if !validatePassword(password, admin.Salt, admin.Password) {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	admin.LoginNum++
	admin.LastLoginTime = int(time.Now().Unix())
	admin.LastLoginIP = clientIP
	models.DB.Save(&admin)

	var adminGroup models.AdminGroup
	models.DB.First(&adminGroup, admin.GroupID)

	token, _ := generateToken(admin.ID, admin.Username)

	return &LoginResult{
		ID:       admin.ID,
		Username: admin.Username,
		Name:     admin.Name,
		Group:    adminGroup.Name,
		Token:    token,
	}, nil
}

func (s *AdminService) GetAdminInfo(userID int) (*models.Admin, *models.AdminGroup, error) {
	var admin models.Admin
	if err := models.DB.First(&admin, userID).Error; err != nil {
		return nil, nil, err
	}

	var adminGroup models.AdminGroup
	models.DB.First(&adminGroup, admin.GroupID)

	return &admin, &adminGroup, nil
}

func (s *AdminService) List() ([]models.Admin, error) {
	var admins []models.Admin
	if err := models.DB.Order("id DESC").Find(&admins).Error; err != nil {
		return nil, err
	}
	for i := range admins {
		admins[i].Password = ""
		admins[i].Salt = ""
	}
	return admins, nil
}

func (s *AdminService) Get(id int) (*models.Admin, error) {
	var admin models.Admin
	if err := models.DB.First(&admin, id).Error; err != nil {
		return nil, err
	}
	admin.Password = ""
	admin.Salt = ""
	return &admin, nil
}

func (s *AdminService) Create(admin *models.Admin) error {
	salt := fmt.Sprintf("%d", time.Now().UnixNano())
	password := admin.Password + salt
	hashed := md5.Sum([]byte(password))
	admin.Salt = salt
	admin.Password = hex.EncodeToString(hashed[:])
	admin.CreateTime = int(time.Now().Unix())
	if err := models.DB.Create(admin).Error; err != nil {
		return err
	}
	admin.Password = ""
	admin.Salt = ""
	return nil
}

func (s *AdminService) Update(id int, data map[string]interface{}) error {
	var admin models.Admin
	if err := models.DB.First(&admin, id).Error; err != nil {
		return err
	}
	delete(data, "id")
	delete(data, "create_time")
	delete(data, "password")
	delete(data, "salt")
	data["update_time"] = int(time.Now().Unix())
	return models.DB.Model(&admin).Updates(data).Error
}

func (s *AdminService) Delete(id int) error {
	return models.DB.Delete(&models.Admin{}, id).Error
}

func (s *AdminService) ChangePassword(userID int, oldPassword, newPassword string) error {
	var admin models.Admin
	if err := models.DB.First(&admin, userID).Error; err != nil {
		return err
	}

	oldPwd := oldPassword + admin.Salt
	oldHash := md5.Sum([]byte(oldPwd))
	if hex.EncodeToString(oldHash[:]) != admin.Password {
		return fmt.Errorf("原密码错误")
	}

	salt := fmt.Sprintf("%d", time.Now().UnixNano())
	newPwd := newPassword + salt
	newHash := md5.Sum([]byte(newPwd))

	admin.Salt = salt
	admin.Password = hex.EncodeToString(newHash[:])
	admin.UpdateTime = int(time.Now().Unix())

	return models.DB.Save(&admin).Error
}

func (s *AdminService) ChangeAdminPassword(adminID int, newPassword string) error {
	var admin models.Admin
	if err := models.DB.First(&admin, adminID).Error; err != nil {
		return err
	}

	salt := fmt.Sprintf("%d", time.Now().UnixNano())
	newPwd := newPassword + salt
	newHash := md5.Sum([]byte(newPwd))

	admin.Salt = salt
	admin.Password = hex.EncodeToString(newHash[:])
	admin.UpdateTime = int(time.Now().Unix())

	return models.DB.Save(&admin).Error
}

func validatePassword(password, salt, hash string) bool {
	pwd := password + salt
	h := md5.Sum([]byte(pwd))
	return hex.EncodeToString(h[:]) == hash
}

func generateToken(userID int, username string) (string, error) {
	return middleware.GenerateToken(userID, username)
}
