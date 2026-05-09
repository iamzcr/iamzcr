package admin

import (
	"iamzcr/models"
	"time"
)

type LangService struct{}

func NewLangService() *LangService {
	return &LangService{}
}

func (s *LangService) List(page, pageSize int) ([]models.Lang, int64, error) {
	var langs []models.Lang
	var total int64

	models.DB.Model(&models.Lang{}).Count(&total)
	err := models.DB.Order("weight DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&langs).Error

	return langs, total, err
}

func (s *LangService) Get(id int) (*models.Lang, error) {
	var lang models.Lang
	if err := models.DB.First(&lang, id).Error; err != nil {
		return nil, err
	}
	return &lang, nil
}

func (s *LangService) Create(lang *models.Lang) error {
	lang.CreateTime = int(time.Now().Unix())
	return models.DB.Create(lang).Error
}

func (s *LangService) Update(id int, data map[string]interface{}) error {
	var lang models.Lang
	if err := models.DB.First(&lang, id).Error; err != nil {
		return err
	}
	delete(data, "id")
	delete(data, "create_time")
	data["update_time"] = int(time.Now().Unix())
	return models.DB.Model(&lang).Updates(data).Error
}

func (s *LangService) Delete(id int) error {
	return models.DB.Delete(&models.Lang{}, id).Error
}
