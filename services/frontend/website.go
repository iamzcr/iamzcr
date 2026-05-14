package frontend

import (
	"iamzcr/models"
)

type WebsiteService struct{}

func NewWebsiteService() *WebsiteService {
	return &WebsiteService{}
}

func (s *WebsiteService) List() (map[string]string, error) {
	var websites []models.Website
	if err := models.DB.Where("is_to_frontend = ?", 1).Find(&websites).Error; err != nil {
		return nil, err
	}
	data := make(map[string]string)
	for _, w := range websites {
		data[w.Key] = w.Value
	}
	return data, nil
}
