package admin

import (
	"iamzcr/models"
	"time"
)

type PermitService struct{}

func NewPermitService() *PermitService {
	return &PermitService{}
}

func (s *PermitService) List(page, pageSize int) ([]models.Permit, int64, error) {
	var permits []models.Permit
	var total int64

	models.DB.Model(&models.Permit{}).Count(&total)
	err := models.DB.Order("weight DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&permits).Error

	return permits, total, err
}

func (s *PermitService) Get(id int) (*models.Permit, error) {
	var permit models.Permit
	if err := models.DB.First(&permit, id).Error; err != nil {
		return nil, err
	}
	return &permit, nil
}

func (s *PermitService) Create(permit *models.Permit) error {
	permit.CreateTime = int(time.Now().Unix())
	return models.DB.Create(permit).Error
}

func (s *PermitService) Update(id int, data map[string]interface{}) error {
	var permit models.Permit
	if err := models.DB.First(&permit, id).Error; err != nil {
		return err
	}
	delete(data, "id")
	delete(data, "create_time")
	data["update_time"] = int(time.Now().Unix())
	return models.DB.Model(&permit).Updates(data).Error
}

func (s *PermitService) Delete(id int) error {
	return models.DB.Delete(&models.Permit{}, id).Error
}
