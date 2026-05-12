package admin

import (
	"iamzcr/models"
)

type ArticleMediaService struct{}

func NewArticleMediaService() *ArticleMediaService {
	return &ArticleMediaService{}
}

func (s *ArticleMediaService) ListByArticle(articleID int) ([]models.ArticleMedia, error) {
	var records []models.ArticleMedia
	err := models.DB.Where("aid = ?", articleID).Order("id DESC").Find(&records).Error
	return records, err
}

func (s *ArticleMediaService) GetByArticleAndPlatform(articleID, platformID int) (*models.ArticleMedia, error) {
	var record models.ArticleMedia
	err := models.DB.Where("aid = ? AND platform_id = ?", articleID, platformID).First(&record).Error
	if err != nil {
		return nil, nil
	}
	return &record, nil
}

func (s *ArticleMediaService) Create(record *models.ArticleMedia) error {
	return models.DB.Create(record).Error
}

func (s *ArticleMediaService) Update(record *models.ArticleMedia) error {
	return models.DB.Save(record).Error
}
