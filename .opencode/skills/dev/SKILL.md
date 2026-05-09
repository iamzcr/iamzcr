---
name: dev
description: Start project dev servers and open browser for admin or frontend stacks
compatibility: opencode
---

## What I do
- Start the Go backend API server
- Start the Vue SPA frontend dev server
- Open the browser to the appropriate URL

## When to use me
Use this when the user asks to "run dev", "start the project", or open the project in browser.

## Available stacks

### Admin stack (管理后台)
- Backend: Admin API on port **8081** → `go run ./cmd/admin/main.go`
- Frontend: Admin SPA on port **3001** → `cd web/admin; npm run dev`
- Browser URL: http://localhost:3001
- Login: `test` / `admin123`

### Frontend stack (博客前台)
- Backend: Frontend API on port **8082** → `go run ./cmd/frontend/main.go`
- Frontend: Blog SPA on port **3000** → `cd web/frontend; npm run dev`
- Browser URL: http://localhost:3000

## Steps

Ask the user which stack to run (admin or frontend), or if both.

For each chosen stack:

1. Start the Go backend in a background/separate terminal process:
   - Admin: `go run ./cmd/admin/main.go` (port 8081)
   - Frontend: `go run ./cmd/frontend/main.go` (port 8082)
2. Start the Vue dev server in a background/separate terminal process:
   - Admin: `cd web/admin; npm run dev` (port 3001)
   - Frontend: `cd web/frontend; npm run dev` (port 3000)
3. Wait a few seconds for the dev server to start, then open the browser:
   - Windows: `start http://localhost:<port>`
   - macOS: `open http://localhost:<port>`
   - Linux: `xdg-open http://localhost:<port>`

## Notes
- Both the Go backend and Vue dev server should run simultaneously in separate terminals
- The Vue dev servers proxy `/api` to their respective Go backend
- Use `air` for hot reload if available, otherwise use `go run`
