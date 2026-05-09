package admin

import (
	"iamzcr/models"
	"time"
)

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

func (s *LogService) List(page, pageSize int, username string) ([]models.Log, int64, error) {
	var logs []models.Log
	var total int64

	query := models.DB.Model(&models.Log{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	query.Count(&total)
	err := query.Order("id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&logs).Error

	return logs, total, err
}

func (s *LogService) Get(id int) (*models.Log, error) {
	var log models.Log
	if err := models.DB.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (s *LogService) Create(log *models.Log) error {
	log.CreateTime = int(time.Now().Unix())
	return models.DB.Create(log).Error
}

func (s *LogService) Delete(id int) error {
	return models.DB.Delete(&models.Log{}, id).Error
}
