# 个人博客系统

基于 Go + Gin + GORM 的后端 API，配合 Vite + Naive UI 前端构建的个人博客系统。

## 技术栈

- **后端**: Go, Gin, GORM, MySQL
- **前端**: Vue 3, TypeScript, Vite, Naive UI
- **数据库**: MySQL (stacklifes)

## 项目结构

```
iamzcr/
├── cmd/
│   ├── admin/          # 后台管理 API 服务 (端口 8081)
│   └── frontend/       # 前台展示 API 服务 (端口 8082)
├── config/             # 配置文件
├── handlers/           # API 处理器
├── models/             # 数据模型
├── web/
│   ├── admin/          # 后台管理前端
│   └── frontend/       # 前台展示前端
├── .env                # 环境配置
├── .env.eg             # 环境配置示例
├── admin.exe           # 后台服务可执行文件
└── frontend.exe        # 前台服务可执行文件
```

## 配置

在项目根目录创建 `.env` 文件，可以从 `.env.eg` 复制：

```bash
cp .env.eg .env
```

配置项说明：

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| DB_HOST | 数据库地址 | localhost |
| DB_PORT | 数据库端口 | 3306 |
| DB_USER | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | root |
| DB_NAME | 数据库名称 | stacklifes |
| ADMIN_PORT | 后台 API 端口 | 8081 |
| FRONTEND_PORT | 前台 API 端口 | 8082 |

## 快速开始

### 1. 启动后端服务

```bash
# 启动后台后端 (8081)
./admin.exe

# 启动前台后端 (8082)
./frontend.exe
```

### 2. 启动前端开发服务器

```bash
# 前台前端 (默认 3000 端口)
cd web/frontend
npm run dev

# 后台前端 (默认 3001 端口)
cd web/admin
npm run dev
```

### 3. 访问

- 前台展示: http://localhost:3000/
- 后台管理: http://localhost:3001/

## API 接口

### 前台 API (8082)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/articles | 获取文章列表 |
| GET | /api/articles/:id | 获取文章详情 |

### 后台 API (8081)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/articles | 获取文章列表 |
| GET | /api/articles/:id | 获取文章详情 |
| POST | /api/articles | 创建文章 |
| PUT | /api/articles/:id | 更新文章 |
| DELETE | /api/articles/:id | 删除文章 |
| GET | /api/categories | 获取分类列表 |
| POST | /api/categories | 创建分类 |
| PUT | /api/categories/:id | 更新分类 |
| DELETE | /api/categories/:id | 删除分类 |
| GET | /api/comments | 获取评论列表 |
| POST | /api/comments | 创建评论 |
| PUT | /api/comments/:id | 更新评论 |
| DELETE | /api/comments/:id | 删除评论 |
| GET | /api/menus | 获取菜单列表 |
| GET | /api/tags | 获取标签列表 |
| GET | /api/website | 获取网站设置 |
| PUT | /api/website | 更新网站设置 |

## 数据库表

系统使用 `stacklifes` 数据库中的现有表：

- `sl_article` - 文章
- `sl_category` - 分类
- `sl_comment` - 评论
- `sl_menu` - 菜单
- `sl_tags` - 标签
- `sl_website` - 网站设置
- `sl_admin` - 管理员
- `sl_admin_group` - 管理员组

## 构建生产版本

```bash
# 前端构建
cd web/frontend
npm run build

cd web/admin
npm run build

# 后端编译
go build -o admin.exe ./cmd/admin
go build -o frontend.exe ./cmd/frontend
```