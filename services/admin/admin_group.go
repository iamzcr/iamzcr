package admin

import (
	"iamzcr/models"
	"time"
)

type AdminGroupService struct{}

func NewAdminGroupService() *AdminGroupService {
	return &AdminGroupService{}
}

func (s *AdminGroupService) List() ([]models.AdminGroup, error) {
	var groups []models.AdminGroup
	err := models.DB.Order("id DESC").Find(&groups).Error
	return groups, err
}

func (s *AdminGroupService) Get(id int) (*models.AdminGroup, error) {
	var group models.AdminGroup
	if err := models.DB.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *AdminGroupService) Create(group *models.AdminGroup) error {
	group.CreateTime = int(time.Now().Unix())
	return models.DB.Create(group).Error
}

func (s *AdminGroupService) Update(id int, data map[string]interface{}) error {
	var group models.AdminGroup
	if err := models.DB.First(&group, id).Error; err != nil {
		return err
	}
	delete(data, "id")
	delete(data, "create_time")
	data["update_time"] = int(time.Now().Unix())
	return models.DB.Model(&group).Updates(data).Error
}

func (s *AdminGroupService) Delete(id int) error {
	return models.DB.Delete(&models.AdminGroup{}, id).Error
}
