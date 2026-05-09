package frontend

import (
	"iamzcr/models"
)

type TagsService struct{}

func NewTagsService() *TagsService {
	return &TagsService{}
}

func (s *TagsService) List() ([]models.Tags, error) {
	var tags []models.Tags
	err := models.DB.Where("status = ?", 1).Find(&tags).Error
	return tags, err
}
