# AGENT.md

## Project Overview

This is a personal blog system called "堆栈人生" (Stack Life), built with:
- **Backend**: Go + Gin + GORM + MySQL
- **Frontend**: Vue 3 + TypeScript + Vite + Naive UI
- **Database**: MySQL (stacklifes)

## Project Structure

```
iamzcr/
├── cmd/
│   ├── admin/          # Admin API server (port 8081)
│   └── frontend/       # Frontend API server (port 8082)
├── config/             # Configuration
├── handlers/           # API handlers
│   ├── admin/          # Admin endpoints
│   └── frontend/       # Frontend endpoints
├── middleware/         # Middleware (CORS, Auth, JWT)
├── models/             # Database models
├── services/           # Business logic
├── sql/                # SQL scripts
├── web/
│   ├── admin/          # Admin frontend (Vue 3, port 3001)
│   └── frontend/       # Blog frontend (Vue 3, port 3000)
├── .env                # Environment configuration
├── README.md           # Project documentation
└── AGENT.md            # Agent instructions (this file)
```

## Environment Configuration

Default ports and settings in `.env`:
- DB_HOST: localhost
- DB_PORT: 3306
- DB_USER: root
- DB_PASSWORD: root
- DB_NAME: stacklifes
- ADMIN_PORT: 8081
- FRONTEND_PORT: 8082

## Running the Project

### Backend

```bash
# Start admin API server (port 8081)
go run ./cmd/admin

# Start frontend API server (port 8082)  
go run ./cmd/frontend
```

### Frontend

```bash
# Frontend blog (port 3000)
cd web/frontend
npm install
npm run dev

# Admin panel (port 3001)
cd web/admin
npm install
npm run dev
```

## Key Files and Modifications

### Frontend Blog (web/frontend)

- `App.vue`: Main layout with navigation, logo "堆栈人生", footer
- `views/Home.vue`: Home page with article list
- `views/Article.vue`: Article detail page
- `api/index.ts`: API calls to frontend API (port 8082)

### Admin Panel (web/admin)

- `App.vue`: Main layout with sidebar menu, header with collapse button
- `views/Index.vue`: Dashboard with stats (articles, categories, tags count)
- `views/ArticleList.vue`: Article list with pagination
- `views/ArticleEdit.vue`: Article create/edit form
- `views/Login.vue`: Login page
- `api/index.ts`: API calls to admin API (port 8081)

### Backend Handlers

All handlers in `handlers/admin/` return paginated response format:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100
  }
}
```

## Known Issues and Fixes

1. **Menu hierarchy**: Admin sidebar loads menus from `sl_menu` table, builds hierarchy based on `parent` field, uses separate key for parent nodes to avoid NaiveUI menu state conflict

2. **Pagination**: Frontend uses NaiveUI DataTable with `itemCount` property (not `total`) and `remote` prop for server-side pagination

3. **Article edit**: Categories/directories/tags API returns `{list, total}`, ArticleEdit.vue handles both formats with fallback

4. **Login page**: Pure white background, proper spacing between title and inputs

5. **Sidebar collapse**: Default collapsed, toggle button in header with icon style

## Common Tasks

### Add new admin list page
1. Create Vue component in `web/admin/src/views/`
2. Add route in `web/admin/src/main.ts`
3. Add menu mapping in `App.vue` handleMenuSelect function

### Modify frontend blog branding
- Edit `web/frontend/src/App.vue` - logo text and footer

### Add new API endpoint
1. Add handler in appropriate file under `handlers/admin/` or `handlers/frontend/`
2. Register route in `cmd/admin/main.go` or `cmd/frontend/main.go`
3. Add API function in corresponding `web/*/src/api/index.ts`