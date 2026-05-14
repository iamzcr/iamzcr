package frontend

import (
	svc "iamzcr/services/frontend"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FrontendHandler struct {
	articleSvc   *svc.ArticleService
	categorySvc  *svc.CategoryService
	directorySvc *svc.DirectoryService
	tagsSvc      *svc.TagsService
	websiteSvc   *svc.WebsiteService
}

func NewFrontendHandler(
	articleSvc *svc.ArticleService,
	categorySvc *svc.CategoryService,
	directorySvc *svc.DirectoryService,
	tagsSvc *svc.TagsService,
	websiteSvc *svc.WebsiteService,
) *FrontendHandler {
	return &FrontendHandler{
		articleSvc:   articleSvc,
		categorySvc:  categorySvc,
		directorySvc: directorySvc,
		tagsSvc:      tagsSvc,
		websiteSvc:   websiteSvc,
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
//	@Description	获取网站配置信息（键值对形式，如站点名称、CDN地址等）
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
