package admin

import (
	"iamzcr/models"
	"time"
)

type PlatformService struct{}

func NewPlatformService() *PlatformService {
	return &PlatformService{}
}

func (s *PlatformService) List(page, pageSize int) ([]models.Platform, int64, error) {
	var platforms []models.Platform
	var total int64

	models.DB.Model(&models.Platform{}).Count(&total)
	err := models.DB.Order("id ASC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&platforms).Error

	return platforms, total, err
}

func (s *PlatformService) Get(id int) (*models.Platform, error) {
	var platform models.Platform
	if err := models.DB.First(&platform, id).Error; err != nil {
		return nil, err
	}
	return &platform, nil
}

func (s *PlatformService) Create(platform *models.Platform) error {
	platform.CreateTime = int(time.Now().Unix())
	platform.UpdateTime = int(time.Now().Unix())
	return models.DB.Create(platform).Error
}

func (s *PlatformService) Update(id int, input *models.Platform) (*models.Platform, error) {
	var platform models.Platform
	if err := models.DB.First(&platform, id).Error; err != nil {
		return nil, err
	}

	platform.Mark = input.Mark
	platform.Name = input.Name
	platform.Status = input.Status
	platform.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&platform).Error; err != nil {
		return nil, err
	}
	return &platform, nil
}

func (s *PlatformService) Delete(id int) error {
	return models.DB.Delete(&models.Platform{}, id).Error
}

func (s *PlatformService) GetByMark(mark string) (*models.Platform, error) {
	var platform models.Platform
	if err := models.DB.Where("mark = ?", mark).First(&platform).Error; err != nil {
		return nil, err
	}
	return &platform, nil
}
