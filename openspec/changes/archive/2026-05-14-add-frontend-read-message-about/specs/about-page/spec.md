## ADDED Requirements

### Requirement: About Me page accessible at /about
The blog SPA SHALL provide an `/about` route displaying a centered page with the WeChat QR code image, a personal bio, and social links.

#### Scenario: View About Me page
- **WHEN** a visitor navigates to `/about`
- **THEN** the page displays:
- **AND** a WeChat QR code image loaded from `<cdn_url>/about_me.jpg`
- **AND** the text "一个喜欢踢足球的游戏开发程序员。"
- **AND** a GitHub link to `https://github.com/iamzcr`
- **AND** a blog link to `http://blog.iamzcr.com/`
- **AND** all content is horizontally centered on the page

### Requirement: About Me link in navigation
The blog SPA navigation SHALL include a "关于我" link that routes to `/about`.

#### Scenario: Navigation includes about
- **WHEN** any page is loaded
- **THEN** the top navigation bar shows a "关于我" link
- **AND** clicking it navigates to `/about`
