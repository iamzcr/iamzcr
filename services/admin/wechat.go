package admin

import (
	"context"
	"errors"
	"iamzcr/models"
	"iamzcr/pkg/md2wechat"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	publishReq "github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/publish/request"
)

type WeChatService struct{}

func NewWeChatService() *WeChatService {
	return &WeChatService{}
}

func (s *WeChatService) getClient() (*officialAccount.OfficialAccount, string, error) {
	var websites []models.Website
	if err := models.DB.Find(&websites).Error; err != nil {
		return nil, "", errors.New("无法读取网站配置")
	}

	settings := make(map[string]string)
	for _, w := range websites {
		settings[w.Key] = w.Value
	}

	appID := settings["wechat_app_id"]
	appSecret := settings["wechat_app_secret"]

	if appID == "" || appSecret == "" {
		return nil, "", errors.New("微信公众号配置不完整，请先在基础设置中配置 AppID 和 AppSecret")
	}

	app, err := officialAccount.NewOfficialAccount(&officialAccount.UserConfig{
		AppID:  appID,
		Secret: appSecret,
	})
	if err != nil {
		return nil, "", errors.New("微信公众号初始化失败: " + err.Error())
	}

	return app, settings["cdn_url"], nil
}

func (s *WeChatService) PublishDraft(article *models.Article) (*models.ArticleMedia, error) {
	app, cdnURL, err := s.getClient()
	if err != nil {
		return nil, err
	}

	if article.Content == "" {
		return nil, errors.New("文章内容为空")
	}

	htmlContent, err := md2wechat.Convert(article.Content, cdnURL)
	if err != nil {
		return nil, err
	}

	req := &publishReq.RequestDraftAdd{
		Articles: []*publishReq.Article{
			{
				Title:   article.Title,
				Author:  article.Author,
				Digest:  truncate(article.Desc, 120),
				Content: htmlContent,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := app.Publish.DraftAdd(ctx, req)
	if err != nil {
		return &models.ArticleMedia{
			Aid:        article.ID,
			Platform:   "wechat",
			Status:     2,
			ErrorMsg:   err.Error(),
			CreateTime: int(time.Now().Unix()),
			UpdateTime: int(time.Now().Unix()),
		}, err
	}

	if result == nil {
		return nil, errors.New("微信公众号返回空响应")
	}

	return &models.ArticleMedia{
		Aid:        article.ID,
		Platform:   "wechat",
		MediaID:    result.MediaID,
		Status:     1,
		CreateTime: int(time.Now().Unix()),
		UpdateTime: int(time.Now().Unix()),
	}, nil
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}
