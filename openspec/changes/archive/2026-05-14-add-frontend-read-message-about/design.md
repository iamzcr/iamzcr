## Context

The frontend SPA (web/frontend/) is a Vue 3 + TS + NaiveUI app served on port 3000, consuming the frontend API (port 8082). The frontend API currently has only 6 GET endpoints — no POST routes exist. This change adds three new capabilities: read tracking, a message guestbook, and an about page.

## Goals / Non-Goals

**Goals:**
- Record real article reads to `sl_read` table on article detail page load
- Return `read_count` from `sl_read` in article detail response
- Provide a guestbook page where visitors can view and submit messages
- Provide a static "About Me" page with bio and links
- Add both new pages to the top navigation bar

**Non-Goals:**
- No admin management for reads/messages (admin handlers already exist)
- No authentication for message submission (open to public)
- No pagination for message list (simple display, manageable volume)
- No email notifications on message submission

## Decisions

### 1. Read tracking on GetArticle (not separate endpoint)

**Choice**: Record read in `GetArticle` handler after confirming article exists. Pass IP and Referer from `c.Request`.

**Rationale**: One fewer API call from frontend. Read records every detail page view (including bots) — acceptable for a small personal blog.

### 2. Read count in GetByID response

**Choice**: Add `read_count` field to `GetByID` response map alongside `article`, `category`, `directory`, `tags`.

**Rationale**: Single query provides everything. Frontend accesses `data.read_count`.

### 3. Frontend message service uses existing Message model

**Choice**: `services/frontend/message.go` uses `models.Message` directly, same as admin message service.

**Rationale**: Both admin and frontend operate on the same `sl_message` table. No need for separate models.

### 4. Form debounce via button disable (not lodash/timer)

**Choice**: On submit, set `submitting = true`, disable button. Reset after response.

**Rationale**: Simple, reliable, no extra dependencies. Covers the practical need (prevent double-submit) without complex timing logic.

### 5. About page — static Vue component

**Choice**: Pure component, no API calls. CDN image URL constructed client-side from `websiteApi.get()`'s `cdn_url`.

**Rationale**: No backend needed for static content.

## Risks / Trade-offs

[Risk] Read count increments on bot traffic → **Mitigation**: Acceptable for initial implementation. Could add IP dedup later.

[Risk] Message spam without captcha → **Mitigation**: Low traffic personal blog; could add captcha later if needed.
