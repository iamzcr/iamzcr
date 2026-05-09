package frontend

import (
	"iamzcr/models"
)

type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

func (s *CategoryService) List() ([]models.Category, error) {
	var categories []models.Category
	err := models.DB.Where("status = ?", 1).Find(&categories).Error
	return categories, err
}
