package frontend

import (
	"iamzcr/models"
)

type ArticleService struct{}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

type ArticleWithTags struct {
	models.Article
	Tags []models.Tags `json:"tags"`
}

func (s *ArticleService) ListPublished(page, pageSize int, cid, did, tid string) ([]ArticleWithTags, int64) {
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
			return []ArticleWithTags{}, 0
		}
		query = query.Where("id IN ?", articleIDs)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("create_time DESC").Find(&articles)

	result := make([]ArticleWithTags, len(articles))
	for i, article := range articles {
		result[i].Article = article
		var articleTags []models.ArticleTags
		models.DB.Where("aid = ?", article.ID).Find(&articleTags)
		var tagIds []int
		for _, at := range articleTags {
			tagIds = append(tagIds, at.Tid)
		}
		if len(tagIds) > 0 {
			var tags []models.Tags
			models.DB.Where("id IN ?", tagIds).Find(&tags)
			result[i].Tags = tags
		}
	}

	return result, total
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
