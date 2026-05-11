package admin

import (
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	articleSvc *svc.ArticleService
	adminSvc   *svc.AdminService
	wechatSvc  *svc.WeChatService
}

func NewAdminHandler(articleSvc *svc.ArticleService, adminSvc *svc.AdminService, wechatSvc *svc.WeChatService) *AdminHandler {
	return &AdminHandler{
		articleSvc: articleSvc,
		adminSvc:   adminSvc,
		wechatSvc:  wechatSvc,
	}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	result, err := h.adminSvc.Login(input.Username, input.Password, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    result,
	})
}

func (h *AdminHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (h *AdminHandler) GetAdminInfo(c *gin.Context) {
	userID := c.GetInt("user_id")

	admin, adminGroup, err := h.adminSvc.GetAdminInfo(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"name":     admin.Name,
			"group":    adminGroup.Name,
		},
	})
}

func (h *AdminHandler) ListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, total := h.articleSvc.ListWithTags(page, pageSize)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  result,
			"total": total,
		},
	})
}

func (h *AdminHandler) GetArticle(c *gin.Context) {
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

func (h *AdminHandler) CreateArticle(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var tagIDs []int
	if rawTagIDs, ok := input["tag_ids"].([]interface{}); ok {
		for _, tid := range rawTagIDs {
			if id, ok := tid.(float64); ok {
				tagIDs = append(tagIDs, int(id))
			}
		}
	}

	publishToWechat, _ := input["publish_to_wechat"].(bool)

	article := h.articleSvc.Create(input, tagIDs)

	responseData := gin.H{"article": article}

	if publishToWechat && h.wechatSvc != nil {
		mediaRecord, err := h.wechatSvc.PublishDraft(article)
		if err != nil {
			responseData["wechat_publish_error"] = err.Error()
		}
		if mediaRecord != nil {
			models.DB.Create(mediaRecord)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    responseData,
	})
}

func (h *AdminHandler) UpdateArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var tagIDs []int
	if rawTagIDs, ok := input["tag_ids"].([]interface{}); ok {
		for _, tid := range rawTagIDs {
			if id, ok := tid.(float64); ok {
				tagIDs = append(tagIDs, int(id))
			}
		}
	}

	publishToWechat, _ := input["publish_to_wechat"].(bool)

	article := h.articleSvc.Update(id, input, tagIDs)

	if article == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Article not found"})
		return
	}

	responseData := gin.H{"article": article}

	if publishToWechat && h.wechatSvc != nil {
		mediaRecord, err := h.wechatSvc.PublishDraft(article)
		if err != nil {
			responseData["wechat_publish_error"] = err.Error()
		}
		if mediaRecord != nil {
			var existing models.ArticleMedia
			result := models.DB.Where("aid = ? AND platform = ?", id, "wechat").First(&existing)
			if result.Error == nil && existing.ID > 0 {
				existing.MediaID = mediaRecord.MediaID
				existing.Status = mediaRecord.Status
				existing.ErrorMsg = mediaRecord.ErrorMsg
				existing.UpdateTime = mediaRecord.UpdateTime
				models.DB.Save(&existing)
			} else {
				models.DB.Create(mediaRecord)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    responseData,
	})
}

func (h *AdminHandler) DeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	success := h.articleSvc.Delete(id)

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
