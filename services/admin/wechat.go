package admin

import (
	"context"
	"errors"
	"fmt"
	"iamzcr/models"
	"iamzcr/pkg/md2wechat"
	"path/filepath"
	"sync"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	publishRequest "github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/publish/request"
)

type WeChatService struct {
	mu   sync.Mutex
	app  *officialAccount.OfficialAccount
}

func NewWeChatService() *WeChatService {
	return &WeChatService{}
}

func (s *WeChatService) readSettings() (appID, appSecret, cdnURL string, err error) {
	var websites []models.Website
	if err := models.DB.Find(&websites).Error; err != nil {
		return "", "", "", errors.New("unable to read website config")
	}
	settings := make(map[string]string)
	for _, w := range websites {
		settings[w.Key] = w.Value
	}
	appID = settings["wechat_app_id"]
	appSecret = settings["wechat_app_secret"]
	if appID == "" || appSecret == "" {
		return "", "", "", errors.New("wechat config incomplete")
	}
	return appID, appSecret, settings["cdn_url"], nil
}

func (s *WeChatService) getApp() (*officialAccount.OfficialAccount, error) {
	if s.app != nil {
		return s.app, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.app != nil {
		return s.app, nil
	}

	appID, appSecret, _, err := s.readSettings()
	if err != nil {
		return nil, err
	}

	app, err := officialAccount.NewOfficialAccount(&officialAccount.UserConfig{
		AppID:  appID,
		Secret: appSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("init wechat client failed: %v", err)
	}

	s.app = app
	return s.app, nil
}

func (s *WeChatService) getPlatformID(mark string) (int, error) {
	var platform models.Platform
	if err := models.DB.Where("mark = ?", mark).First(&platform).Error; err != nil {
		return 0, fmt.Errorf("platform %s not found", mark)
	}
	return platform.ID, nil
}

func (s *WeChatService) PublishDraft(article *models.Article) (*models.ArticleMedia, error) {
	_, _, cdnURL, err := s.readSettings()
	if err != nil {
		return nil, err
	}

	if article.Content == "" {
		return nil, errors.New("article content is empty")
	}

	htmlContent, err := md2wechat.Convert(article.Content, cdnURL)
	if err != nil {
		return nil, err
	}

	app, err := s.getApp()
	if err != nil {
		return nil, err
	}

	platformID, err := s.getPlatformID("wechat")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	draftReq := &publishRequest.RequestDraftAdd{
		Articles: []*publishRequest.Article{
			{
				Title:              article.Title,
				Author:             article.Author,
				Digest:             truncate(article.Desc, 120),
				Content:            htmlContent,
				ContentSourceUrl:   "",
				NeedOpenComment:    0,
				OnlyFansCanComment: 0,
			},
		},
	}

	result, err := app.Publish.DraftAdd(ctx, draftReq)
	if err != nil {
		return &models.ArticleMedia{
			Aid:        article.ID,
			PlatformID: platformID,
			Status:     2,
			ErrorMsg:   err.Error(),
			CreateTime: int(time.Now().Unix()),
			UpdateTime: int(time.Now().Unix()),
		}, err
	}

	if result.MediaID == "" {
		errMsg := fmt.Sprintf("publish failed: errcode=%d errmsg=%s", result.ErrCode, result.ErrMsg)
		return &models.ArticleMedia{
			Aid:        article.ID,
			PlatformID: platformID,
			Status:     2,
			ErrorMsg:   errMsg,
			CreateTime: int(time.Now().Unix()),
			UpdateTime: int(time.Now().Unix()),
		}, errors.New(errMsg)
	}

	return &models.ArticleMedia{
		Aid:        article.ID,
		PlatformID: platformID,
		MediaID:    result.MediaID,
		Status:     1,
		CreateTime: int(time.Now().Unix()),
		UpdateTime: int(time.Now().Unix()),
	}, nil
}

func (s *WeChatService) UploadAttachMedia(attach *models.Attach) (*models.AttachMedia, error) {
	app, err := s.getApp()
	if err != nil {
		return nil, err
	}

	platformID, err := s.getPlatformID("wechat")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	attachPath := attach.Path
	ext := filepath.Ext(attachPath)

	now := int(time.Now().Unix())
	makeFailed := func(errMsg string, origErr error) (*models.AttachMedia, error) {
		return &models.AttachMedia{
			AttachID:   attach.ID,
			PlatformID: platformID,
			Status:     2,
			ErrorMsg:   errMsg,
			CreateTime: now,
			UpdateTime: now,
		}, origErr
	}

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		result, err := app.Material.UploadImage(ctx, attachPath)
		if err != nil {
			return makeFailed(err.Error(), err)
		}
		if result.MediaID == "" && result.URL == "" {
			return makeFailed(fmt.Sprintf("upload failed: errcode=%d errmsg=%s", result.ErrCode, result.ErrMsg), errors.New(result.ErrMsg))
		}
		return &models.AttachMedia{
			AttachID:   attach.ID,
			PlatformID: platformID,
			MediaID:    result.MediaID,
			MediaURL:   result.URL,
			Status:     1,
			CreateTime: now,
			UpdateTime: now,
		}, nil

	case ".mp4":
		result, err := app.Material.UploadVideo(ctx, attachPath, attach.Name, "")
		if err != nil {
			return makeFailed(err.Error(), err)
		}
		if result.MediaID == "" && result.URL == "" {
			return makeFailed(fmt.Sprintf("upload failed: errcode=%d errmsg=%s", result.ErrCode, result.ErrMsg), errors.New(result.ErrMsg))
		}
		return &models.AttachMedia{
			AttachID:   attach.ID,
			PlatformID: platformID,
			MediaID:    result.MediaID,
			MediaURL:   result.URL,
			Status:     1,
			CreateTime: now,
			UpdateTime: now,
		}, nil

	case ".mp3", ".wav", ".amr":
		result, err := app.Material.UploadVoice(ctx, attachPath)
		if err != nil {
			return makeFailed(err.Error(), err)
		}
		if result.MediaID == "" && result.URL == "" {
			return makeFailed(fmt.Sprintf("upload failed: errcode=%d errmsg=%s", result.ErrCode, result.ErrMsg), errors.New(result.ErrMsg))
		}
		return &models.AttachMedia{
			AttachID:   attach.ID,
			PlatformID: platformID,
			MediaID:    result.MediaID,
			MediaURL:   result.URL,
			Status:     1,
			CreateTime: now,
			UpdateTime: now,
		}, nil

	default:
		result, err := app.Material.UploadArticleImage(ctx, attachPath)
		if err != nil {
			return makeFailed(err.Error(), err)
		}
		if result.MediaID == "" && result.URL == "" {
			return makeFailed(fmt.Sprintf("upload failed: errcode=%d errmsg=%s", result.ErrCode, result.ErrMsg), errors.New(result.ErrMsg))
		}
		return &models.AttachMedia{
			AttachID:   attach.ID,
			PlatformID: platformID,
			MediaID:    result.MediaID,
			MediaURL:   result.URL,
			Status:     1,
			CreateTime: now,
			UpdateTime: now,
		}, nil
	}
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}
