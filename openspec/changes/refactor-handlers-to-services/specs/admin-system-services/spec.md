## ADDED Requirements

### Requirement: Admin user service with auth
The admin user service (`services/admin/admin.go`) SHALL provide methods for admin user CRUD (`List`, `Get`, `Create`, `Update`, `Delete`), password management (`ChangePassword`, `ChangeAdminPassword`), and login authentication (`Login`). The `Login` method SHALL handle the test user bypass (`test`/`admin123`), database user lookup with password validation (MD5+salt), update last-login fields, and generate JWT tokens.

#### Scenario: Login with test user bypass
- **WHEN** `AdminService.Login("test", "admin123", "127.0.0.1")` is called
- **THEN** it returns admin info with id=999, username="test", group="super admin", and a valid JWT token

#### Scenario: Login with database user
- **WHEN** `AdminService.Login("realuser", "password", "127.0.0.1")` is called with valid credentials
- **THEN** it looks up the user in the database, validates the MD5+salted password, updates `login_num` and `last_login_time`, and returns a JWT token

#### Scenario: Create admin with password hashing
- **WHEN** `AdminService.Create(admin)` is called
- **THEN** a random salt is generated, the password is hashed with MD5(password+salt), and `CreateTime` is set

#### Scenario: Change own password
- **WHEN** `AdminService.ChangePassword(userID, oldPwd, newPwd)` is called with correct old password
- **THEN** the old password is validated, the new password is hashed with a new salt, and the admin record is updated

### Requirement: Admin group service
The admin group service (`services/admin/admin_group.go`) SHALL provide standard CRUD methods for admin groups with List returning all records (no pagination, matching existing behavior).

#### Scenario: List all admin groups
- **WHEN** `AdminGroupService.List()` is called
- **THEN** it returns all admin groups without pagination

### Requirement: Admin attach service with file upload
The admin attach service (`services/admin/attach.go`) SHALL provide standard CRUD methods for attaches plus a `Upload` method. The `Upload` method SHALL accept a multipart file header, save the file to disk in the asset directory, and create an attach database record.

#### Scenario: Upload file
- **WHEN** `AttachService.Upload(fileHeader, assetDir)` is called with a valid multipart file
- **THEN** the file is saved to `assetDir/somehash.ext`, and an attach record is created in the database with the file path

#### Scenario: List attaches with pagination
- **WHEN** `AttachService.List(1, 10)` is called
- **THEN** it returns attaches ordered by create_time DESC with total count

### Requirement: Admin lang service
The admin lang service (`services/admin/lang.go`) SHALL provide standard CRUD methods for language entries.

#### Scenario: List langs
- **WHEN** `LangService.List()` is called
- **THEN** it returns all lang records without pagination (CONVENTION B endpoint)

### Requirement: Admin log service
The admin log service (`services/admin/log.go`) SHALL provide List, Get, Create, Delete methods (no Update). List SHALL support pagination.

#### Scenario: List logs with pagination
- **WHEN** `LogService.List(1, 10)` is called
- **THEN** it returns logs ordered by create_time DESC with total count

### Requirement: Admin message service
The admin message service (`services/admin/message.go`) SHALL provide standard CRUD methods for messages.

#### Scenario: List messages
- **WHEN** `MessageService.List(1, 10)` is called
- **THEN** it returns messages ordered by create_time DESC with total count

### Requirement: Admin permit service
The admin permit service (`services/admin/permit.go`) SHALL provide standard CRUD methods for permits.

#### Scenario: List permits
- **WHEN** `PermitService.List()` is called
- **THEN** it returns all permit records (CONVENTION B, currently returns all without pagination)

### Requirement: Admin read service
The admin read service (`services/admin/read.go`) SHALL provide List, Get, Create, Delete methods (no Update). List SHALL support pagination.

#### Scenario: List reads with pagination
- **WHEN** `ReadService.List(1, 10)` is called
- **THEN** it returns reads ordered by create_time DESC with total count

### Requirement: Admin website service
The admin website service (`services/admin/website.go`) SHALL provide List (all records), GetByKey (returns single website setting by key), Upsert (create or update by key), and Delete methods.

#### Scenario: Upsert website setting
- **WHEN** `WebsiteService.Upsert(key, value)` is called with an existing key
- **THEN** the existing record is updated; if the key does not exist, a new record is created

#### Scenario: Get all website settings
- **WHEN** `WebsiteService.List()` is called
- **THEN** it returns all website records
