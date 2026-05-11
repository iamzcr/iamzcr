package admin

import (
	"fmt"
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleMediaHandler struct {
	articleMediaSvc *svc.ArticleMediaService
	articleSvc      *svc.ArticleService
	wechatSvc       *svc.WeChatService
}

func NewArticleMediaHandler(articleMediaSvc *svc.ArticleMediaService, articleSvc *svc.ArticleService, wechatSvc *svc.WeChatService) *ArticleMediaHandler {
	return &ArticleMediaHandler{
		articleMediaSvc: articleMediaSvc,
		articleSvc:      articleSvc,
		wechatSvc:       wechatSvc,
	}
}

func (h *ArticleMediaHandler) ListMedia(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	records, err := h.articleMediaSvc.ListByArticle(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	if records == nil {
		records = []models.ArticleMedia{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    records,
	})
}

func (h *ArticleMediaHandler) PublishToMedia(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PublishToMedia panic: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("服务器内部错误: %v", r)})
		}
	}()

	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Platforms []string `json:"platforms"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	article := h.articleSvc.Get(id)
	if article == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "文章不存在"})
		return
	}

	results := make([]interface{}, 0)
	publishError := ""

	for _, platform := range input.Platforms {
		switch platform {
		case "wechat":
			mediaRecord, err := h.wechatSvc.PublishDraft(article)
			if err != nil {
				publishError = err.Error()
			}
			if mediaRecord != nil {
				existingRecord, _ := h.articleMediaSvc.GetByArticleAndPlatform(id, "wechat")
				if existingRecord != nil && existingRecord.ID > 0 {
					existingRecord.MediaID = mediaRecord.MediaID
					existingRecord.Status = mediaRecord.Status
					existingRecord.ErrorMsg = mediaRecord.ErrorMsg
					existingRecord.UpdateTime = mediaRecord.UpdateTime
					if updateErr := h.articleMediaSvc.Update(existingRecord); updateErr != nil {
						publishError = updateErr.Error()
					}
					results = append(results, existingRecord)
				} else {
					if createErr := h.articleMediaSvc.Create(mediaRecord); createErr != nil {
						publishError = createErr.Error()
					} else {
						results = append(results, mediaRecord)
					}
				}
			}
		}
	}

	responseData := gin.H{
		"records": results,
	}
	if publishError != "" {
		responseData["error"] = publishError
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    responseData,
	})
}
