package admin

import (
	"iamzcr/models"
	"time"
)

type DirectoryService struct{}

func NewDirectoryService() *DirectoryService {
	return &DirectoryService{}
}

func (s *DirectoryService) List(page, pageSize int, cid string) ([]models.Directory, int64, error) {
	var dirs []models.Directory
	var total int64

	query := models.DB.Model(&models.Directory{})
	if cid != "" {
		query = query.Where("cid = ?", cid)
	}
	query.Count(&total)
	err := query.Order("weight DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&dirs).Error

	return dirs, total, err
}

func (s *DirectoryService) Get(id int) (*models.Directory, error) {
	var dir models.Directory
	if err := models.DB.First(&dir, id).Error; err != nil {
		return nil, err
	}
	return &dir, nil
}

func (s *DirectoryService) Create(dir *models.Directory) error {
	dir.CreateTime = int(time.Now().Unix())
	dir.UpdateTime = int(time.Now().Unix())
	return models.DB.Create(dir).Error
}

func (s *DirectoryService) Update(id int, input *models.Directory) (*models.Directory, error) {
	var dir models.Directory
	if err := models.DB.First(&dir, id).Error; err != nil {
		return nil, err
	}

	dir.Cid = input.Cid
	dir.Type = input.Type
	dir.Parent = input.Parent
	dir.Mark = input.Mark
	dir.Author = input.Author
	dir.Name = input.Name
	dir.Weight = input.Weight
	dir.Status = input.Status
	dir.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&dir).Error; err != nil {
		return nil, err
	}
	return &dir, nil
}

func (s *DirectoryService) Delete(id int) error {
	return models.DB.Delete(&models.Directory{}, id).Error
}
