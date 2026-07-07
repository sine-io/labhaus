# Labhaus 快速开始指南

## 前置要求

- **Docker Desktop** (macOS/Windows) 或 Docker + Docker Compose (Linux)
- **Git**
- **Node.js** 20+ 和 pnpm 9+ (仅开发模式需要)

## 1. 克隆项目

```bash
git clone https://github.com/sine-io/labhaus.git
cd labhaus
```

## 2. 环境配置

复制环境变量模板：

```bash
cp backend/.env.example backend/.env
```

编辑 `backend/.env`，配置必要的环境变量：

```bash
# Go API 配置
LABHAUS_SERVER_PORT=8080
LABHAUS_JWT_SECRET_KEY=your-secret-key-change-in-production

# 图像 Provider（必须显式配置）
LABHAUS_IMAGE_PROVIDER_BASE_URL=http://localhost:8089
LABHAUS_IMAGE_PROVIDER_API_KEY=replace-with-provider-api-key
```

Docker Compose 也会读取 `LABHAUS_IMAGE_PROVIDER_BASE_URL` 和 `LABHAUS_IMAGE_PROVIDER_API_KEY`。启动前请在 shell 中导出，或放入 Compose 可读取的 `.env` 文件。

## 3. 启动服务

### 方式 A: Docker Compose（推荐）

```bash
# 启动所有服务
docker compose up -d

# 查看日志
docker compose logs -f api

# 停止服务
docker compose down
```

访问：

- API: http://localhost:8080
- Web: http://localhost:3000（单独运行 `apps/web`）
- PostgreSQL: localhost:5432
- Redis: localhost:6379
- MinIO: http://localhost:9001

### 方式 B: 本地开发模式

```bash
# 安装依赖
pnpm install

# 启动基础设施
docker compose up -d postgres redis minio

# 启动 Go API
cd backend
go run cmd/api/main.go

# 启动前端
cd ../apps/web
pnpm dev
```

## 4. 验证安装

### 健康检查

```bash
curl http://localhost:8080/api/health
```

预期响应：

```json
{
  "status": "healthy",
  "version": "0.1.0"
}
```

### 注册并登录

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

### 测试样式库 API

```bash
# 获取样式列表
curl "http://localhost:8080/api/styles?limit=5" \
  -H "Authorization: Bearer $TOKEN"

# 样式推荐
curl -X POST http://localhost:8080/api/styles/recommend \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"query": "modern UI design", "limit": 5}'
```

## 5. 用户注册和认证

### 注册账号

前端访问 `http://localhost:3000/auth`，登录或注册后会保存 Bearer Token。样式推荐和批量生图页面会自动附带 Authorization。

### 使用认证

```bash
# 使用 token 访问受保护的端点
curl http://localhost:8080/api/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 6. 运行测试

```bash
cd packages/workflow

# 工作流参考包测试
pnpm test

# 前端类型检查
cd ../../apps/web
pnpm typecheck
```

## 7. 下一步

- 📖 阅读 [本地开发指南](local-development.md) 了解开发流程
- 🏗️ 查看 [后端文档](../../backend/README.md)
- 🔐 查看 Go 后端用户接口：`/api/users/register`、`/api/users/login`、`/api/users/me`
- 📦 查看 [部署指南](../DEPLOYMENT.md)

## 常见问题

### Q: Docker 容器启动失败

**A**: 检查端口占用：

```bash
# 检查 8080 端口
lsof -i :8080

# 检查 5432 端口（PostgreSQL）
lsof -i :5432
```

### Q: 数据库连接失败

**A**: 确保 PostgreSQL 容器正在运行：

```bash
docker compose ps postgres
docker compose logs postgres
```

### Q: 样式库数据为空

**A**: 当前 Go 主线需要补齐样式数据导入流程。`apps/api/src/scripts/import-styles.ts` 仅作为 legacy 参考，不再作为主线启动步骤。

### Q: pnpm 安装依赖慢

**A**: 配置国内镜像：

```bash
pnpm config set registry https://registry.npmmirror.com
```

## 获取帮助

- 📋 [GitHub Issues](https://github.com/sine-io/labhaus/issues)
- 💬 [Discussions](https://github.com/sine-io/labhaus/discussions)
- 📧 Email: support@labhaus.io
