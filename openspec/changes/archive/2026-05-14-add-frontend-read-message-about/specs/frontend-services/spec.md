## MODIFIED Requirements

### Requirement: Frontend article service
The frontend article service (`services/frontend/article.go`) SHALL provide `ListPublished` and `GetByID` methods for public read access. `ListPublished` SHALL only return articles with `status = 1` and SHALL support filtering by category ID (`cid`), directory ID (`did`), and tag ID (`tid`). The List method SHALL return articles with their associated tags pre-loaded. `GetByID` SHALL return an article with its category, directory, tags, and the read count from `sl_read`.

#### Scenario: List published articles with tag filter
- **WHEN** `ArticleService.ListPublished(1, 10, "5", "", "3")` is called
- **THEN** it returns only published articles in category 5 that have tag 3, with each article including its `[]models.Tags`

#### Scenario: List all published articles unfiltered
- **WHEN** `ArticleService.ListPublished(1, 10, "", "", "")` is called
- **THEN** it returns all published articles with pagination and injected tags

#### Scenario: Get article detail with relations and read count
- **WHEN** `ArticleService.GetByID(42)` is called
- **THEN** it returns a map with `article`, `category`, `directory`, `tags`, and `read_count` keys
- **AND** `read_count` equals the number of rows in `sl_read` where `aid = 42`

#### Scenario: Get non-existent article
- **WHEN** `ArticleService.GetByID(99999)` is called with a non-existent ID
- **THEN** it returns `nil, error`

## ADDED Requirements

### Requirement: Frontend message service
The frontend message service (`services/frontend/message.go`) SHALL provide `List` and `Create` methods. `List` SHALL return all messages ordered by `create_time DESC`. `Create` SHALL insert a new message with name, email, url, content, and IP, and set `is_reply = 0`.

#### Scenario: List all messages
- **WHEN** `MessageService.List()` is called
- **THEN** it returns all `sl_message` rows ordered by `create_time` descending

#### Scenario: Create a new message
- **WHEN** `MessageService.Create("name", "a@b.com", "", "Hello", "1.2.3.4")` is called
- **THEN** a new row is inserted into `sl_message` with `is_reply = 0` and `create_time` set to current timestamp

### Requirement: Frontend API message endpoints
The frontend API SHALL expose `GET /api/messages` and `POST /api/messages` endpoints handled by the frontend handler's `GetMessages` and `CreateMessage` methods.

#### Scenario: Get messages via API
- **WHEN** `GET /api/messages` is requested
- **THEN** it returns Convention A response with message list in `data`

#### Scenario: Create message via API
- **WHEN** `POST /api/messages` is requested with valid email and content
- **THEN** it creates a message record and returns the created message
