package admin

import (
	"iamzcr/models"
	"time"
)

type MenuService struct{}

func NewMenuService() *MenuService {
	return &MenuService{}
}

func (s *MenuService) List(page, pageSize int) ([]models.Menu, int64, error) {
	var menus []models.Menu
	var total int64

	models.DB.Model(&models.Menu{}).Count(&total)
	err := models.DB.Order("weight DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&menus).Error

	return menus, total, err
}

func (s *MenuService) Get(id int) (*models.Menu, error) {
	var menu models.Menu
	if err := models.DB.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (s *MenuService) Create(menu *models.Menu) error {
	menu.CreateTime = int(time.Now().Unix())
	menu.UpdateTime = int(time.Now().Unix())
	return models.DB.Create(menu).Error
}

func (s *MenuService) Update(id int, input *models.Menu) (*models.Menu, error) {
	var menu models.Menu
	if err := models.DB.First(&menu, id).Error; err != nil {
		return nil, err
	}

	menu.Type = input.Type
	menu.Mark = input.Mark
	menu.Author = input.Author
	menu.Name = input.Name
	menu.Url = input.Url
	menu.Parent = input.Parent
	menu.Icon = input.Icon
	menu.Weight = input.Weight
	menu.Status = input.Status
	menu.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&menu).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (s *MenuService) Delete(id int) error {
	return models.DB.Delete(&models.Menu{}, id).Error
}
