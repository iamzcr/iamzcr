package admin

import (
	"iamzcr/models"
	"time"
)

type TagsService struct{}

func NewTagsService() *TagsService {
	return &TagsService{}
}

func (s *TagsService) List(page, pageSize int) ([]models.Tags, int64, error) {
	var tags []models.Tags
	var total int64

	models.DB.Model(&models.Tags{}).Count(&total)
	err := models.DB.Order("weight DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&tags).Error

	return tags, total, err
}

func (s *TagsService) Get(id int) (*models.Tags, error) {
	var tag models.Tags
	if err := models.DB.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *TagsService) Create(tag *models.Tags) error {
	tag.CreateTime = int(time.Now().Unix())
	tag.UpdateTime = int(time.Now().Unix())
	return models.DB.Create(tag).Error
}

func (s *TagsService) Update(id int, input *models.Tags) (*models.Tags, error) {
	var tag models.Tags
	if err := models.DB.First(&tag, id).Error; err != nil {
		return nil, err
	}

	tag.Type = input.Type
	tag.Mark = input.Mark
	tag.Author = input.Author
	tag.Name = input.Name
	tag.Weight = input.Weight
	tag.Status = input.Status
	tag.IsHot = input.IsHot
	tag.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *TagsService) Delete(id int) error {
	return models.DB.Delete(&models.Tags{}, id).Error
}
