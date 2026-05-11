## ADDED Requirements

### Requirement: Markdown to WeChat HTML conversion
The system SHALL provide a Markdown-to-WeChat-HTML converter that transforms article markdown content into HTML compatible with WeChat Official Account's draft API. The converter SHALL strip or adapt HTML elements unsupported by WeChat and SHALL resolve relative image URLs to absolute URLs using the `cdn_url` website setting.

#### Scenario: Basic Markdown conversion
- **WHEN** a Markdown article with headers, paragraphs, bold, italic, links, and images is converted
- **THEN** the output SHALL be valid HTML with `<h1>`-`<h6>`, `<p>`, `<strong>`, `<em>`, `<a>`, and `<img>` tags

#### Scenario: Code block adaptation
- **WHEN** a Markdown article contains a fenced code block
- **THEN** the output SHALL wrap the code in `<pre><code>` tags (WeChat does not support syntax highlighting)

#### Scenario: Relative image URL resolution
- **WHEN** a Markdown article contains an image with a relative path like `![alt](/upload/image.png)` and `cdn_url` is `https://example.com`
- **THEN** the output SHALL contain `<img src="https://example.com/upload/image.png">`

#### Scenario: Unsupported element removal
- **WHEN** a Markdown article contains HTML elements unsupported by WeChat (e.g., `<iframe>`, `<script>`, `<style>`)
- **THEN** those elements SHALL be removed or replaced with placeholder text in the output

### Requirement: WeChat draft publishing
The system SHALL publish article content to WeChat Official Account as drafts using PowerWeChat's Draft API (`draft/add`). The article title, author, digest (summary), content (converted HTML), and cover image SHALL be sent. On success, the returned `media_id` SHALL be stored in `sl_article_media`.

#### Scenario: Successful WeChat draft publish
- **WHEN** `WeChatService.PublishDraft(article)` is called with valid WeChat credentials and a complete article
- **THEN** the PowerWeChat Draft API SHALL be called, the WeChat `media_id` SHALL be extracted from the response, and an `ArticleMedia` record SHALL be created with `platform: "wechat"`, `status: 1`, `media_id: <from WeChat>`

#### Scenario: WeChat API error
- **WHEN** `WeChatService.PublishDraft(article)` is called and the WeChat API returns an error (e.g., invalid credential, rate limit)
- **THEN** the error SHALL be logged, an `ArticleMedia` record SHALL be created with `status: 2` and `error_msg` containing the WeChat error detail

#### Scenario: Article with no content
- **WHEN** `WeChatService.PublishDraft(article)` is called with an article that has empty content
- **THEN** the system SHALL return an error "文章内容为空" without calling the WeChat API

### Requirement: Publish triggered from article create/edit
When an article is created or updated through the admin API with a `publish_to_wechat: true` flag in the request body, the system SHALL save the article first, then SHALL attempt to publish it to WeChat. The save SHALL succeed even if the WeChat publish fails.

#### Scenario: Create article with publish flag
- **WHEN** `POST /api/articles` is called with `publish_to_wechat: true` and valid article data
- **THEN** the article SHALL be created, and then the article SHALL be published to WeChat

#### Scenario: Create article without publish flag
- **WHEN** `POST /api/articles` is called without `publish_to_wechat` or with `publish_to_wechat: false`
- **THEN** the article SHALL be created normally with no WeChat publish attempt

#### Scenario: Publish fails but article saved
- **WHEN** `POST /api/articles` is called with `publish_to_wechat: true` and the WeChat publish fails
- **THEN** the article SHALL still be created successfully with `code: 0`; the response data SHALL include a `wechat_publish_error` field

### Requirement: WeChatService as a reusable service
The WeChat service (`services/admin/wechat.go`) SHALL provide `PublishDraft(article *models.Article) (*models.ArticleMedia, error)` and a `NewWeChatService(websiteSettings map[string]string) *WeChatService` constructor. It SHALL hold the PowerWeChat client internally and SHALL NOT expose it.

#### Scenario: Construct WeChatService with settings
- **WHEN** `NewWeChatService(settings)` is called with a map containing `wechat_app_id` and `wechat_app_secret`
- **THEN** a `WeChatService` SHALL be returned with an initialized PowerWeChat Official Account client

#### Scenario: Construct WeChatService without required settings
- **WHEN** `NewWeChatService(settings)` is called with an empty or incomplete map
- **THEN** a `WeChatService` SHALL be returned with a nil client; calling `PublishDraft` SHALL return an error
