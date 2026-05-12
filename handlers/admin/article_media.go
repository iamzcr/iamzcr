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
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("server error: %v", r)})
		}
	}()

	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		PlatformIDs []int `json:"platform_ids"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	article := h.articleSvc.Get(id)
	if article == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "article not found"})
		return
	}

	results := make([]interface{}, 0)
	publishError := ""

	for _, platformID := range input.PlatformIDs {
		var platform models.Platform
		if err := models.DB.First(&platform, platformID).Error; err != nil {
			publishError = fmt.Sprintf("platform %d not found", platformID)
			continue
		}

		switch platform.Mark {
		case "wechat":
			mediaRecord, err := h.wechatSvc.PublishDraft(article)
			if err != nil {
				publishError = err.Error()
			}
			if mediaRecord != nil {
				existingRecord, _ := h.articleMediaSvc.GetByArticleAndPlatform(id, platformID)
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
