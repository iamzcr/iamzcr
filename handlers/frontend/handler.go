package frontend

import (
	svc "iamzcr/services/frontend"
	"iamzcr/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FrontendHandler struct {
	articleSvc   *svc.ArticleService
	categorySvc  *svc.CategoryService
	directorySvc *svc.DirectoryService
	tagsSvc      *svc.TagsService
	websiteSvc   *svc.WebsiteService
	messageSvc   *svc.MessageService
}

func NewFrontendHandler(
	articleSvc *svc.ArticleService,
	categorySvc *svc.CategoryService,
	directorySvc *svc.DirectoryService,
	tagsSvc *svc.TagsService,
	websiteSvc *svc.WebsiteService,
	messageSvc *svc.MessageService,
) *FrontendHandler {
	return &FrontendHandler{
		articleSvc:   articleSvc,
		categorySvc:  categorySvc,
		directorySvc: directorySvc,
		tagsSvc:      tagsSvc,
		websiteSvc:   websiteSvc,
		messageSvc:   messageSvc,
	}
}

// ListArticles godoc
//
//	@Summary		文章列表
//	@Description	获取已发布文章列表，支持分页和按分类/目录/标签过滤
//	@Tags			Frontend API
//	@Produce		json
//	@Param			page		query		int	false	"页码"	default(1)
//	@Param			page_size	query		int	false	"每页数量"	default(10)
//	@Param			cid			query		int	false	"分类ID"
//	@Param			did			query		int	false	"目录ID"
//	@Param			tid			query		int	false	"标签ID"
//	@Success		200			{object}	ArticleListResponse
//	@Router			/articles [get]
func (h *FrontendHandler) ListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	cid := c.DefaultQuery("cid", "")
	did := c.DefaultQuery("did", "")
	tid := c.DefaultQuery("tid", "")

	result, total := h.articleSvc.ListPublished(page, pageSize, cid, did, tid)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  result,
			"total": total,
		},
	})
}

// GetArticle godoc
//
//	@Summary		文章详情
//	@Description	获取单篇文章详情，包含分类、目录和标签信息
//	@Tags			Frontend API
//	@Produce		json
//	@Param			id	path		int	true	"文章ID"
//	@Success		200	{object}	ArticleDetailResponse
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/articles/{id} [get]
func (h *FrontendHandler) GetArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	article := h.articleSvc.GetByID(id)

	if article == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Article not found"})
		return
	}

	models.DB.Create(&models.Read{
		Aid:        id,
		Referer:    c.Request.Referer(),
		IP:         c.ClientIP(),
		CreateTime: int(time.Now().Unix()),
		UpdateTime: int(time.Now().Unix()),
	})

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    article,
	})
}

// GetCategories godoc
//
//	@Summary		分类列表
//	@Description	获取所有启用的文章分类
//	@Tags			Frontend API
//	@Produce		json
//	@Success		200	{object}	CategoryListResponse
//	@Router			/categories [get]
func (h *FrontendHandler) GetCategories(c *gin.Context) {
	categories, err := h.categorySvc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    categories,
	})
}

// GetDirectories godoc
//
//	@Summary		目录列表
//	@Description	获取所有启用的文章目录
//	@Tags			Frontend API
//	@Produce		json
//	@Success		200	{object}	DirectoryListResponse
//	@Router			/directories [get]
func (h *FrontendHandler) GetDirectories(c *gin.Context) {
	directories, err := h.directorySvc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    directories,
	})
}

// GetTags godoc
//
//	@Summary		标签列表
//	@Description	获取所有启用的文章标签
//	@Tags			Frontend API
//	@Produce		json
//	@Success		200	{object}	TagListResponse
//	@Router			/tags [get]
func (h *FrontendHandler) GetTags(c *gin.Context) {
	tags, err := h.tagsSvc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tags,
	})
}

// GetWebsite godoc
//
//	@Summary		网站设置
//	@Description	获取前端可见的网站配置信息（is_to_frontend=1 的键值对，如站点名称、CDN地址等）
//	@Tags			Frontend API
//	@Produce		json
//	@Success		200	{object}	WebsiteResponse
//	@Router			/website [get]
func (h *FrontendHandler) GetWebsite(c *gin.Context) {
	data, err := h.websiteSvc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

// GetMessages godoc
//
//	@Summary		留言列表
//	@Description	获取所有留言记录
//	@Tags			Frontend API
//	@Produce		json
//	@Success		200	{object}	MessageListResponse
//	@Router			/messages [get]
func (h *FrontendHandler) GetMessages(c *gin.Context) {
	messages, err := h.messageSvc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	if messages == nil {
		messages = []models.Message{}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    messages,
	})
}

// CreateMessage godoc
//
//	@Summary		提交留言
//	@Description	提交一条新留言，邮箱和内容为必填
//	@Tags			Frontend API
//	@Accept			json
//	@Produce		json
//	@Param			body	body		object	true	"留言内容 {name, email, url, content}"
//	@Success		200		{object}	MessageListResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Router			/messages [post]
func (h *FrontendHandler) CreateMessage(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求格式错误"})
		return
	}

	email, _ := input["email"].(string)
	content, _ := input["content"].(string)
	if email == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请输入邮箱和留言内容"})
		return
	}

	name, _ := input["name"].(string)
	url, _ := input["url"].(string)

	msg, err := h.messageSvc.Create(name, email, url, content, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    []models.Message{*msg},
	})
}
