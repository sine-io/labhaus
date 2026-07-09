# Labhaus 快速开始指南

本指南用于从干净 checkout 跑通当前 MVP：

> 注册/登录 -> 样式推荐 -> 批量生图 -> MinIO 预签名图片链接

## 前置要求

- Docker 和 Docker Compose
- Node.js 20+
- pnpm 11（前端工程当前使用 `pnpm@11.8.0`）
- `curl` 和 `jq`（运行 smoke 脚本需要）
- Go 1.25+（裸跑 Go API 或运行 Go 测试需要）

## 1. 克隆项目

```bash
git clone https://github.com/sine-io/labhaus.git
cd labhaus
```

## 2. 启动后端依赖和 API

```bash
docker compose -f infra/docker-compose.yml up -d --build
```

这会启动：

- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`
- MinIO: `localhost:9000`
- MinIO Console: `http://localhost:9001`
- Mock Image Provider: `http://localhost:8089`
- Go API: `http://localhost:8080`

Docker Compose 会自动给 API 注入本地 demo Provider：

```bash
LABHAUS_IMAGE_PROVIDER_BASE_URL=http://mock-image-provider:8089
LABHAUS_IMAGE_PROVIDER_API_KEY=dev-mock-key
```

## 3. 导入样式 seed

当前本地 demo seed 包含 12 条样式。后续产品目标是接入 500+ 样式库，但当前仓库内 seed 不是完整样式库。

```bash
docker compose -f infra/docker-compose.yml exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql
docker compose -f infra/docker-compose.yml restart api
```

需要重启 API 的原因：样式推荐器在 API 启动时加载样式快照。

## 4. 启动前端

```bash
cd frontend
pnpm install
cp .env.example .env.local
pnpm dev
```

`frontend/.env.local` 当前只需要：

```bash
BACKEND_URL=http://localhost:8080
```

访问 http://localhost:3000。

## 5. 用 Web 跑 MVP

1. 进入 `/auth` 注册或登录。
2. 登录成功后 Token 会保存在浏览器 `localStorage`。
3. 进入 `/styles/recommend` 输入创意描述，例如 `modern UI dashboard for a SaaS product`。
4. 复制推荐样式中的 prompt 或直接参考其风格。
5. 进入 `/images/generate`，每行输入一个 prompt。
6. 点击生成，等待图片列表返回。
7. 点击下载，打开 MinIO 预签名 URL。

## 6. 用脚本验证 MVP

```bash
infra/scripts/mvp-smoke.sh
```

脚本会：

1. 检查 `/api/health`。
2. 注册 demo 用户（已存在时允许继续）。
3. 登录并获取 token。
4. 调用 `/api/styles/recommend`。
5. 调用 `/api/images/generate`。
6. 检查推荐结果和生成图片结果非空。

可覆盖默认变量：

```bash
API_URL=http://localhost:8080 \
EMAIL=demo@example.com \
PASSWORD=SecurePassword123! \
NAME="Demo User" \
infra/scripts/mvp-smoke.sh
```

## 7. 常用 API 手动验证

### 健康检查

```bash
curl http://localhost:8080/api/health
```

预期：

```json
{
  "status": "healthy",
  "version": "0.1.0"
}
```

### 注册和登录

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@labhaus.io",
    "password": "SecurePassword123!",
    "name": "Demo User"
  }'

TOKEN=$(curl -s -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@labhaus.io",
    "password": "SecurePassword123!"
  }' | jq -r .token)
```

### 样式推荐

```bash
curl -X POST http://localhost:8080/api/styles/recommend \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"query": "modern UI dashboard", "limit": 5}'
```

### 批量生图

```bash
curl -X POST http://localhost:8080/api/images/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "prompts": ["modern dashboard hero", "minimal product card"],
    "width": 512,
    "height": 512,
    "quality": "standard"
  }'
```

## 8. 运行测试

```bash
# 前端 helper 测试
(cd frontend && pnpm test)

# 前端类型检查
(cd frontend && pnpm typecheck)

# Go 测试（需要 Go 1.25+）
cd backend
go test ./...
```

## 常见问题

### 样式推荐为空

确认已导入 seed 并重启 API：

```bash
docker compose -f infra/docker-compose.yml exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql
docker compose -f infra/docker-compose.yml restart api
```

### 生图失败

检查 mock provider 和 API 日志：

```bash
docker compose -f infra/docker-compose.yml ps mock-image-provider api
docker compose -f infra/docker-compose.yml logs mock-image-provider api
```

### Web 请求返回 401

进入 `/auth` 重新登录，或清除 Token 后再登录。

### 端口被占用

检查常用端口：

```bash
lsof -i :3000
lsof -i :8080
lsof -i :5432
lsof -i :9000
```

## 下一步

- 本地开发：`docs/guides/local-development.md`
- API 契约：`docs/architecture/api-design.md`
- 部署：`docs/DEPLOYMENT.md`
