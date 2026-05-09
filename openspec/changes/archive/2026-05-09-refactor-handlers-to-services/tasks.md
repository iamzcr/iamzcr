## 1. Directory Setup

- [x] 1.1 Create `services/admin/` directory
- [x] 1.2 Create `services/frontend/` directory

## 2. Refactor Existing Article Service

- [x] 2.1 Create `services/admin/article.go` — move `ArticleService.List`, `Create`, `Update`, `Delete` from `services/article.go`; add `ListWithTags` method that includes tag enrichment logic (currently in `handlers/admin/article.go:ListArticles` lines 127–140); move `toInt` helper from `handlers/admin/handler.go` into this file
- [x] 2.2 Create `services/frontend/article.go` — move `ArticleService.ListPublished`, `GetByID` from `services/article.go`; add tag enrichment to `ListPublished` (currently in `handlers/frontend/handler.go:ListArticles` lines 31–50)
- [x] 2.3 Update `handlers/admin/article.go` — change `AdminHandler` to inject `*admin.ArticleService` instead of `*services.ArticleService`; remove inline tag enrichment from `ListArticles`, `CreateArticle`, `UpdateArticle`, `DeleteArticle` (delegate all to service); remove inline DB calls for ArticleTags in UpdateArticle and DeleteArticle; keep Login/Logout/GetAdminInfo in handler (these will move later)
- [x] 2.4 Update `handlers/frontend/handler.go` — change `FrontendHandler` to inject `*frontend.ArticleService` instead of `*services.ArticleService`; remove inline tag enrichment from `ListArticles`

## 3. Admin Content Services (Standard CRUD Pattern)

- [x] 3.1 Create `services/admin/category.go` — `CategoryService` with `List(page, pageSize)`, `Get(id)`, `Create(input)`, `Update(id, input)`, `Delete(id)`; move SQL from `handlers/admin/category.go`
- [x] 3.2 Update `handlers/admin/category.go` — inject `*admin.CategoryService`; delegate all DB operations to service; keep HTTP response formatting
- [x] 3.3 Create `services/admin/directory.go` — `DirectoryService` with standard CRUD methods; move SQL from `handlers/admin/directory.go`
- [x] 3.4 Update `handlers/admin/directory.go` — inject `*admin.DirectoryService`; delegate operations
- [x] 3.5 Create `services/admin/tags.go` — `TagsService` with standard CRUD methods; move SQL from `handlers/admin/tags.go`
- [x] 3.6 Update `handlers/admin/tags.go` — inject `*admin.TagsService`; delegate operations
- [x] 3.7 Create `services/admin/menu.go` — `MenuService` with standard CRUD methods; move SQL from `handlers/admin/menu.go`
- [x] 3.8 Update `handlers/admin/menu.go` — inject `*admin.MenuService`; delegate operations
- [x] 3.9 Create `services/admin/comment.go` — `CommentService` with standard CRUD methods; move SQL from `handlers/admin/comment.go`
- [x] 3.10 Update `handlers/admin/comment.go` — inject `*admin.CommentService`; delegate operations

## 4. Admin System Services

- [x] 4.1 Create `services/admin/admin.go` — `AdminService` with `Login(username, password, clientIP)`, `List()`, `Get(id)`, `Create(admin)`, `Update(id, admin)`, `Delete(id)`, `ChangePassword(userID, oldPwd, newPwd)`, `ChangeAdminPassword(adminID, newPwd)`; move `validatePassword`, MD5+salt hashing, `generateTestToken` from `handlers/admin/handler.go`; move all SQL from `handlers/admin/admin.go`
- [x] 4.2 Update `handlers/admin/admin.go` and `handlers/admin/article.go` — inject `*admin.AdminService` into `AdminHandler`; delegate Login/GetAdminInfo/ListAdmins/GetAdmin/CreateAdmin/UpdateAdmin/DeleteAdmin/ChangePassword/ChangeAdminPassword to service; keep HTTP response formatting in handler
- [x] 4.3 Create `services/admin/admin_group.go` — `AdminGroupService` with standard CRUD methods; move SQL from `handlers/admin/admin.go:AdminGroupHandler`
- [x] 4.4 Update `handlers/admin/admin.go:AdminGroupHandler` — inject `*admin.AdminGroupService`; delegate operations
- [x] 4.5 Create `services/admin/attach.go` — `AttachService` with `List`, `Get`, `Create`, `Update`, `Delete`, `Upload(fileHeader, assetDir)`; move SQL from `handlers/admin/attach.go`; move file save logic from `Upload` handler
- [x] 4.6 Update `handlers/admin/attach.go` — inject `*admin.AttachService`; delegate operations including file upload
- [x] 4.7 Create `services/admin/website.go` — `WebsiteService` with `List()`, `Get(id)`, `Upsert(key, value)`, `Delete(id)`; move SQL from `handlers/admin/website.go`
- [x] 4.8 Update `handlers/admin/website.go` — inject `*admin.WebsiteService`; delegate operations
- [x] 4.9 Create `services/admin/lang.go` — `LangService` with standard CRUD; move SQL from `handlers/admin/lang.go`
- [x] 4.10 Update `handlers/admin/lang.go` — inject `*admin.LangService`; delegate operations
- [x] 4.11 Create `services/admin/log.go` — `LogService` with `List`, `Get`, `Create`, `Delete`; move SQL from `handlers/admin/log.go`
- [x] 4.12 Update `handlers/admin/log.go` — inject `*admin.LogService`; delegate operations
- [x] 4.13 Create `services/admin/message.go` — `MessageService` with standard CRUD; move SQL from `handlers/admin/message.go`
- [x] 4.14 Update `handlers/admin/message.go` — inject `*admin.MessageService`; delegate operations
- [x] 4.15 Create `services/admin/permit.go` — `PermitService` with standard CRUD; move SQL from `handlers/admin/permit.go`
- [x] 4.16 Update `handlers/admin/permit.go` — inject `*admin.PermitService`; delegate operations
- [x] 4.17 Create `services/admin/read.go` — `ReadService` with `List`, `Get`, `Create`, `Delete`; move SQL from `handlers/admin/read.go`
- [x] 4.18 Update `handlers/admin/read.go` — inject `*admin.ReadService`; delegate operations

## 5. Frontend Services

- [x] 5.1 Create `services/frontend/category.go` — `CategoryService` with `List()` returning categories where status=1; move SQL from `handlers/frontend/handler.go:GetCategories`
- [x] 5.2 Create `services/frontend/directory.go` — `DirectoryService` with `List()` returning directories where status=1; move SQL from `handlers/frontend/handler.go:GetDirectories`
- [x] 5.3 Create `services/frontend/tags.go` — `TagsService` with `List()` returning tags where status=1; move SQL from `handlers/frontend/handler.go:GetTags`
- [x] 5.4 Create `services/frontend/website.go` — `WebsiteService` with `List()` returning `map[string]string`; move SQL from `handlers/frontend/handler.go:GetWebsite`
- [x] 5.5 Update `handlers/frontend/handler.go` — inject all 5 frontend services; delegate all remaining methods to services; remove all direct `models.DB` calls

## 6. Wire Up Services in main.go

- [x] 6.1 Update `cmd/admin/main.go` — create all admin service instances and pass them into handler constructors (AdminHandler gets articleSvc + adminSvc; other handlers each get their respective service)
- [x] 6.2 Update `cmd/frontend/main.go` — create all frontend service instances and pass them into `FrontendHandler` constructor

## 7. Cleanup and Verify

- [x] 7.1 Delete `services/article.go` (all logic migrated to `services/admin/article.go` and `services/frontend/article.go`)
- [x] 7.2 Remove `toInt`, `generateTestToken`, `validatePassword` from `handlers/admin/handler.go` (migrated to services)
- [x] 7.3 Run `go build ./cmd/admin/...` and `go build ./cmd/frontend/...` to verify compilation
- [ ] 7.4 Run the project and smoke-test: admin login, admin article list/create/edit/delete, frontend article list/detail, frontend categories/directories/tags/website
