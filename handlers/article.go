package handlers

import (
	"blog/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct{}

func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

func (h *ArticleHandler) List(c *gin.Context) {
	var articles []models.Article
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	var total int64
	models.DB.Model(&models.Article{}).Where("status = ?", 1).Count(&total)

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	offset := (pageInt - 1) * pageSizeInt
	models.DB.Where("status = ?", 1).Offset(offset).Limit(pageSizeInt).Order("create_time DESC").Find(&articles)

	type ArticleWithTags struct {
		models.Article
		Tags []models.Tags `json:"tags"`
	}

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

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  result,
			"total": total,
		},
	})
}

func (h *ArticleHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	if err := models.DB.Where("status = ?", 1).First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Article not found"})
		return
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

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"article":   article,
			"category":  category,
			"directory": directory,
			"tags":      tags,
		},
	})
}

func (h *ArticleHandler) Create(c *gin.Context) {
	var input struct {
		Cid        int    `json:"cid"`
		Did        int    `json:"did"`
		Title      string `json:"title"`
		Desc       string `json:"desc"`
		Keyword    string `json:"keyword"`
		Author     string `json:"author"`
		Thumb      string `json:"thumb"`
		Summary    string `json:"summary"`
		Content    string `json:"content"`
		IsHot      int    `json:"is_hot"`
		IsNew      int    `json:"is_new"`
		IsRecom    int    `json:"is_recom"`
		Weight     int    `json:"weight"`
		PublicTime int    `json:"public_time"`
		Status     int    `json:"status"`
		Month      string `json:"month"`
		TagIds     []int  `json:"tag_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	article := models.Article{
		Cid:        input.Cid,
		Did:        input.Did,
		Title:      input.Title,
		Desc:       input.Desc,
		Keyword:    input.Keyword,
		Author:     input.Author,
		Thumb:      input.Thumb,
		Summary:    input.Summary,
		Content:    input.Content,
		IsHot:      input.IsHot,
		IsNew:      input.IsNew,
		IsRecom:    input.IsRecom,
		Weight:     input.Weight,
		PublicTime: input.PublicTime,
		Status:     input.Status,
		Month:      input.Month,
		CreateTime: int(time.Now().Unix()),
		UpdateTime: int(time.Now().Unix()),
	}

	if err := models.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	for _, tid := range input.TagIds {
		at := models.ArticleTags{
			Aid:        article.ID,
			Tid:        tid,
			CreateTime: int(time.Now().Unix()),
			UpdateTime: int(time.Now().Unix()),
		}
		models.DB.Create(&at)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    article,
	})
}

func (h *ArticleHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	if err := models.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Article not found"})
		return
	}

	var input struct {
		Cid        int    `json:"cid"`
		Did        int    `json:"did"`
		Title      string `json:"title"`
		Desc       string `json:"desc"`
		Keyword    string `json:"keyword"`
		Author     string `json:"author"`
		Thumb      string `json:"thumb"`
		Summary    string `json:"summary"`
		Content    string `json:"content"`
		IsHot      int    `json:"is_hot"`
		IsNew      int    `json:"is_new"`
		IsRecom    int    `json:"is_recom"`
		Weight     int    `json:"weight"`
		PublicTime int    `json:"public_time"`
		Status     int    `json:"status"`
		Month      string `json:"month"`
		TagIds     []int  `json:"tag_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	article.Cid = input.Cid
	article.Did = input.Did
	article.Title = input.Title
	article.Desc = input.Desc
	article.Keyword = input.Keyword
	article.Author = input.Author
	article.Thumb = input.Thumb
	article.Summary = input.Summary
	article.Content = input.Content
	article.IsHot = input.IsHot
	article.IsNew = input.IsNew
	article.IsRecom = input.IsRecom
	article.Weight = input.Weight
	article.PublicTime = input.PublicTime
	article.Status = input.Status
	article.Month = input.Month
	article.UpdateTime = int(time.Now().Unix())

	if err := models.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	models.DB.Where("aid = ?", article.ID).Delete(&models.ArticleTags{})
	for _, tid := range input.TagIds {
		at := models.ArticleTags{
			Aid:        article.ID,
			Tid:        tid,
			CreateTime: int(time.Now().Unix()),
			UpdateTime: int(time.Now().Unix()),
		}
		models.DB.Create(&at)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    article,
	})
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	models.DB.Where("aid = ?", id).Delete(&models.ArticleTags{})
	if err := models.DB.Delete(&models.Article{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
