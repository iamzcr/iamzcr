package services

import (
	"iamzcr/models"
	"time"
)

type ArticleService struct{}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

func (s *ArticleService) List(page, pageSize int) ([]models.Article, int64) {
	var articles []models.Article
	var total int64

	query := models.DB.Model(&models.Article{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("create_time DESC").Find(&articles)

	return articles, total
}

func (s *ArticleService) Get(id int) *models.Article {
	var article models.Article
	if err := models.DB.First(&article, id).Error; err != nil {
		return nil
	}
	return &article
}

func (s *ArticleService) Create(data map[string]interface{}) *models.Article {
	article := models.Article{
		Cid:        data["cid"].(int),
		Did:        data["did"].(int),
		Title:      data["title"].(string),
		Desc:       data["desc"].(string),
		Keyword:    data["keyword"].(string),
		Author:     data["author"].(string),
		Thumb:      data["thumb"].(string),
		Summary:    data["summary"].(string),
		Content:    data["content"].(string),
		IsHot:      data["is_hot"].(int),
		IsNew:      data["is_new"].(int),
		IsRecom:    data["is_recom"].(int),
		Weight:     data["weight"].(int),
		PublicTime: data["public_time"].(int),
		Status:     data["status"].(int),
		Month:      data["month"].(string),
		CreateTime: int(time.Now().Unix()),
		UpdateTime: int(time.Now().Unix()),
	}

	models.DB.Create(&article)
	return &article
}

func (s *ArticleService) Update(id int, data map[string]interface{}) *models.Article {
	var article models.Article
	if err := models.DB.First(&article, id).Error; err != nil {
		return nil
	}

	if v, ok := data["cid"].(int); ok {
		article.Cid = v
	}
	if v, ok := data["did"].(int); ok {
		article.Did = v
	}
	if v, ok := data["title"].(string); ok {
		article.Title = v
	}
	if v, ok := data["desc"].(string); ok {
		article.Desc = v
	}
	if v, ok := data["keyword"].(string); ok {
		article.Keyword = v
	}
	if v, ok := data["author"].(string); ok {
		article.Author = v
	}
	if v, ok := data["thumb"].(string); ok {
		article.Thumb = v
	}
	if v, ok := data["summary"].(string); ok {
		article.Summary = v
	}
	if v, ok := data["content"].(string); ok {
		article.Content = v
	}
	if v, ok := data["is_hot"].(int); ok {
		article.IsHot = v
	}
	if v, ok := data["is_new"].(int); ok {
		article.IsNew = v
	}
	if v, ok := data["is_recom"].(int); ok {
		article.IsRecom = v
	}
	if v, ok := data["weight"].(int); ok {
		article.Weight = v
	}
	if v, ok := data["public_time"].(int); ok {
		article.PublicTime = v
	}
	if v, ok := data["status"].(int); ok {
		article.Status = v
	}
	if v, ok := data["month"].(string); ok {
		article.Month = v
	}

	article.UpdateTime = int(time.Now().Unix())
	models.DB.Save(&article)
	return &article
}

func (s *ArticleService) Delete(id int) bool {
	result := models.DB.Delete(&models.Article{}, id)
	return result.RowsAffected > 0
}

func (s *ArticleService) GetByID(id int) map[string]interface{} {
	var article models.Article
	if err := models.DB.First(&article, id).Error; err != nil {
		return nil
	}

	var category models.Category
	if article.Cid > 0 {
		models.DB.First(&category, article.Cid)
	}

	var directory models.Directory
	if article.Did > 0 {
		models.DB.First(&directory, article.Did)
	}

	var articleTags []models.ArticleTags
	models.DB.Where("aid = ?", article.ID).Find(&articleTags)

	var tagIds []int
	for _, at := range articleTags {
		tagIds = append(tagIds, at.Tid)
	}

	var tags []models.Tags
	if len(tagIds) > 0 {
		models.DB.Where("id IN ?", tagIds).Find(&tags)
	}

	return map[string]interface{}{
		"article":   article,
		"category":  category,
		"directory": directory,
		"tags":      tags,
	}
}
