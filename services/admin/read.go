package admin

import (
	"iamzcr/models"
	"time"
)

type ReadService struct{}

func NewReadService() *ReadService {
	return &ReadService{}
}

func (s *ReadService) List(page, pageSize int, aid string) ([]models.Read, int64, error) {
	var reads []models.Read
	var total int64

	query := models.DB.Model(&models.Read{})
	if aid != "" {
		query = query.Where("aid = ?", aid)
	}
	query.Count(&total)
	err := query.Order("id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&reads).Error

	return reads, total, err
}

func (s *ReadService) Get(id int) (*models.Read, error) {
	var read models.Read
	if err := models.DB.First(&read, id).Error; err != nil {
		return nil, err
	}
	return &read, nil
}

func (s *ReadService) Create(read *models.Read) error {
	read.CreateTime = int(time.Now().Unix())
	return models.DB.Create(read).Error
}

func (s *ReadService) Delete(id int) error {
	return models.DB.Delete(&models.Read{}, id).Error
}
