## MODIFIED Requirements

### Requirement: Admin article service CRUD
The admin article service (`services/admin/article.go`) SHALL provide methods to list all articles (including drafts), create articles with tag associations and optional WeChat publishing, update articles with tag replacement and optional WeChat publishing, delete articles with cascade tag cleanup, and retrieve a single article with its category/directory/tags. The List method SHALL return articles enriched with their associated tags (not plain articles). The Create and Update methods SHALL accept an optional `publish_to_wechat` boolean flag; when true, the article SHALL be published to WeChat after the save operation completes.

#### Scenario: List articles with tags
- **WHEN** `ArticleService.ListWithTags(1, 10)` is called
- **THEN** it returns `([]ArticleWithTags, int64, error)` where each article includes its `[]models.Tags`

#### Scenario: Create article with tags
- **WHEN** `ArticleService.Create(data, tagIDs)` is called with article fields and a list of tag IDs
- **THEN** the article is created and corresponding `ArticleTags` join records are inserted

#### Scenario: Create article with WeChat publish
- **WHEN** `ArticleService.Create(data, tagIDs)` is called and `data["publish_to_wechat"]` is `true`
- **THEN** after the article is saved, the system SHALL attempt to publish it to WeChat; if publish fails, the article SHALL still be saved and the error SHALL be returned alongside the article data

#### Scenario: Update article and replace tags
- **WHEN** `ArticleService.Update(id, data, tagIDs)` is called
- **THEN** existing `ArticleTags` for that article are deleted and replaced with new ones for the provided tag IDs

#### Scenario: Update article with WeChat publish
- **WHEN** `ArticleService.Update(id, data, tagIDs)` is called and `data["publish_to_wechat"]` is `true`
- **THEN** after the article is updated, the system SHALL attempt to publish it to WeChat; if already published, SHALL re-publish (update the existing WeChat draft)

#### Scenario: Delete article with tag cleanup
- **WHEN** `ArticleService.Delete(id)` is called
- **THEN** the article and all associated `ArticleTags` records are deleted
