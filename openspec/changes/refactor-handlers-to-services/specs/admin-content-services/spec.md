## ADDED Requirements

### Requirement: Admin article service CRUD
The admin article service (`services/admin/article.go`) SHALL provide methods to list all articles (including drafts), create articles with tag associations, update articles with tag replacement, delete articles with cascade tag cleanup, and retrieve a single article with its category/directory/tags. The List method SHALL return articles enriched with their associated tags (not plain articles).

#### Scenario: List articles with tags
- **WHEN** `ArticleService.ListWithTags(1, 10)` is called
- **THEN** it returns `([]ArticleWithTags, int64, error)` where each article includes its `[]models.Tags`

#### Scenario: Create article with tags
- **WHEN** `ArticleService.Create(data, tagIDs)` is called with article fields and a list of tag IDs
- **THEN** the article is created and corresponding `ArticleTags` join records are inserted

#### Scenario: Update article and replace tags
- **WHEN** `ArticleService.Update(id, data, tagIDs)` is called
- **THEN** existing `ArticleTags` for that article are deleted and replaced with new ones for the provided tag IDs

#### Scenario: Delete article with tag cleanup
- **WHEN** `ArticleService.Delete(id)` is called
- **THEN** the article and all associated `ArticleTags` records are deleted

### Requirement: Admin category service
The admin category service (`services/admin/category.go`) SHALL provide standard CRUD methods: `List`, `Get`, `Create`, `Update`, `Delete`. The `List` method SHALL support pagination with page and pageSize parameters. The `Create` and `Update` methods SHALL set `create_time` and `update_time` fields as Unix timestamps.

#### Scenario: List categories with pagination
- **WHEN** `CategoryService.List(1, 10)` is called
- **THEN** it returns `([]models.Category, int64, error)` ordered by weight DESC

#### Scenario: Create category sets timestamps
- **WHEN** `CategoryService.Create(category)` is called
- **THEN** `CreateTime` and `UpdateTime` are set to `int(time.Now().Unix())`

### Requirement: Admin directory service
The admin directory service (`services/admin/directory.go`) SHALL provide standard CRUD methods for directory entries with pagination on List.

#### Scenario: List directories
- **WHEN** `DirectoryService.List(1, 10)` is called
- **THEN** it returns directories ordered by weight DESC with total count

### Requirement: Admin tags service
The admin tags service (`services/admin/tags.go`) SHALL provide standard CRUD methods for tags with pagination on List.

#### Scenario: List tags with pagination
- **WHEN** `TagsService.List(1, 10)` is called
- **THEN** it returns tags ordered by weight DESC with total count

### Requirement: Admin menu service
The admin menu service (`services/admin/menu.go`) SHALL provide standard CRUD methods for menus with pagination on List.

#### Scenario: List menus
- **WHEN** `MenuService.List(1, 10)` is called
- **THEN** it returns menus ordered by weight DESC

### Requirement: Admin comment service
The admin comment service (`services/admin/comment.go`) SHALL provide standard CRUD methods for comments with pagination on List.

#### Scenario: List comments with pagination
- **WHEN** `CommentService.List(1, 10)` is called
- **THEN** it returns comments ordered by create_time DESC with total count
