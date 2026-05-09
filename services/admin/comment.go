package admin

import (
	"iamzcr/models"
	"time"
)

type CommentService struct{}

func NewCommentService() *CommentService {
	return &CommentService{}
}

func (s *CommentService) List(page, pageSize int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	models.DB.Model(&models.Comment{}).Count(&total)
	err := models.DB.Order("create_time DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&comments).Error

	return comments, total, err
}

func (s *CommentService) Get(id int) (*models.Comment, error) {
	var comment models.Comment
	if err := models.DB.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *CommentService) Create(comment *models.Comment) error {
	comment.CreateTime = int(time.Now().Unix())
	comment.UpdateTime = int(time.Now().Unix())
	return models.DB.Create(comment).Error
}

func (s *CommentService) Update(id int, content string) (*models.Comment, error) {
	var comment models.Comment
	if err := models.DB.First(&comment, id).Error; err != nil {
		return nil, err
	}

	comment.Content = content
	comment.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *CommentService) Delete(id int) error {
	return models.DB.Delete(&models.Comment{}, id).Error
}
