## ADDED Requirements

### Requirement: Frontend article service
The frontend article service (`services/frontend/article.go`) SHALL provide `ListPublished` and `GetByID` methods for public read access. `ListPublished` SHALL only return articles with `status = 1` and SHALL support filtering by category ID (`cid`), directory ID (`did`), and tag ID (`tid`). The List method SHALL return articles with their associated tags pre-loaded. `GetByID` SHALL return an article with its category, directory, and tags.

#### Scenario: List published articles with tag filter
- **WHEN** `ArticleService.ListPublished(1, 10, "5", "", "3")` is called
- **THEN** it returns only published articles in category 5 that have tag 3, with each article including its `[]models.Tags`

#### Scenario: List all published articles unfiltered
- **WHEN** `ArticleService.ListPublished(1, 10, "", "", "")` is called
- **THEN** it returns all published articles with pagination and injected tags

#### Scenario: Get article detail with relations
- **WHEN** `ArticleService.GetByID(42)` is called
- **THEN** it returns a map with `article`, `category`, `directory`, and `tags` keys

#### Scenario: Get non-existent article
- **WHEN** `ArticleService.GetByID(99999)` is called with a non-existent ID
- **THEN** it returns `nil, error`

### Requirement: Frontend category service
The frontend category service (`services/frontend/category.go`) SHALL provide a `List` method that returns all categories with `status = 1` (no pagination).

#### Scenario: List active categories
- **WHEN** `CategoryService.List()` is called
- **THEN** it returns all categories where status equals 1

### Requirement: Frontend directory service
The frontend directory service (`services/frontend/directory.go`) SHALL provide a `List` method that returns all directories with `status = 1` (no pagination).

#### Scenario: List active directories
- **WHEN** `DirectoryService.List()` is called
- **THEN** it returns all directories where status equals 1

### Requirement: Frontend tags service
The frontend tags service (`services/frontend/tags.go`) SHALL provide a `List` method that returns all tags with `status = 1` (no pagination).

#### Scenario: List active tags
- **WHEN** `TagsService.List()` is called
- **THEN** it returns all tags where status equals 1

### Requirement: Frontend website service
The frontend website service (`services/frontend/website.go`) SHALL provide a `List` method that returns all website settings as a `map[string]string` (key-value pairs).

#### Scenario: List website settings as map
- **WHEN** `WebsiteService.List()` is called
- **THEN** it returns `map[string]string` where each key is the website setting name and value is its setting value
