package frontend

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
	err := models.DB.Order("create_time DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&messages).Error
	return messages, total, err
}

func (s *MessageService) Create(name, email, url, content, ip string) (*models.Message, error) {
	msg := models.Message{
		Name:       name,
		Email:      email,
		URL:        url,
		Content:    content,
		IP:         ip,
		IsReply:    0,
		CreateTime: int(time.Now().Unix()),
		UpdateTime: int(time.Now().Unix()),
	}
	if err := models.DB.Create(&msg).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}
