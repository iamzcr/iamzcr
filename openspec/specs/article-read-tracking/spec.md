## ADDED Requirements

### Requirement: Record article read on detail view
The system SHALL insert a record into `sl_read` when a visitor views an article detail page, capturing the article ID, visitor IP, HTTP Referer header, and current timestamp.

#### Scenario: Read recorded on article detail
- **WHEN** a visitor requests `GET /api/articles/42`
- **AND** article 42 exists
- **THEN** a new row is inserted into `sl_read` with `aid=42`, `ip=<visitor IP>`, `referer=<Referer header>`, `create_time=<now>`

#### Scenario: Read not recorded for missing article
- **WHEN** a visitor requests `GET /api/articles/99999`
- **AND** article 99999 does not exist
- **THEN** no row is inserted into `sl_read`

### Requirement: Return real read count with article detail
The system SHALL include `read_count` (the number of rows in `sl_read` where `aid = <article id>`) in the `GET /api/articles/:id` response data.

#### Scenario: Read count returned
- **WHEN** a visitor requests `GET /api/articles/42`
- **THEN** the response `data` object contains `read_count: <integer>`
- **AND** `read_count` equals `SELECT COUNT(*) FROM sl_read WHERE aid = 42`

### Requirement: Frontend displays real read count
The blog SPA article detail page SHALL display `read_count` instead of `article.weight` as the reading count.

#### Scenario: Read count displayed
- **WHEN** a visitor navigates to `/article/42`
- **THEN** the page displays the value of `read_count` (e.g., "阅读: 128")
