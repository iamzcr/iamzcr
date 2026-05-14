package frontend

import (
	"iamzcr/models"
	"time"
)

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) List() ([]models.Message, error) {
	var messages []models.Message
	err := models.DB.Order("create_time DESC").Find(&messages).Error
	return messages, err
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
