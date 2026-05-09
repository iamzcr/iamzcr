package admin

import (
	"iamzcr/models"
	"time"
)

type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

func (s *CategoryService) List(page, pageSize int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	models.DB.Model(&models.Category{}).Count(&total)
	err := models.DB.Order("weight DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&categories).Error

	return categories, total, err
}

func (s *CategoryService) Get(id int) (*models.Category, error) {
	var category models.Category
	if err := models.DB.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) Create(category *models.Category) error {
	category.CreateTime = int(time.Now().Unix())
	category.UpdateTime = int(time.Now().Unix())
	return models.DB.Create(category).Error
}

func (s *CategoryService) Update(id int, input *models.Category) (*models.Category, error) {
	var category models.Category
	if err := models.DB.First(&category, id).Error; err != nil {
		return nil, err
	}

	category.Type = input.Type
	category.Parent = input.Parent
	category.Mark = input.Mark
	category.Author = input.Author
	category.Name = input.Name
	category.Weight = input.Weight
	category.Status = input.Status
	category.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) Delete(id int) error {
	return models.DB.Delete(&models.Category{}, id).Error
}
