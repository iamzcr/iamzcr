## 1. Setup

- [x] 1.1 Add swaggo dependencies (`swaggo/swag`, `swaggo/gin-swagger`, `swaggo/files`) to go.mod
- [x] 1.2 Install `swag` CLI tool (`go install github.com/swaggo/swag/cmd/swag@latest`)

## 2. Response Models

- [x] 2.1 Create `handlers/frontend/response.go` with typed response structs for swaggo annotations (ArticleListData, ArticleDetailData, etc.)

## 3. API Documentation Annotations

- [x] 3.1 Add main-level Swagger annotations to `cmd/frontend/main.go` (`@title`, `@version`, `@host`, `@BasePath`)
- [x] 3.2 Add handler-level Swagger annotations to all 6 endpoints in `handlers/frontend/handler.go`

## 4. Swagger UI Route

- [x] 4.1 Import swagger docs and gin-swagger middleware in `cmd/frontend/main.go`
- [x] 4.2 Register Swagger UI route in `router/frontend.go`

## 5. Generate Documentation

- [x] 5.1 Run `swag init` to generate `docs/` directory (swagger.json, swagger.yaml, docs.go)
- [x] 5.2 Add `make docs` target to Makefile

## 6. Verify

- [x] 6.1 Build and run frontend API, verify Swagger UI loads at `/api/swagger/index.html`
- [x] 6.2 Verify all 6 endpoints are listed and their request/response schemas are correct
