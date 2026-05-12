package admin

import (
	svc "iamzcr/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttachMediaHandler struct {
	attachMediaSvc *svc.AttachMediaService
	attachSvc      *svc.AttachService
	wechatSvc      *svc.WeChatService
}

func NewAttachMediaHandler(attachMediaSvc *svc.AttachMediaService, attachSvc *svc.AttachService, wechatSvc *svc.WeChatService) *AttachMediaHandler {
	return &AttachMediaHandler{
		attachMediaSvc: attachMediaSvc,
		attachSvc:      attachSvc,
		wechatSvc:      wechatSvc,
	}
}

func (h *AttachMediaHandler) List(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("page_size", "10"))

	records, total, err := h.attachMediaSvc.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  records,
			"total": total,
		},
	})
}

func (h *AttachMediaHandler) Get(c *gin.Context) {
	id := parseInt(c.Param("id"))

	record, err := h.attachMediaSvc.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    record,
	})
}

func (h *AttachMediaHandler) SyncToWechat(c *gin.Context) {
	id := parseInt(c.Param("id"))

	attach, err := h.attachSvc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "attachment not found"})
		return
	}

	mediaRecord, err := h.wechatSvc.UploadAttachMedia(attach)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	if mediaRecord != nil {
		existingRecord, _ := h.attachMediaSvc.GetByAttachAndPlatform(attach.ID, mediaRecord.PlatformID)
		if existingRecord != nil && existingRecord.ID > 0 {
			existingRecord.MediaID = mediaRecord.MediaID
			existingRecord.MediaURL = mediaRecord.MediaURL
			existingRecord.Status = mediaRecord.Status
			existingRecord.ErrorMsg = mediaRecord.ErrorMsg
			existingRecord.UpdateTime = mediaRecord.UpdateTime
			if updateErr := h.attachMediaSvc.Update(existingRecord); updateErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": updateErr.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "success",
				"data":    existingRecord,
			})
			return
		}

		if createErr := h.attachMediaSvc.Create(mediaRecord); createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": createErr.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    mediaRecord,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "sync failed"})
	}
}

func (h *AttachMediaHandler) Delete(c *gin.Context) {
	id := parseInt(c.Param("id"))

	if err := h.attachMediaSvc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
