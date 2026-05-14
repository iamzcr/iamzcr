## Why

We are about to develop a WeChat Mini Program frontend that will consume the frontend API (port 8082). The Mini Program team needs clear, machine-readable API documentation to understand all endpoints, request/response formats, and data models. Currently, the project has zero API documentation — all 6 public endpoints are undocumented, forcing developers to read Go source code to understand the API contract.

## What Changes

- Add `swaggo/swag` + `github.com/swaggo/gin-swagger` dependencies to generate OpenAPI 2.0 (Swagger) documentation from Go annotations
- Add Swagger annotations (`@Summary`, `@Tags`, `@Param`, `@Success`, `@Failure`, `@Router`) to all 6 frontend API handlers
- Add main-level annotations (`@title`, `@version`, `@host`, `@BasePath`) in `cmd/frontend/main.go`
- Define typed response model structs with `swaggo` annotations to replace opaque `gin.H` maps
- Add Swagger UI route (`/api/swagger/*any`) to serve interactive documentation
- Add `make docs` command to regenerate `docs/` directory via `swag init`

## Capabilities

### New Capabilities

- `frontend-api-docs`: Auto-generated Swagger UI at `/api/swagger/index.html` serving as the single source of truth for the frontend API contract, viewable in any browser with "Try it out" functionality.

### Modified Capabilities

<!-- No existing specs are being modified -->

## Impact

- **Dependencies**: Add `github.com/swaggo/swag` (CLI), `github.com/swaggo/gin-swagger` (middleware), `github.com/swaggo/files` (swagger UI assets)
- **Code**: Annotate `cmd/frontend/main.go`, `handlers/frontend/handler.go`; modify `router/frontend.go` to register Swagger UI route; add `docs/` directory with generated files
- **Build**: New `make docs` target; `docs/` checked into version control
- **Breaking**: None — all existing API routes and response formats remain unchanged
