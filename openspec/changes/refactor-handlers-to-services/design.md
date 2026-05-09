## Context

The project "堆棧人生" (Stack Life) is a Go + Gin + GORM blog system with two API servers (Admin on 8081, Frontend on 8082). Currently, 16 handler files contain all business logic inline — SQL queries, data transformations, password hashing, file upload logic, and N+1 tag enrichment all live at the HTTP layer. Only one service file exists (`services/article.go` with `ArticleService`), which is used by both admin and frontend handlers but still incomplete (tag enrichment remains in handlers).

The target architecture is:
- **Handlers**: Parse HTTP requests, call services, format HTTP responses
- **Services**: All business logic — DB queries, data transformations, auth logic, file operations
- **Models**: Data structures (already well-separated)

The service layer will mirror the handler admin/frontend split:
```
services/
├── admin/
│   ├── article.go      (ArticleService)
│   ├── category.go     (CategoryService)
│   ├── comment.go      (CommentService)
│   ├── directory.go    (DirectoryService)
│   ├── tags.go         (TagsService)
│   ├── menu.go         (MenuService)
│   ├── website.go      (WebsiteService)
│   ├── admin.go        (AdminService)
│   ├── admin_group.go  (AdminGroupService)
│   ├── attach.go       (AttachService)
│   ├── lang.go         (LangService)
│   ├── log.go          (LogService)
│   ├── message.go      (MessageService)
│   ├── permit.go       (PermitService)
│   ├── read.go         (ReadService)
├── frontend/
│   ├── article.go      (ArticleService)
│   ├── category.go     (CategoryService)
│   ├── directory.go    (DirectoryService)
│   ├── tags.go         (TagsService)
│   ├── website.go      (WebsiteService)
```

The existing `services/article.go` will be deleted after migration.

## Goals / Non-Goals

**Goals:**
- Every handler method delegates business logic to a service method; handlers only handle HTTP concerns
- Services organized by domain boundary: `services/admin/` for admin write/read operations, `services/frontend/` for public read-only operations
- Tag-enrichment N+1 logic moved into article services (not duplicated in handlers)
- Auth/login logic, password hashing/validation extracted into admin user service
- File upload logic extracted into attach service
- Existing response conventions preserved per entity (Convention A: `code:0` / Convention B: `code:200`)

**Non-Goals:**
- No API contract changes — all endpoints, request/response formats remain identical
- No new Go packages or third-party dependencies
- No GORM association/relationship changes (stick with manual multi-query pattern)
- No input validation framework changes (keep bare type assertions in handlers per existing pattern)
- No pagination behavior changes (some endpoints return all records, some are paginated — preserve as-is)
- No timestamp handling changes (manual `int(time.Now().Unix())` stays)

## Decisions

### 1. Service constructor injection pattern

**Decision**: Handlers accept service pointers via constructor functions. `main.go` is the composition root where services are created and injected.

```go
// services/admin/article.go
type ArticleService struct{}
func NewArticleService() *ArticleService { return &ArticleService{} }

// handlers/admin/article.go (refactored)
type AdminHandler struct {
    articleSvc *admin.ArticleService
    ...other services...
}
func NewAdminHandler(articleSvc *admin.ArticleService, ...) *AdminHandler { ... }
```

**Rationale**: Follows the existing pattern established by `AdminHandler` / `ArticleService`. Enables future unit testing by allowing mock service injection. Consistent with Go community conventions.

**Alternatives considered**:
- Package-level singleton services (like `models.DB`): Rejected — prevents testing and creates hidden dependencies
- Interface-based injection: Rejected for now — over-engineering for a service layer that always uses the same DB implementation. Can add later if needed.

### 2. Service method signatures

**Decision**: Services return domain data + error. No `gin.Context` or HTTP concerns in services.

```go
// Service returns domain types
func (s *CategoryService) List(page, pageSize int) ([]models.Category, int64, error)

// Handler converts to HTTP response
func (h *CategoryHandler) List(c *gin.Context) {
    categories, total, err := h.categorySvc.List(page, pageSize)
    if err != nil {
        c.JSON(500, gin.H{"code": 500, "message": err.Error()})
        return
    }
    c.JSON(200, gin.H{"code": 0, "message": "success", "data": ...})
}
```

**Rationale**: Clean separation of concerns. Services are testable without HTTP context. Return values are plain Go types usable by any caller.

**Error return**: Currently handlers return `nil` on not-found. We'll add an `error` return for service methods, keeping the pattern of `nil` checks for get-by-ID operations. Services return `gorm.ErrRecordNotFound` for not-found; handlers map it to 404.

### 3. Response convention preservation

**Decision**: Response formatting (Convention A vs B) stays in handlers. Services are convention-agnostic.

**Convention A handlers** (articles, categories, comments, menus, tags, directories, website, frontend): return `{"code": 0, "message": "success", "data": ...}`
**Convention B handlers** (admins, admin_groups, attaches, langs, logs, messages, permits, reads): return `{"code": 200, "msg": "success", "data": ...}`

**Rationale**: Services shouldn't know about HTTP response conventions. Changing conventions later (e.g., unifying) would only touch handlers, not services.

### 4. Article tag association handling

**Decision**: Tag-enrichment logic (N+1: query article_tags join table, then query tags table) moves into `services/admin/article.go` and `services/frontend/article.go`. The `List` method returns articles with tags as a composite result type.

```go
// services/admin/article.go
type ArticleWithTags struct {
    models.Article
    Tags []models.Tags
}

func (s *ArticleService) ListWithTags(page, pageSize int) ([]ArticleWithTags, int64, error)
```

**Rationale**: The tag enrichment is currently duplicated in both `handlers/admin/article.go:ListArticles` and `handlers/frontend/handler.go:ListArticles`. Moving it into the service eliminates duplication and N+1 queries from the handler layer.

### 5. Login/auth service

**Decision**: Login logic (test user bypass, DB user lookup, password validation, token generation, last-login update) moves into `services/admin/admin.go` as `AdminService.Login(username, password, clientIP)`.

**Rationale**: Login is business logic — token generation, password hash validation, and login audit updates should not live in an HTTP handler. Same pattern for `ChangePassword` and `ChangeAdminPassword`.

### 6. File upload logic

**Decision**: File upload (multipart form handling, file saving to disk, DB record creation) moves into `services/admin/attach.go`. The handler delegates the `*multipart.FileHeader` to the service.

**Rationale**: Upload involves both filesystem I/O and DB writes — business logic territory. The handler remains responsible for binding the multipart form; the service handles the file save + DB insert.

### 7. Shared helpers

**Decision**: The `toInt` helper and password hashing functions (`generateTestToken`, `validatePassword`) move from `handlers/admin/handler.go` into the relevant service files where they're used. No shared `helpers.go` utility file — helpers live in the service that uses them.

**Rationale**: `toInt` is used by article service (type coercion from `interface{}`). Password functions belong in admin service. Avoiding a shared utility package prevents unnecessary coupling between services.

## Risks / Trade-offs

- **Risk**: Refactoring all handlers simultaneously could introduce regressions → **Mitigation**: Refactor one entity at a time (e.g., start with category — simple, low risk), verify with manual endpoint testing before proceeding. Follow the order in tasks.md.
- **Risk**: Import cycle between `services/admin` and `services/frontend` if services reference each other → **Mitigation**: Services never import each other. They only import `models`. All cross-service orchestration happens in handlers.
- **Risk**: Changing handler constructor signatures breaks `main.go` at build time → **Mitigation**: Update `main.go` immediately when a constructor changes. The compiler catches all wiring errors.
- **Trade-off**: Increased file count (19 new service files) adds directory complexity → Acceptable because each file maps 1:1 with a handler, making navigation predictable.
- **Trade-off**: No shared service interfaces means handlers are coupled to concrete service types → Acceptable for a monolith at this scale. Can add interfaces later if testing requires it.

## Open Questions

- Should `AdminHandler` be split into separate handler files (article handler, user handler) to match the 1:1 service:handler mapping? Currently `AdminHandler` contains both article methods and user management methods. → **Decision deferred** — out of scope for this refactor; can be done later.
</text>
