## Context

The blog system ("堆栈人生") is a Go + Vue 3 application with admin and frontend APIs. Articles are written in Markdown and stored in `sl_article.content`. Currently there is no mechanism to publish articles to external platforms; content only lives on the blog site itself. The admin UI has an article editor (v-md-editor for Markdown) and an article list with actions (edit/delete). Settings are managed through a key-value `sl_website` table with an admin UI page (WebsiteSettings.vue).

## Goals / Non-Goals

**Goals:**
- Add WeChat Official Account (微信公众号) article publishing capability
- Store WeChat API credentials as website settings (reusing existing `sl_website` key-value table)
- Allow publishing during article creation and as a separate action from the article list
- Convert Markdown article content to WeChat-compatible HTML format
- Design the `sl_article_media` table to support future platform additions (B站, 小红书, etc.)
- Integrate PowerWeChat Go SDK for WeChat API communication

**Non-Goals:**
- Publishing to B站 or 小红书 (only foundational table design for them)
- WeChat authentication / OAuth flows (we manage credentials in settings)
- Frontend rendering of articles in WeChat format (backend conversion only)
- Scheduling or batch publishing
- WeChat article edit/delete sync (publish-only; draft management in WeChat backend)

## Decisions

### 1. PowerWeChat v3 as WeChat SDK
**Choice:** Use `github.com/ArtisanCloud/PowerWeChat/v3`  
**Rationale:** It is the most actively maintained Go WeChat SDK with a dedicated Official Account module that supports draft publishing (draft/add, draft/get). It handles access token management internally.  
**Alternatives considered:**
- `silenceper/wechat` – older, less maintained, more complex API surface
- Direct HTTP calls to WeChat API – more code to write and maintain (token refresh, error handling)

### 2. WeChat config in sl_website (key-value)
**Choice:** Store WeChat credentials (`wechat_app_id`, `wechat_app_secret`, `wechat_token`, `wechat_aes_key`, `wechat_original_id`) as rows in the existing `sl_website` table.  
**Rationale:** Avoids a new config table; reuse existing WebsiteSettings admin UI for edits; simple to add contextual help/placeholder text in the frontend. Consistency with how other system configs are stored.  
**Alternatives considered:**
- New `sl_wechat_config` table – overkill for 5 key-value pairs; would need new handler, service, and UI
- Environment variables – not manageable through admin UI

### 3. sl_article_media table for per-article platform tracking
**Choice:** New table `sl_article_media(article_id, platform, media_id, status, publish_url, error_msg, create_time, update_time)`. One row per article per platform.  
**Rationale:** Decouples article storage from publishing state. Adding a new platform means inserting new rows with a different `platform` value – no schema changes. Querying "has this article been published to WeChat?" is a simple lookup.  
**Alternatives considered:**
- Add columns to `sl_article` (e.g., `wechat_media_id`, `wechat_status`) – not extensible; each platform adds 2+ columns
- JSON column on `sl_article` – hard to query, no referential clarity

### 4. Markdown → WeChat HTML conversion
**Choice:** Use `github.com/yuin/goldmark` (already used indirectly via v-md-editor concepts) to convert Markdown to HTML, then apply WeChat-specific sanitization (strip unsupported tags, inline CSS for WeChat's limited HTML support, handle images with absolute URLs).  
**Rationale:** goldmark is the Go standard for Markdown → HTML. WeChat supports a limited HTML subset; a post-processing pass strips/handles unsupported elements. Image src URLs must be absolute (prepend CDN URL from website settings).  
**Alternatives considered:**
- blackfriday – v2 unmaintained
- Custom regex-based conversion – fragile, doesn't handle nested structures

### 5. Publishing flow triggers
**Choice:** Two entry points:
1. **Article create/edit form** – a toggle "发布到微信公众号" on the form. If checked on save, after the article is saved, attempt WeChat publish.
2. **Article list** – a "发布" action button (or "已发布" badge if already published). Click triggers the publish API.
**Rationale:** Cover both workflows: publishing a new article immediately, and publishing an existing/draft article later.  
**Alternatives considered:**
- Only article list action – misses the create-and-publish workflow
- Auto-publish on save – too aggressive; user should opt in

### 6. API endpoint design (Convention A)
**Choice:** New endpoints under `/api/articles/:id/media` (resource-based URL, consistent with existing GET/PUT/DELETE pattern):
- `GET /api/articles/:id/media` – list media publish status for an article
- `POST /api/articles/:id/media/publish` – publish article to specified platform(s)
**Rationale:** Convention A (`code: 0, message: "success"`) matches most existing handlers. Resource nesting under articles keeps routes organized.  
**Alternatives considered:**
- `/api/media/article/:id/publish` – less natural, creates a new top-level resource
- PUT for publish – POST is more idiomatic for "trigger an action"

### 7. New Go package organization
**Choice:**
- `services/admin/article_media.go` – CRUD for sl_article_media records
- `services/admin/wechat.go` – WeChat API client wrapper (holds PowerWeChat instance, methods for draft publish)
- `handlers/admin/article_media.go` – HTTP handlers for media status and publish
- `pkg/md2wechat/` – Markdown to WeChat HTML converter
- `models/models.go` – add `ArticleMedia` struct
**Rationale:** Follow existing codebase pattern (service per resource, handler per resource). Converter as a reusable package under `pkg/`.

## Risks / Trade-offs

- **WeChat API rate limiting**: WeChat limits draft publishing frequency. → Mitigation: Return clear error to user, record error_msg in sl_article_media, allow re-publish.
- **Access token expiry**: PowerWeChat handles this internally, but if the credentials are wrong, publish will fail. → Mitigation: Validate credentials on save in settings (future enhancement; initial implementation shows error on first publish attempt).
- **Markdown → WeChat fidelity loss**: WeChat supports only a limited HTML subset (no code fences with syntax highlights, limited table support, no custom CSS). → Mitigation: Convert as best-effort; document limitations; code blocks rendered as pre>code with monospace.
- **Image URL resolution**: Article images in Markdown may use relative paths. → Mitigation: Prepend `cdn_url` from website settings to relative image paths during conversion.
- **Non-idempotent publish**: Publishing the same article twice to WeChat creates a new draft each time. → Mitigation: Check sl_article_media for existing record and confirm overwrite in UI.

## Open Questions

- Should we support WeChat cover image (thumb_media_id) upload? The article.thumb field could map to it, but requires separate image upload to WeChat first.
- Should article publish status in WeChat (draft vs published) be synced back? PowerWeChat supports draft status queries.
