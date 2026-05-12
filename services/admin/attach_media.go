package admin

import (
	"iamzcr/models"
	"time"
)

type AttachMediaService struct{}

func NewAttachMediaService() *AttachMediaService {
	return &AttachMediaService{}
}

func (s *AttachMediaService) List(page, pageSize int) ([]models.AttachMedia, int64, error) {
	var records []models.AttachMedia
	var total int64

	models.DB.Model(&models.AttachMedia{}).Count(&total)
	err := models.DB.Order("id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&records).Error

	return records, total, err
}

func (s *AttachMediaService) Get(id int) (*models.AttachMedia, error) {
	var record models.AttachMedia
	if err := models.DB.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *AttachMediaService) GetByAttachAndPlatform(attachID, platformID int) (*models.AttachMedia, error) {
	var record models.AttachMedia
	err := models.DB.Where("attach_id = ? AND platform_id = ?", attachID, platformID).First(&record).Error
	if err != nil {
		return nil, nil
	}
	return &record, nil
}

func (s *AttachMediaService) Create(record *models.AttachMedia) error {
	record.CreateTime = int(time.Now().Unix())
	record.UpdateTime = int(time.Now().Unix())
	return models.DB.Create(record).Error
}

func (s *AttachMediaService) Update(record *models.AttachMedia) error {
	record.UpdateTime = int(time.Now().Unix())
	return models.DB.Save(record).Error
}

func (s *AttachMediaService) Delete(id int) error {
	return models.DB.Delete(&models.AttachMedia{}, id).Error
}
