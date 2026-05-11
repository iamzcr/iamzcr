## Why

The blog currently only publishes articles on the site itself. Many readers consume content through WeChat Official Accounts (微信公众号), and manually reformatting and republishing Markdown articles there is tedious and error-prone. Adding automated WeChat publishing expands the blog's reach and creates a foundation for multi-platform content distribution (B站, 小红书, etc.).

## What Changes

- Add WeChat Official Account configuration (AppID, AppSecret, etc.) as key-value settings in the existing `sl_website` table accessible through the admin "Website Settings" page
- Integrate `github.com/ArtisanCloud/PowerWeChat/v3` Go library for WeChat Official Account API calls
- Add a `sl_article_media` table to track per-article publishing status across multiple platforms (designed for extensibility)
- Add a "publish to third-party media" toggle in the article edit/create form
- Add a "publish" action button in the article list view to publish existing articles to WeChat
- Implement Markdown-to-WeChat-format conversion for article content
- Add backend API endpoints for article media publishing operations

## Capabilities

### New Capabilities
- `article-media-publishing`: Track and manage article publishing status across third-party media platforms. Covers the `sl_article_media` table, backend services for creating/querying media publish records, and API endpoints for publishing and status queries.
- `wechat-configuration`: Manage WeChat Official Account credentials and settings through the existing website settings key-value store. Covers the configuration keys, admin UI for setting them, and a backend WeChat API client initialized from those settings.
- `wechat-article-publish`: Publish articles to WeChat Official Account as drafts. Covers Markdown-to-WeChat conversion, the WeChat API integration via PowerWeChat, and the publish flow triggered from article create/edit and article list.

### Modified Capabilities
- `admin-content-services`: Article create and update handlers will accept an optional "publish to third-party media" flag and trigger platform publishing when set.

## Impact

- **New Go dependency**: `github.com/ArtisanCloud/PowerWeChat/v3`
- **New database table**: `sl_article_media` (article_id, platform, media_id, status, publish_url, create_time, update_time)
- **New Go files**: WeChat service, article media service, article media handler, markdown-to-WeChat converter
- **New frontend files**: WeChat publish action component in article list, media publish toggle in article edit form
- **Modified Go files**: `cmd/admin/main.go` (register new routes), `services/admin/article.go` (integrate publish logic), `handlers/admin/article.go` (add publish toggle to create/update), `models/models.go` (add ArticleMedia struct)
- **Modified frontend files**: `web/admin/src/views/ArticleEdit.vue` (add publish toggle), `web/admin/src/views/ArticleList.vue` (add publish action), `web/admin/src/views/WebsiteSettings.vue` (WeChat config keys), `web/admin/src/api/index.ts` (new API functions)
