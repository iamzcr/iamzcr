## ADDED Requirements

### Requirement: WeChat configuration in website settings
The system SHALL store WeChat Official Account credentials as key-value pairs in the existing `sl_website` table using the following keys: `wechat_app_id`, `wechat_app_secret`, `wechat_token`, `wechat_aes_key`, `wechat_original_id`. All keys SHALL use the `staus` column default of 1 (active).

#### Scenario: Read WeChat configuration
- **WHEN** `WebsiteService.Get()` is called
- **THEN** the returned map SHALL include all configured WeChat keys with their values

#### Scenario: Save WeChat configuration
- **WHEN** `PUT /api/website` is called with `{"wechat_app_id": "wx123", "wechat_app_secret": "secret456"}`
- **THEN** the system SHALL upsert the keys into `sl_website` and return `code: 0`

#### Scenario: Missing WeChat configuration
- **WHEN** a publish attempt is made but any required WeChat key (`wechat_app_id`, `wechat_app_secret`) is missing or empty
- **THEN** the system SHALL return `code: 1` with message "微信公众号配置不完整，请先在基础设置中配置"

### Requirement: WeChat Admin UI in WebsiteSettings
The admin WebsiteSettings page SHALL display WeChat configuration fields with contextual placeholder text for each key. The fields SHALL be editable through the existing key-value add/edit modal.

#### Scenario: View WeChat settings in table
- **WHEN** the admin navigates to Website Settings
- **THEN** the settings table SHALL display all WeChat config keys alongside other settings, identifiable by their `wechat_` prefixed keys

#### Scenario: Add WeChat config key
- **WHEN** the admin creates a new setting with key `wechat_app_id` and value `wx1234567890`
- **THEN** the setting SHALL be saved via `PUT /api/website` and appear in the table

### Requirement: WeChat API client initialization
The system SHALL initialize a PowerWeChat Official Account client from the stored website settings when the admin API server starts. If WeChat credentials are missing, the client SHALL be nil and publish attempts SHALL return a clear configuration error.

#### Scenario: Client initialized with valid config
- **WHEN** the admin server starts and `wechat_app_id` and `wechat_app_secret` are both non-empty
- **THEN** a PowerWeChat Official Account instance SHALL be created and held by the WeChatService

#### Scenario: Client not initialized due to missing config
- **WHEN** the admin server starts and `wechat_app_id` is empty
- **THEN** the WeChatService client SHALL be nil; subsequent publish attempts SHALL return a configuration error
