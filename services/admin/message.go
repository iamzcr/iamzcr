package admin

import (
	"iamzcr/models"
	"time"
)

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) List(page, pageSize int) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	models.DB.Model(&models.Message{}).Count(&total)
	err := models.DB.Order("id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&messages).Error

	return messages, total, err
}

func (s *MessageService) Get(id int) (*models.Message, error) {
	var message models.Message
	if err := models.DB.First(&message, id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (s *MessageService) Create(message *models.Message) error {
	message.CreateTime = int(time.Now().Unix())
	return models.DB.Create(message).Error
}

func (s *MessageService) Update(id int, data map[string]interface{}) error {
	var message models.Message
	if err := models.DB.First(&message, id).Error; err != nil {
		return err
	}
	delete(data, "id")
	delete(data, "create_time")
	data["update_time"] = int(time.Now().Unix())
	return models.DB.Model(&message).Updates(data).Error
}

func (s *MessageService) Delete(id int) error {
	return models.DB.Delete(&models.Message{}, id).Error
}
