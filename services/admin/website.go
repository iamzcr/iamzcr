package admin

import (
	"iamzcr/models"
	"time"
)

type WebsiteService struct{}

func NewWebsiteService() *WebsiteService {
	return &WebsiteService{}
}

func (s *WebsiteService) List() ([]models.Website, error) {
	var websites []models.Website
	err := models.DB.Order("id DESC").Find(&websites).Error
	return websites, err
}

func (s *WebsiteService) Get() (map[string]interface{}, error) {
	var websites []models.Website
	if err := models.DB.Find(&websites).Error; err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	for _, w := range websites {
		data[w.Key] = w.Value
	}
	return data, nil
}

type WebsiteInput struct {
	Key          string `json:"key"`
	Value        string `json:"value"`
	IsToFrontend int    `json:"is_to_frontend"`
}

func (s *WebsiteService) Upsert(input map[string]string) error {
	for key, value := range input {
		if err := s.upsertOne(key, value, 0); err != nil {
			return err
		}
	}
	return nil
}

func (s *WebsiteService) UpsertAll(items []WebsiteInput) error {
	for _, item := range items {
		if err := s.upsertOne(item.Key, item.Value, item.IsToFrontend); err != nil {
			return err
		}
	}
	return nil
}

func (s *WebsiteService) upsertOne(key, value string, isToFrontend int) error {
	var website models.Website
	if err := models.DB.Where("`key` = ?", key).First(&website).Error; err == nil {
		website.Value = value
		website.IsToFrontend = isToFrontend
		website.UpdateTime = int(time.Now().Unix())
		models.DB.Save(&website)
	} else {
		if isToFrontend == 0 {
			isToFrontend = 1
		}
		website := models.Website{
			Key:          key,
			Value:        value,
			Staus:        1,
			IsToFrontend: isToFrontend,
			CreateTime:   int(time.Now().Unix()),
			UpdateTime:   int(time.Now().Unix()),
		}
		models.DB.Create(&website)
	}
	return nil
}

func (s *WebsiteService) Delete(id int) error {
	return models.DB.Delete(&models.Website{}, id).Error
}
