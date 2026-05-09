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
