## ADDED Requirements

### Requirement: List messages
The frontend API SHALL expose `GET /api/messages` returning all messages ordered by creation time descending.

#### Scenario: Retrieve message list
- **WHEN** a visitor requests `GET /api/messages`
- **THEN** the system returns all `sl_message` rows in descending order by `create_time`
- **AND** the response follows Convention A with `data` as an array of message objects

### Requirement: Submit message with required fields
The frontend API SHALL expose `POST /api/messages` accepting a JSON body with `name`, `email`, `url`, and `content`. `email` and `content` SHALL be required; requests missing either SHALL return HTTP 400 with an error message.

#### Scenario: Successful message submission
- **WHEN** a visitor submits `POST /api/messages` with `{"name": "test", "email": "a@b.com", "url": "", "content": "Hello"}`
- **THEN** the system inserts a new row into `sl_message` with the provided fields plus the visitor's IP
- **AND** returns HTTP 200 with the created message

#### Scenario: Missing email
- **WHEN** a visitor submits `POST /api/messages` with `{"content": "Hello"}` and no `email` field
- **THEN** the system returns HTTP 400 with `{"code": 400, "message": "请输入邮箱和留言内容"}`

#### Scenario: Empty content
- **WHEN** a visitor submits `POST /api/messages` with `{"email": "a@b.com", "content": ""}`
- **THEN** the system returns HTTP 400 with `{"code": 400, "message": "请输入邮箱和留言内容"}`

### Requirement: Message page with form and list
The blog SPA SHALL provide a `/messages` route displaying a submission form (name, email, url, content fields) and a chronological list of all messages.

#### Scenario: View messages page
- **WHEN** a visitor navigates to `/messages`
- **THEN** the page shows a message submission form and a list of existing messages below it

#### Scenario: Submit message from page
- **WHEN** a visitor fills in the form and clicks submit
- **THEN** the form data is sent to `POST /api/messages`
- **AND** the submit button is disabled during the request
- **AND** on success, the message list refreshes and the content field is cleared
- **AND** the submit button is re-enabled after the request completes

### Requirement: Messages link in navigation
The blog SPA navigation SHALL include a "留言" link that routes to `/messages`.

#### Scenario: Navigation includes messages
- **WHEN** any page is loaded
- **THEN** the top navigation bar shows a "留言" link
- **AND** clicking it navigates to `/messages`
