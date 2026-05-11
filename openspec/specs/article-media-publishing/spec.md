## ADDED Requirements

### Requirement: Article media status tracking
The system SHALL maintain a `sl_article_media` table with one record per article per platform to track publishing status across third-party media platforms. Each record SHALL include `article_id`, `platform` (a string identifier like `wechat`, `bilibili`, `xiaohongshu`), `media_id` (the external platform's article/draft ID), `status` (0=pending, 1=published, 2=failed), `publish_url` (optional URL on the platform), `error_msg` (error detail if status=failed), `create_time`, and `update_time`.

#### Scenario: Query media status for an article
- **WHEN** `GET /api/articles/:id/media` is called
- **THEN** the response SHALL return a list of media publish records for that article with `code: 0`

#### Scenario: Article has no media records
- **WHEN** `GET /api/articles/:id/media` is called for an article with no media records
- **THEN** the response SHALL return an empty list with `code: 0`

#### Scenario: Create media record on publish
- **WHEN** an article is published to a platform
- **THEN** a new `sl_article_media` record SHALL be created with `status: 1` and the platform's media ID

#### Scenario: Record publish failure
- **WHEN** an article publish attempt fails
- **THEN** a `sl_article_media` record SHALL be created or updated with `status: 2` and `error_msg` containing the failure reason

### Requirement: Article media CRUD service
The article media service (`services/admin/article_media.go`) SHALL provide `ListByArticle(articleID)`, `GetByArticleAndPlatform(articleID, platform)`, `Create(record)`, and `Update(record)` methods. It SHALL NOT provide `Delete` (records are immutable audit trail).

#### Scenario: List by article
- **WHEN** `ArticleMediaService.ListByArticle(42)` is called
- **THEN** it returns all `[]ArticleMedia` records where `article_id = 42`

#### Scenario: Upsert behavior
- **WHEN** `ArticleMediaService.GetByArticleAndPlatform(42, "wechat")` is called and a record exists
- **THEN** it returns the existing record; if none exists, it returns nil with no error

### Requirement: Media publish API endpoint
The admin API SHALL expose `POST /api/articles/:id/media/publish` that accepts a JSON body with `platforms` (array of platform names, e.g. `["wechat"]`) and triggers publishing to each specified platform. The endpoint requires JWT authentication.

#### Scenario: Publish to WeChat
- **WHEN** `POST /api/articles/:id/media/publish` is called with `{"platforms": ["wechat"]}`
- **THEN** the system SHALL attempt to publish the article to WeChat and return the media record in the response

#### Scenario: Publish to multiple platforms
- **WHEN** `POST /api/articles/:id/media/publish` is called with `{"platforms": ["wechat", "bilibili"]}`
- **THEN** the system SHALL attempt publishing to each platform independently and return all media records; failure on one platform SHALL NOT prevent others

#### Scenario: Article not found
- **WHEN** `POST /api/articles/:id/media/publish` is called for a non-existent article ID
- **THEN** the response SHALL return `code: 404` with message "文章不存在"
