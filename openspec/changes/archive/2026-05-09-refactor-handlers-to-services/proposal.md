## Why

Currently, all business logic (SQL queries, data transformation, business rules) lives directly in HTTP handler functions across 16 handler files. The only exception is `services/article.go` which provides partial service-layer coverage for articles — proving the pattern works but is incomplete. This makes the codebase hard to test, mixes HTTP and domain concerns, and duplicates common patterns (pagination, timestamps, password hashing) across handlers. Extracting a proper service layer split by admin/frontend responsibilities will improve testability, readability, and maintainability.

## What Changes

- Create `services/admin/` package with one service file per resource entity (14 services for admin operations: articles, categories, comments, menus, tags, directories, website, admins, admin_groups, attaches, langs, logs, messages, permits, reads)
- Create `services/frontend/` package with one service file per resource entity (5 services for public read operations: articles, categories, directories, tags, website)
- Refactor existing `services/article.go` into `services/admin/article.go` and `services/frontend/article.go` — split admin CRUD from public read logic
- Move all SQL queries, data transformations, and business logic from `handlers/admin/*.go` into the corresponding `services/admin/` files
- Move all SQL queries and data transformations from `handlers/frontend/handler.go` into `services/frontend/` files
- Move tag-enrichment N+1 logic from handlers into article services (currently duplicated in both admin and frontend handlers)
- Move login/auth logic from `AdminHandler.Login` into a dedicated admin user service
- Move password hashing/validation helpers from `handlers/admin/handler.go` into admin user service
- Move file upload logic from `handlers/admin/attach.go` into admin attach service
- **BREAKING**: Handler structs change constructor signatures to accept injected services
- Update `cmd/admin/main.go` and `cmd/frontend/main.go` to wire services into handlers at startup

## Capabilities

### New Capabilities

- `service-layer-foundation`: Base service patterns including constructor injection, shared response helpers, type coercion utilities (`toInt`), and service interface consistency across all entities
- `admin-content-services`: Business logic for admin-side content CRUD — articles (with tag associations), categories, directories, tags, menus, comments
- `admin-system-services`: Business logic for admin-side system CRUD — admins (with password hashing), admin_groups, attaches (with file upload), langs, logs, messages, permits, reads, website
- `frontend-services`: Business logic for public read-only queries — articles (published, filtered, with tags), categories, directories, tags, website

### Modified Capabilities

<!-- No existing specs to modify -->

## Impact

- **Affected code**: All 16 handler files in `handlers/admin/` and `handlers/frontend/`; both `cmd/admin/main.go` and `cmd/frontend/main.go`; `services/article.go` (split and migrated); new `services/admin/` and `services/frontend/` package directories
- **APIs**: No HTTP API changes — all endpoints, request/response formats, and response conventions remain identical
- **Dependencies**: No new third-party dependencies required; services use existing `models.DB` singleton
- **Risk**: Moderate — this is a pure internal refactor with no API contract changes, but touches most Go files. Mitigate by refactoring one entity at a time, testing each before proceeding
