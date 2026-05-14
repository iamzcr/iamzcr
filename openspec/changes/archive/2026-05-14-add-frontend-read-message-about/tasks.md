## 1. Article Read Tracking (Backend)

- [x] 1.1 Add `read_count` to `GetByID` response in `services/frontend/article.go`
- [x] 1.2 Insert `sl_read` record in `GetArticle` handler in `handlers/frontend/handler.go`

## 2. Article Read Tracking (Frontend)

- [x] 2.1 Update `Article.vue` to display `read_count` instead of `article.weight`

## 3. Message API (Backend)

- [x] 3.1 Create `services/frontend/message.go` with `List()` and `Create()` methods
- [x] 3.2 Add `GetMessages` and `CreateMessage` handlers in `handlers/frontend/handler.go`
- [x] 3.3 Add `GET /api/messages` and `POST /api/messages` routes in `router/frontend.go`
- [x] 3.4 Add `MessageListResponse` to `handlers/frontend/response.go`

## 4. Message Page (Frontend)

- [x] 4.1 Add `messageApi` to `web/frontend/src/api/index.ts`
- [x] 4.2 Create `web/frontend/src/views/Messages.vue` with form and list
- [x] 4.3 Add `/messages` route in `web/frontend/src/main.ts`

## 5. About Me Page

- [x] 5.1 Create `web/frontend/src/views/About.vue` with centered content
- [x] 5.2 Add `/about` route in `web/frontend/src/main.ts`

## 6. Navigation Update

- [x] 6.1 Add "留言" and "关于我" links to App.vue nav

## 7. Docs

- [x] 7.1 Regenerate Swagger docs
