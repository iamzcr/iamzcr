# 个人博客系统

基于 Go + Gin + GORM 的后端 API，配合 Vite + Vue 3 + Naive UI 前端构建的个人博客系统。

## 技术栈

- **后端**: Go, Gin, GORM, MySQL
- **前端**: Vue 3, TypeScript, Vite, Naive UI
- **数据库**: MySQL (stacklifes)
- **部署**: Nginx

## 项目结构

```
iamzcr/
├── cmd/
│   ├── admin/          # 后台管理 API 服务 (端口 8081)
│   └── frontend/       # 前台展示 API 服务 (端口 8082)
├── config/             # 配置文件
├── handlers/           # API 处理器
│   ├── admin/          # 后台管理接口
│   └── frontend/       # 前台展示接口
├── middleware/         # 中间件 (CORS, Auth, JWT)
├── models/             # 数据模型
├── services/           # 业务逻辑层
├── sql/                # SQL 脚本
├── web/
│   ├── admin/          # 后台管理前端 (Vue 3 + NaiveUI)
│   └── frontend/       # 前台展示前端 (Vue 3 + NaiveUI)
├── .env                # 环境配置
├── .env.eg             # 环境配置示例
├── Makefile            # 构建脚本
├── admin               # 后台服务可执行文件 (Linux)
├── admin.exe           # 后台服务可执行文件 (Windows)
├── frontend           # 前台服务可执行文件 (Linux)
└── frontend.exe        # 前台服务可执行文件 (Windows)
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

### 1. 初始化数据库

导入 SQL 文件到 MySQL：

```bash
mysql -u root -p stacklifes < sql/stacklifes.sql
```

### 2. 启动后端服务

```bash
# Linux/Mac
make run-admin
make run-frontend

# Windows
./admin.exe
./frontend.exe
```

### 3. 启动前端开发服务器

```bash
# 前台前端 (默认 3000 端口)
cd web/frontend
npm install
npm run dev

# 后台前端 (默认 3001 端口)
cd web/admin
npm install
npm run dev
```

### 4. 访问

- 前台展示: http://localhost:3000/
- 后台管理: http://localhost:3001/
- 后台 API: http://localhost:8081/
- 前台 API: http://localhost:8082/

## 生产部署

### 1. 构建后端

```bash
# 使用 Makefile
make build

# 或手动编译
go build -o admin ./cmd/admin
go build -o frontend ./cmd/frontend
```

### 2. 构建前端

```bash
# 前台前端
cd web/frontend
npm install
npm run build

# 后台前端
cd web/admin
npm install
npm run build
```

构建产物位于 `web/frontend/dist` 和 `web/admin/dist`

### 3. Nginx 配置

将前端构建产物部署到 Nginx，配置示例：

```nginx
# /etc/nginx/conf.d/iamzcr.conf

upstream admin_api {
    server 127.0.0.1:8081;
}

upstream frontend_api {
    server 127.0.0.1:8082;
}

server {
    listen 80;
    server_name iamzcr.com www.iamzcr.com;

    # 前台静态文件
    location / {
        root /var/www/iamzcr/frontend/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
        
        # 缓存静态资源
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 30d;
            add_header Cache-Control "public, immutable";
        }
    }

    # 前台 API 代理
    location /api/ {
        proxy_pass http://frontend_api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 后台管理静态文件
    location /admin/ {
        alias /var/www/iamzcr/admin/dist/;
        index index.html;
        try_files $uri $uri/ /admin/index.html;
    }

    # 后台 API 代理
    location /admin-api/ {
        rewrite ^/admin-api/(.*)$ /$1 break;
        proxy_pass http://admin_api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 上传文件大小限制
    client_max_body_size 50M;
}
```

如果需要让后台 API 通过 `/api/` 统一访问：

```nginx
server {
    listen 80;
    server_name iamzcr.com www.iamzcr.com;

    # 前台静态文件
    location / {
        root /var/www/iamzcr/frontend/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    # 后台管理静态文件
    location /admin {
        alias /var/www/iamzcr/admin/dist;
        index index.html;
        try_files $uri $uri/ /admin/index.html;
    }

    # 所有 API 代理 (后台 + 前台)
    location /api/ {
        # 根据路径区分前后台 API
        location ~ ^/api/admin(/.*)$ {
            proxy_pass http://admin_api$1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
        
        location ~ ^/api(/.*)$ {
            proxy_pass http://frontend_api$1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
    }
}
```

### 4. 部署目录结构

```
/var/www/iamzcr/
├── frontend/
│   └── dist/          # 前台构建产物
├── admin/
│   └── dist/          # 后台构建产物
├── logs/              # 日志目录
└── data/              # 数据目录 (如需)
```

### 5. Systemd 服务配置 (Linux)

创建 `/etc/systemd/system/iamzcr-admin.service`:

```ini
[Unit]
Description=IAMZCR Admin API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/iamzcr
ExecStart=/var/www/iamzcr/admin
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

创建 `/etc/systemd/system/iamzcr-frontend.service`:

```ini
[Unit]
Description=IAMZCR Frontend API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/iamzcr
ExecStart=/var/www/iamzcr/frontend
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable iamzcr-admin iamzcr-frontend
sudo systemctl start iamzcr-admin iamzcr-frontend
```

## API 接口

### 后台 API (8081)

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/login | 登录 |
| POST | /api/logout | 登出 |
| GET | /api/admin/info | 获取管理员信息 |
| GET | /api/articles | 获取文章列表 |
| GET | /api/articles/:id | 获取文章详情 |
| POST | /api/articles | 创建文章 |
| PUT | /api/articles/:id | 更新文章 |
| DELETE | /api/articles/:id | 删除文章 |
| GET | /api/categories | 获取分类列表 |
| POST | /api/categories | 创建分类 |
| PUT | /api/categories/:id | 更新分类 |
| DELETE | /api/categories/:id | 删除分类 |
| GET | /api/directories | 获取目录列表 |
| POST | /api/directories | 创建目录 |
| PUT | /api/directories/:id | 更新目录 |
| DELETE | /api/directories/:id | 删除目录 |
| GET | /api/tags | 获取标签列表 |
| POST | /api/tags | 创建标签 |
| PUT | /api/tags/:id | 更新标签 |
| DELETE | /api/tags/:id | 删除标签 |
| GET | /api/comments | 获取评论列表 |
| POST | /api/comments | 创建评论 |
| PUT | /api/comments/:id | 更新评论 |
| DELETE | /api/comments/:id | 删除评论 |
| GET | /api/menus | 获取菜单列表 |
| POST | /api/menus | 创建菜单 |
| PUT | /api/menus/:id | 更新菜单 |
| DELETE | /api/menus/:id | 删除菜单 |
| GET | /api/website | 获取网站设置 |
| PUT | /api/website | 更新网站设置 |
| GET | /api/attaches | 获取附件列表 |
| POST | /api/attaches | 创建附件 |
| PUT | /api/attaches/:id | 更新附件 |
| DELETE | /api/attaches/:id | 删除附件 |
| GET | /api/langs | 获取语言列表 |
| POST | /api/langs | 创建语言 |
| PUT | /api/langs/:id | 更新语言 |
| DELETE | /api/langs/:id | 删除语言 |
| GET | /api/logs | 获取日志列表 |
| DELETE | /api/logs/:id | 删除日志 |
| GET | /api/messages | 获取留言列表 |
| POST | /api/messages | 创建留言 |
| PUT | /api/messages/:id | 更新留言 |
| DELETE | /api/messages/:id | 删除留言 |
| GET | /api/permits | 获取权限列表 |
| POST | /api/permits | 创建权限 |
| PUT | /api/permits/:id | 更新权限 |
| DELETE | /api/permits/:id | 删除权限 |
| GET | /api/reads | 获取阅读记录 |
| DELETE | /api/reads/:id | 删除阅读记录 |
| GET | /api/admins | 获取管理员列表 |
| POST | /api/admins | 创建管理员 |
| PUT | /api/admins/:id | 更新管理员 |
| DELETE | /api/admins/:id | 删除管理员 |
| POST | /api/admin/password | 修改密码 |
| GET | /api/admin_groups | 获取用户组列表 |
| POST | /api/admin_groups | 创建用户组 |
| PUT | /api/admin_groups/:id | 更新用户组 |
| DELETE | /api/admin_groups/:id | 删除用户组 |

### 前台 API (8082)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/articles | 获取文章列表 |
| GET | /api/articles/:id | 获取文章详情 |
| GET | /api/categories | 获取分类列表 |
| GET | /api/directories | 获取目录列表 |
| GET | /api/tags | 获取标签列表 |
| GET | /api/menus | 获取菜单列表 |
| GET | /api/website | 获取网站设置 |
| GET | /api/comments | 获取评论列表 |
| POST | /api/comments | 创建评论 |
| GET | /api/langs | 获取语言列表 |
| GET | /api/reads | 获取阅读记录 |
| POST | /api/reads | 创建阅读记录 |
| GET | /api/messages | 获取留言列表 |
| POST | /api/messages | 创建留言 |

## 数据库表

系统使用 `stacklifes` 数据库中的表：

| 表名 | 说明 |
|------|------|
| sl_article | 文章 |
| sl_article_tags | 文章标签关联 |
| sl_category | 分类 |
| sl_comment | 评论 |
| sl_directory | 目录 |
| sl_menu | 菜单 |
| sl_tags | 标签 |
| sl_website | 网站设置 |
| sl_admin | 管理员 |
| sl_admin_group | 管理员组 |
| sl_attach | 附件 |
| sl_lang | 语言 |
| sl_log | 日志 |
| sl_message | 留言 |
| sl_permit | 权限 |
| sl_read | 阅读记录 |

## Makefile 命令

```bash
make help          # 显示帮助信息
make build         # 编译后端服务
make build-admin   # 编译后台服务
make build-frontend # 编译前台服务
make run-admin     # 运行后台服务
make run-frontend  # 运行前台服务
make clean         # 清理构建产物
make all           # 编译所有服务
```

## 默认账号

- 用户名: admin
- 密码: admin
