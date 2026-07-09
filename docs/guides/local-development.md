# 本地开发指南

当前主线：

- Go API：`backend/`
- Next.js 前端：`frontend/`
- 本地依赖：PostgreSQL、Redis、MinIO、mock image provider

## 环境要求

- Node.js >= 20
- pnpm 11
- Docker 和 Docker Compose
- Go 1.25+（开发后端需要）

## 初始化

```bash
git clone https://github.com/sine-io/labhaus.git
cd labhaus
cd frontend
pnpm install
cp .env.example .env.local
cd ..
cp backend/.env.example backend/.env
```

`frontend/.env.local`：

```bash
BACKEND_URL=http://localhost:8080
```

## 启动服务

### 一次性启动本地演示环境

```bash
docker compose -f infra/docker-compose.yml up -d --build
docker compose -f infra/docker-compose.yml exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql
docker compose -f infra/docker-compose.yml restart api
(cd frontend && pnpm dev)
```

### 分开启动

只启动依赖服务：

```bash
docker compose -f infra/docker-compose.yml up -d postgres redis minio mock-image-provider
docker compose -f infra/docker-compose.yml exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql
```

裸跑 Go API：

```bash
cd backend
go run cmd/api/main.go
```

启动前端：

```bash
(cd frontend && pnpm dev)
```

## 常用命令

```bash
# 前端开发
(cd frontend && pnpm dev)

# 前端测试
(cd frontend && pnpm test)
(cd frontend && pnpm typecheck)
(cd frontend && pnpm lint)

# 前端构建和格式化
(cd frontend && pnpm build)
(cd frontend && pnpm format)
(cd frontend && pnpm format:check)

# 后端
cd backend
go test ./...
go build ./...
go run cmd/api/main.go

# 本地 smoke
infra/scripts/mvp-smoke.sh
```

## Docker 服务

### PostgreSQL

- Host: `localhost`
- Port: `5432`
- Database: `labhaus`
- User: `labhaus`
- Password: `labhaus_dev_password`

### Redis

- Host: `localhost`
- Port: `6379`

### MinIO

- API: `http://localhost:9000`
- Console: `http://localhost:9001`
- User: `minioadmin`
- Password: `minioadmin`

### Mock Image Provider

- URL: `http://localhost:8089`
- Health: `GET /health`
- Generate: `POST /v1/generate`
- API key: `dev-mock-key`

## 当前开发注意事项

- `backend/seeds/styles.sql` 当前只有 12 条 demo 样式；导入后需要重启 API。
- `LABHAUS_IMAGE_PROVIDER_BASE_URL` 和 `LABHAUS_IMAGE_PROVIDER_API_KEY` 对裸跑 API 是必填。
- Web 前端通过 Next.js route handler 代理到 Go API，不直接在浏览器读取 Go API URL。
- 登录 token 存在浏览器 `localStorage` 的 `labhaus_bearer_token`。
- Workflow API 现在只管理元数据和状态，不执行完整视频工作流。

## 代码规范

- TypeScript 使用严格模式。
- 前端使用 ESLint 和 Prettier。
- Go 代码使用 `gofmt`。
- 新功能需要补测试；文档变更需要同步 quick-start / API / README 中相关内容。

## 故障排查

### 样式推荐结果为空

```bash
docker compose -f infra/docker-compose.yml exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql
docker compose -f infra/docker-compose.yml restart api
```

### API 启动失败，提示 Provider 配置缺失

裸跑 API 时确认已设置：

```bash
export LABHAUS_IMAGE_PROVIDER_BASE_URL=http://localhost:8089
export LABHAUS_IMAGE_PROVIDER_API_KEY=dev-mock-key
```

### Go 测试缺依赖服务

存储和集成测试可能需要 PostgreSQL / MinIO 正在运行。先运行：

```bash
docker compose -f infra/docker-compose.yml up -d postgres minio
```

### pnpm 安装失败

```bash
pnpm store prune
rm -rf frontend/node_modules
(cd frontend && pnpm install)
```
