## Context

The frontend API (port 8082) serves a minimal set of 6 public endpoints for the blog frontend SPA. A WeChat Mini Program is being developed that will consume these same endpoints. The current handler code uses opaque `gin.H` maps for all responses, and there is zero API documentation. The team needs auto-generated, always-up-to-date Swagger documentation.

Project constraints from AGENTS.md:
- Response format: Convention A (`{"code": 0, "message": "success", "data": ...}`)
- All timestamp fields are `int` (Unix epoch seconds)
- No GORM auto-timestamps or associations
- Go module: `iamzcr`

## Goals / Non-Goals

**Goals:**
- Generate Swagger 2.0 (OpenAPI) documentation from Go source annotations
- Serve interactive Swagger UI at `/api/swagger/index.html`
- Document all 6 existing public endpoints with accurate request/response schemas
- Keep existing API behavior and response formats unchanged
- Add `make docs` command for easy regeneration

**Non-Goals:**
- Migrate to OpenAPI 3.0 (use Swagger 2.0 since swaggo is the standard Go tool)
- Document admin API (port 8081) — focused on frontend API only
- Add request validation or change handler logic
- Auto-generate client SDKs from the spec

## Decisions

### 1. Use swaggo suite (not hand-written OpenAPI spec)

**Choice**: `github.com/swaggo/swag` (CLI) + `github.com/swaggo/gin-swagger` (Gin middleware) + `github.com/swaggo/files` (UI assets).

**Rationale**: swaggo is the de facto standard for Go Swagger generation. Annotations live in the source code alongside handler logic, ensuring docs stay in sync. Alternatives considered:
- Hand-written OpenAPI YAML → easy to drift from code
- `go-swagger` → heavier, requires separate spec file
- `ogen` → targets OpenAPI 3.0 code generation, overkill for documentation

### 2. Define typed response model structs instead of documenting gin.H

**Choice**: Create a new file `handlers/frontend/response.go` with typed response structs that existing handlers continue to populate via `gin.H`. The structs exist solely for swaggo annotation purposes — they are NOT used at runtime by handlers (keeping existing handler code untouched).

**Rationale**: swaggo cannot introspect `gin.H` (which is `map[string]interface{}`). Typed structs with `json` tags are required for swaggo to generate accurate response schemas. Using separate annotation-only structs avoids touching existing handler logic.

### 3. Single Swagger tag: "Frontend API"

**Choice**: All 6 endpoints share one tag `Frontend API`. No sub-grouping.

**Rationale**: Only 6 endpoints exist — splitting into sub-tags would add noise. The Mini Program team will consume the entire API as one surface.

### 4. Response envelope: wrap actual data in a generic structure

**Choice**: Define `APIResponse` / `APIError` wrapper structs matching Convention A (`{"code": 0, "message": "success", "data": <T>}`), reused across all endpoints via swaggo's `@Success` annotation with `{object}` type reference.

**Rationale**: Avoids duplicating the envelope definition for each endpoint.

## Risks / Trade-offs

[Risk] Annotations could drift from actual handler behavior if handlers are modified without updating annotations → **Mitigation**: Add `make docs` command for easy regeneration; the generated `docs/` is committed so changes are visible in PRs.

[Risk] swaggo version compatibility with Go 1.25 → **Mitigation**: Use latest swaggo releases; Go 1.25 should be backward compatible.

[Risk] Typed response structs in `response.go` may become inconsistent with actual response construction in handlers → **Mitigation**: Keep structs simple and mirror the handler's `gin.H` construction; document in code comments that these are annotation-only structs.
