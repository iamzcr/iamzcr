## ADDED Requirements

### Requirement: Swagger UI endpoint
The system SHALL serve an interactive Swagger UI at `/api/swagger/index.html` that displays all frontend API endpoints with request/response documentation.

#### Scenario: Access Swagger UI
- **WHEN** a developer navigates to `http://localhost:8082/api/swagger/index.html`
- **THEN** the system displays the Swagger UI page showing all documented API endpoints grouped under "Frontend API" tag

#### Scenario: Swagger JSON endpoint
- **WHEN** a developer requests `GET /api/swagger/doc.json`
- **THEN** the system returns the generated OpenAPI 2.0 JSON specification including all 6 endpoints with their request parameters and response schemas

### Requirement: Document article list endpoint
The system SHALL document `GET /api/articles` with query parameters for pagination (`page`, `page_size`) and filtering (`cid`, `did`, `tid`), and describe the response structure including the article object with all fields.

#### Scenario: Article list documentation visible
- **WHEN** Swagger UI renders the document for GET /api/articles
- **THEN** it shows query parameters: page (default 1), page_size (default 10), cid, did, tid as optional integer filters
- **AND** the 200 response schema describes an object with code, message, data.list (array of articles), and data.total (integer)

### Requirement: Document article detail endpoint
The system SHALL document `GET /api/articles/:id` with the `id` path parameter and describe the response structure including article, category, directory, and tags.

#### Scenario: Article detail documentation visible
- **WHEN** Swagger UI renders the document for GET /api/articles/{id}
- **THEN** it shows the required integer path parameter `id`
- **AND** the 200 response schema describes the data object containing article, category, directory, and tags fields

### Requirement: Document static data endpoints
The system SHALL document `GET /api/categories`, `GET /api/directories`, and `GET /api/tags` as endpoints that return lists without pagination parameters.

#### Scenario: Categories documentation visible
- **WHEN** Swagger UI renders the document for GET /api/categories
- **THEN** it shows the endpoint returns a list of category objects with id and name fields

#### Scenario: Directories documentation visible
- **WHEN** Swagger UI renders the document for GET /api/directories
- **THEN** it shows the endpoint returns a list of directory objects with id, name, and cid fields

#### Scenario: Tags documentation visible
- **WHEN** Swagger UI renders the document for GET /api/tags
- **THEN** it shows the endpoint returns a list of tag objects with id and name fields

### Requirement: Document website settings endpoint
The system SHALL document `GET /api/website` as an endpoint that returns a key-value map of website configuration settings.

#### Scenario: Website settings documentation visible
- **WHEN** Swagger UI renders the document for GET /api/website
- **THEN** it shows the endpoint returns a map of string key-value pairs in the data field

### Requirement: Regeneratable documentation
The system SHALL provide a `make docs` command that runs `swag init` to regenerate the `docs/` directory from source annotations without modifying any runtime code.

#### Scenario: Regenerate docs
- **WHEN** a developer runs `make docs` from the project root
- **THEN** the `docs/` directory is regenerated with updated swagger.json, swagger.yaml, and docs.go files reflecting current annotations
