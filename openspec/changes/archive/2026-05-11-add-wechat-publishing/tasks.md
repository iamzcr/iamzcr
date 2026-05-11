## 1. Database & Models

- [x] 1.1 Add `sl_article_media` table SQL to `sql/stacklifes.sql` (columns: id, article_id, platform, media_id, status, publish_url, error_msg, create_time, update_time)
- [x] 1.2 Add `ArticleMedia` struct to `models/models.go` with GORM tags matching `sl_article_media`
- [x] 1.3 Add `MD` package-level import to `go.mod` if not already present

## 2. Markdown to WeChat Converter

- [x] 2.1 Create `pkg/md2wechat/md2wechat.go` with `Convert(markdown string, cdnURL string) (string, error)` function
- [x] 2.2 Implement goldmark-based Markdown to HTML conversion
- [x] 2.3 Implement post-processing: strip unsupported tags (`<iframe>`, `<script>`, `<style>`), resolve relative image URLs to absolute using cdnURL
- [x] 2.4 Handle code blocks: wrap in `<pre><code>` without syntax highlighting

## 3. WeChat Service

- [x] 3.1 Run `go get github.com/ArtisanCloud/PowerWeChat/v3` to add the dependency
- [x] 3.2 Create `services/admin/wechat.go` with `WeChatService` struct holding PowerWeChat Official Account client
- [x] 3.3 Implement `NewWeChatService(websiteSettings map[string]string) *WeChatService` constructor
- [x] 3.4 Implement `PublishDraft(article *models.Article) (*models.ArticleMedia, error)` method using PowerWeChat Draft API
- [x] 3.5 Implement error handling: return descriptive error when credentials missing or API call fails

## 4. Article Media Service

- [x] 4.1 Create `services/admin/article_media.go` with `ArticleMediaService` struct
- [x] 4.2 Implement `ListByArticle(articleID int) ([]models.ArticleMedia, error)`
- [x] 4.3 Implement `GetByArticleAndPlatform(articleID int, platform string) (*models.ArticleMedia, error)`
- [x] 4.4 Implement `Create(record *models.ArticleMedia) error`
- [x] 4.5 Implement `Update(record *models.ArticleMedia) error` (for status/error_msg updates)

## 5. Article Media API Handler

- [x] 5.1 Create `handlers/admin/article_media.go` with `ArticleMediaHandler` struct
- [x] 5.2 Implement `ListMedia` – `GET /api/articles/:id/media` – returns media publish records for an article
- [x] 5.3 Implement `PublishToMedia` – `POST /api/articles/:id/media/publish` – accepts `{"platforms": ["wechat"]}`, triggers publish for each platform
- [x] 5.4 Use Convention A response format (`code: 0, message: "success", data: ...`)

## 6. Route Registration

- [x] 6.1 Initialize `ArticleMediaService` and `WeChatService` in `cmd/admin/main.go`
- [x] 6.2 Register new routes: `GET /api/articles/:id/media`, `POST /api/articles/:id/media/publish` (JWT-protected)

## 7. Modify Article Create/Update for WeChat Publish Toggle

- [x] 7.1 Modify `services/admin/article.go` `Create()` to accept and process `publish_to_wechat` flag – if true, call WeChatService.PublishDraft after save
- [x] 7.2 Modify `services/admin/article.go` `Update()` to accept and process `publish_to_wechat` flag – if true, call WeChatService.PublishDraft after save
- [x] 7.3 Ensure article save succeeds even if WeChat publish fails; include `wechat_publish_error` in response when applicable
- [x] 7.4 Update `handlers/admin/article.go` – pass `publish_to_wechat` from request body to service call

## 8. Frontend – API Layer

- [x] 8.1 Add `getArticleMedia(articleId: number)` API function to `web/admin/src/api/index.ts`
- [x] 8.2 Add `publishToMedia(articleId: number, platforms: string[])` API function to `web/admin/src/api/index.ts`

## 9. Frontend – Article Edit (WeChat Publish Toggle)

- [x] 9.1 Add `publish_to_wechat` toggle switch (NSwitch) to `web/admin/src/views/ArticleEdit.vue` form
- [x] 9.2 Include `publish_to_wechat` in the data sent to `articleApi.create()` / `articleApi.update()`

## 10. Frontend – Article List (Publish Action)

- [x] 10.1 Add "发布" action button in `web/admin/src/views/ArticleList.vue` article row actions
- [x] 10.2 On click, call `publishToMedia()` API and show success/error message via `useMessage()`
- [x] 10.3 Show publish status badge (e.g., "已发布" green tag if already published to WeChat)
- [x] 10.4 Fetch article media status for list items to determine badge visibility

## 11. Frontend – Website Settings (WeChat Config)

- [x] 11.1 Add hint/guide text in `web/admin/src/views/WebsiteSettings.vue` listing the WeChat config key names (`wechat_app_id`, `wechat_app_secret`, etc.)
- [x] 11.2 Add input placeholder text for each WeChat config key to help admin configure correctly

## 12. Verification

- [x] 12.1 Run `go build ./...` to verify Go compilation
- [x] 12.2 Run `cd web/admin && npm run build` to verify frontend type-check and build
- [ ] 12.3 Manual test: configure WeChat settings, create article with publish toggle, verify record in sl_article_media
- [ ] 12.4 Manual test: publish existing article from article list, verify record in sl_article_media
