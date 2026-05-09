## ADDED Requirements

### Requirement: Service constructor pattern
All service structs SHALL follow a consistent constructor pattern: zero-argument `New*Service()` function that returns a pointer to the service struct. Services SHALL NOT accept dependencies through their constructor (they use the package-level `models.DB` singleton). Service structs SHALL be defined in `services/admin/` or `services/frontend/` matching the handler domain boundary.

#### Scenario: Creating an admin category service
- **WHEN** `services/admin.NewCategoryService()` is called
- **THEN** a pointer to `CategoryService` is returned with no initialization needed

#### Scenario: Creating a frontend article service
- **WHEN** `services/frontend.NewArticleService()` is called
- **THEN** a pointer to `ArticleService` is returned ready for use

### Requirement: Service method signatures
All service methods SHALL accept only Go primitive types, models, or map[string]interface{} as parameters. Service methods SHALL NOT accept `*gin.Context` or any HTTP framework types. Service methods SHALL return domain data types (model structs, slices, maps) plus an `error` value. List methods SHALL return `([]ModelType, int64, error)` where the second value is total count.

#### Scenario: Category list method
- **WHEN** `CategoryService.List(1, 10)` is called
- **THEN** the method returns `([]models.Category, int64, error)` with categories, total count, and nil error on success

#### Scenario: Record not found
- **WHEN** a Get method queries for a non-existent ID
- **THEN** the method returns `nil, gorm.ErrRecordNotFound` or equivalent error

### Requirement: Service isolation from HTTP
Service methods SHALL NOT format HTTP response bodies or set HTTP status codes. Response formatting (gin.H with code/message/data) SHALL remain in handler functions. Service errors SHALL be returned as Go errors; handlers SHALL map them to appropriate HTTP responses.

#### Scenario: Service returns error
- **WHEN** a service method encounters a database error
- **THEN** it returns the error to the caller without formatting it as a JSON response

### Requirement: Handler service injection
Handler structs SHALL accept service pointers via their constructor functions. Each handler struct field for a service SHALL be unexported. The `cmd/*/main.go` files SHALL be the composition root where all services are created and injected into handlers.

#### Scenario: Admin handler construction
- **WHEN** `admin.NewCategoryHandler(categorySvc)` is called with a valid service pointer
- **THEN** the handler is initialized with access to the category service

#### Scenario: Frontend handler construction
- **WHEN** `frontend.NewFrontendHandler(articleSvc, categorySvc, directorySvc, tagsSvc, websiteSvc)` is called
- **THEN** the handler is initialized with all service dependencies

### Requirement: Shared helper migration
The `toInt` type coercion helper SHALL be moved from `handlers/admin/handler.go` to the service file where it is used. Password hashing functions (`validatePassword`, MD5+salt hashing) SHALL be moved into `services/admin/admin.go`. Token generation SHALL be moved into `services/admin/admin.go`.

#### Scenario: Type coercion helper location
- **WHEN** looking for the `toInt` function
- **THEN** it is defined in `services/admin/article.go` (its only consumer)

#### Scenario: Password validation helper location
- **WHEN** looking for the `validatePassword` function
- **THEN** it is defined in `services/admin/admin.go`
