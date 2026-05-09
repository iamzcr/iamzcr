package frontend

import (
	"iamzcr/models"
)

type DirectoryService struct{}

func NewDirectoryService() *DirectoryService {
	return &DirectoryService{}
}

func (s *DirectoryService) List() ([]models.Directory, error) {
	var directories []models.Directory
	err := models.DB.Where("status = ?", 1).Find(&directories).Error
	return directories, err
}
