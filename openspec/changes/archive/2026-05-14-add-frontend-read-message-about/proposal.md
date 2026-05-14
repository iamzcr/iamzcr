## Why

The blog SPA currently lacks visitor engagement features. Article read counts display the `weight` field (an SEO tuning value) instead of actual read data. There is no way for visitors to leave messages, and no personal "About Me" page. These are standard expectations for a personal blog.

## What Changes

- Record article reads in `sl_read` table when a visitor views an article detail page, capturing IP, referer, and timestamp
- Display real read count from `sl_read` COUNT in article detail page (replacing `article.weight`)
- Add `/messages` page with a guestbook form (name, email, url, content) and message list — email and content are required
- Add `/about` static page with WeChat QR code, bio, and links (GitHub, blog)
- Add "留言" and "关于我" to the top navigation bar
- Form submission uses button-disable debounce to prevent duplicate submissions

## Capabilities

### New Capabilities

- `article-read-tracking`: Record article read events in `sl_read` table, return accurate read count with article detail
- `frontend-messages`: Visitors can view and submit messages via a guestbook page; email and content are required fields
- `about-page`: Static "About Me" page displaying WeChat QR code, personal bio, and social links

### Modified Capabilities

- `frontend-services`: Article `GetByID` now returns `read_count` in response data; new `MessageService` with `List` and `Create` methods added

## Impact

- **New API**: `GET /api/messages` (list), `POST /api/messages` (submit), read count added to `GET /api/articles/:id` response
- **New Vue pages**: `Messages.vue`, `About.vue`
- **Modified Vue**: `App.vue` (nav links), `Article.vue` (read count display), `main.ts` (new routes)
- **Backend**: New `services/frontend/message.go`, modified `services/frontend/article.go`, modified `handlers/frontend/handler.go`, modified `router/frontend.go`
- **Breaking**: None — message submission is anonymous (no auth), read tracking is transparent
