package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iamzcr/models"
	"iamzcr/pkg/md2wechat"
	"io"
	"net/http"
	"time"
)

type WeChatService struct{}

func NewWeChatService() *WeChatService {
	return &WeChatService{}
}

func (s *WeChatService) readSettings() (appID, appSecret, cdnURL string, err error) {
	var websites []models.Website
	if err := models.DB.Find(&websites).Error; err != nil {
		return "", "", "", errors.New("无法读取网站配置")
	}
	settings := make(map[string]string)
	for _, w := range websites {
		settings[w.Key] = w.Value
	}
	appID = settings["wechat_app_id"]
	appSecret = settings["wechat_app_secret"]
	if appID == "" || appSecret == "" {
		return "", "", "", errors.New("微信公众号配置不完整，请先在基础设置中配置 AppID 和 AppSecret")
	}
	return appID, appSecret, settings["cdn_url"], nil
}

func (s *WeChatService) getAccessToken(appID, appSecret string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appID, appSecret)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("获取access_token失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		AccessToken string `json:"access_token"`
		Errcode     int    `json:"errcode"`
		Errmsg      string `json:"errmsg"`
	}
	json.Unmarshal(body, &result)
	if result.AccessToken == "" {
		return "", fmt.Errorf("获取access_token失败: %s", string(body))
	}
	return result.AccessToken, nil
}

func (s *WeChatService) PublishDraft(article *models.Article) (*models.ArticleMedia, error) {
	appID, appSecret, cdnURL, err := s.readSettings()
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

	token, err := s.getAccessToken(appID, appSecret)
	if err != nil {
		return nil, err
	}

	draft := map[string]interface{}{
		"articles": []map[string]interface{}{
			{
				"title":              article.Title,
				"author":             article.Author,
				"digest":             truncate(article.Desc, 120),
				"content":            htmlContent,
				"content_source_url": "",
				"need_open_comment":  0,
				"only_fans_can_comment": 0,
			},
		},
	}

	body, _ := json.Marshal(draft)
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/draft/add?access_token=%s", token)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
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
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		MediaID string `json:"media_id"`
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	json.Unmarshal(respBody, &result)

	if result.MediaID == "" {
		errMsg := fmt.Sprintf("发布失败: %s", string(respBody))
		return &models.ArticleMedia{
			Aid:        article.ID,
			Platform:   "wechat",
			Status:     2,
			ErrorMsg:   errMsg,
			CreateTime: int(time.Now().Unix()),
			UpdateTime: int(time.Now().Unix()),
		}, errors.New(errMsg)
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
