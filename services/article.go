package services

import (
	"fmt"
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

func (s *ArticleService) ListPublished(page, pageSize int, cid, did, tid string) ([]models.Article, int64) {
	var articles []models.Article
	var total int64

	query := models.DB.Model(&models.Article{}).Where("status = ?", 1)
	if cid != "" {
		query = query.Where("cid = ?", cid)
	}
	if did != "" {
		query = query.Where("did = ?", did)
	}
	if tid != "" {
		var articleIDs []int
		models.DB.Model(&models.ArticleTags{}).Where("tid = ?", tid).Pluck("aid", &articleIDs)
		if len(articleIDs) == 0 {
			return []models.Article{}, 0
		}
		query = query.Where("id IN ?", articleIDs)
	}

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

func toInt(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case int64:
		return int(val)
	case string:
		var i int
		fmt.Sscanf(val, "%d", &i)
		return i
	default:
		return 0
	}
}

func (s *ArticleService) Create(data map[string]interface{}) *models.Article {
	article := models.Article{
		Cid:        toInt(data["cid"]),
		Did:        toInt(data["did"]),
		Title:      data["title"].(string),
		Desc:       data["desc"].(string),
		Keyword:    data["keyword"].(string),
		Author:     data["author"].(string),
		Thumb:      data["thumb"].(string),
		Summary:    data["summary"].(string),
		Content:    data["content"].(string),
		IsHot:      toInt(data["is_hot"]),
		IsNew:      toInt(data["is_new"]),
		IsRecom:    toInt(data["is_recom"]),
		Weight:     toInt(data["weight"]),
		PublicTime: toInt(data["public_time"]),
		Status:     toInt(data["status"]),
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

	if v, ok := data["cid"]; ok {
		article.Cid = toInt(v)
	}
	if v, ok := data["did"]; ok {
		article.Did = toInt(v)
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
	if v, ok := data["is_hot"]; ok {
		article.IsHot = toInt(v)
	}
	if v, ok := data["is_new"]; ok {
		article.IsNew = toInt(v)
	}
	if v, ok := data["is_recom"]; ok {
		article.IsRecom = toInt(v)
	}
	if v, ok := data["weight"]; ok {
		article.Weight = toInt(v)
	}
	if v, ok := data["public_time"]; ok {
		article.PublicTime = toInt(v)
	}
	if v, ok := data["status"]; ok {
		article.Status = toInt(v)
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
