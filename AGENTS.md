# AGENTS.md

個人博客系統 "堆棧人生" (Stack Life). Go + Gin + GORM backend, Vue 3 + TS + Vite + NaiveUI frontend, MySQL (`stacklifes`).

## Architecture

Two **separate** Go API servers and two **separate** Vue SPAs — they are independent processes:

| Process | Port | Dir | Auth |
|---|---|---|---|
| Admin API | 8081 | `cmd/admin/` | JWT (all routes except `/api/login`) |
| Frontend API | 8082 | `cmd/frontend/` | None (public) |
| Admin SPA | 3001 | `web/admin/` | Token in localStorage |
| Blog SPA | 3000 | `web/frontend/` | None (public) |

The Vue SPAs proxy `/api` to their respective Go backend via Vite dev server config. No shared code between the two SPAs.

## Development commands

```bash
# Go hot reload (requires `go install github.com/air-verse/air@latest`)
air -c .air_admin.toml       # admin API with hot reload
air -c .air_frontend.toml    # frontend API with hot reload

# Go run without hot reload
make run-admin               # port 8081
make run-frontend            # port 8082
make build                   # compile both binaries

# Frontend dev servers
cd web/frontend && npm run dev   # port 3000
cd web/admin && npm run dev      # port 3001

# Frontend build (type-check first, then build)
cd web/frontend && npm run build  # runs vue-tsc -b && vite build
cd web/admin && npm run build     # runs vue-tsc -b && vite build

# Database init
mysql -u root -p stacklifes < sql/stacklifes.sql
```

## Critical gotchas

### Response format inconsistency (backend)
Two different response conventions exist in handlers. **Match whichever one the file you're editing already uses:**

- **Convention A** (most handlers): `{"code": 0, "message": "success", "data": ...}` — success code is `0`
- **Convention B** (admin user management, attaches, langs, logs, messages, permits, reads): `{"code": 200, "msg": "success", "data": ...}` — success code is `200`, uses `"msg"` not `"message"`

The frontends expect Convention A (check `res.data.code === 0`). When adding new endpoints, use Convention A unless the existing file already uses B.

### JWT auth (hardcoded)
- JWT secret: `"blog-admin-secret-key-2024"` (in `middleware/jwt.go`)
- Test user: `test` / `admin123` bypasses DB auth; always gets `user_id=1, username="test"` in token
- Token expiry: 24 hours, HS256
- Auth failure returns `{"code": 401, "message": "未登录"}` or `"登录已过期，请重新登录"`

### NaiveUI conventions (frontend)
- Use `v-model:value` (NOT bare `v-model`) for all NaiveUI form inputs
- NDataTable pagination: use `remote` prop + `:pagination` with `itemCount` (NOT `total`) + `@update:page`/`@update:page-size`
- Custom cell rendering: use `h()` from Vue in `render` functions (e.g., `render: (row) => h(NTag, { type: 'success' }, () => '已发布')`)
- Form submission uses `useMessage()` for toast notifications

### No path aliases
Imports must use relative paths (`../../api`, `../views/...`). There is no `@/` alias configured in `tsconfig` or `vite.config`.

### No state management library
No Pinia, Vuex, or provide/inject. State is managed via component-local `ref()`/`reactive()` + `localStorage` for auth tokens and admin info.

### Vue Router quirks
- `Home.vue` serves **4 routes** (`/`, `/category/:id`, `/directory/:id`, `/tag/:id`) — it inspects `route.name` to determine the filter
- `ArticleEdit.vue` serves **both** create (`/articles/new`) and edit (`/articles/edit/:id`) — checks `!!route.params.id`
- Admin router uses `window.location.href = '/login'` (full page reload) for 401 redirects, NOT `router.push()`

### Database/model quirks
- All timestamp fields are `int` (Unix epoch seconds), set manually in handler code — not GORM auto-timestamps
- No GORM associations; all joins are manual multi-query (N+1 pattern in article listing with tags)
- Model `Website.Staus` has a typo (missing `t`), mapped to DB column `staus`
- Pagination is inconsistent: some handlers return all records without pagination (admins, admin_groups, frontend categories/directories/tags)
- `models.DB` is a package-level `*gorm.DB` singleton; all handlers access it directly

### No input validation
Handlers use bare type assertions (`data["title"].(string)`), which will panic if a field is missing or the wrong type. Be defensive.

## Adding a new resource

1. Add handler in `handlers/admin/` or `handlers/frontend/` with 5 methods: `List`, `Get`, `Create`, `Update`, `Delete` (follow existing naming)
2. Register routes in `cmd/admin/main.go` or `cmd/frontend/main.go` using the `New*Handler` / route-registration pattern
3. Add API module in `web/admin/src/api/index.ts` or `web/frontend/src/api/index.ts`
4. Create Vue view in `web/admin/src/views/` and add route in `web/admin/src/main.ts`

## Test user
- Login to admin: `test` / `admin123`
- Token is stored in `localStorage` as `admin_token`; admin info cached as `admin_info`
