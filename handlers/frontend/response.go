package frontend

import "iamzcr/models"

type ArticleWithTags struct {
	models.Article
	Tags []models.Tags `json:"tags"`
}

type ArticleListData struct {
	List  []ArticleWithTags `json:"list"`
	Total int64             `json:"total"`
}

type ArticleListResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    ArticleListData `json:"data"`
}

type ArticleDetailData struct {
	Article     models.Article   `json:"article"`
	Category    models.Category  `json:"category"`
	Directory   models.Directory `json:"directory"`
	Tags        []models.Tags    `json:"tags"`
	ReadCount   int64            `json:"read_count"`
	PrevArticle *models.Article  `json:"prev_article"`
	NextArticle *models.Article  `json:"next_article"`
}

type ArticleDetailResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    ArticleDetailData `json:"data"`
}

type CategoryListResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    []models.Category `json:"data"`
}

type DirectoryListResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    []models.Directory   `json:"data"`
}

type TagListResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    []models.Tags `json:"data"`
}

type WebsiteData map[string]string

type WebsiteResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    WebsiteData `json:"data"`
}

type MessageListData struct {
	List  []models.Message `json:"list"`
	Total int64            `json:"total"`
}

type MessageListResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    MessageListData `json:"data"`
}
