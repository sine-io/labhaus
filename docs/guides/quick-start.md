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
cp .env.example .env.local
```

编辑 `.env.local`，配置必要的环境变量：

```bash
# 数据库
DATABASE_URL=postgresql://labhaus:labhaus@postgres:5432/labhaus

# JWT 密钥（生产环境请使用强随机字符串）
JWT_SECRET=your-secret-key-change-in-production
JWT_REFRESH_SECRET=your-refresh-secret-key

# API 配置
API_PORT=3001
NODE_ENV=development

# CORS（前端地址）
CORS_ORIGIN=http://localhost:3000
```

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
- API: http://localhost:3001
- PostgreSQL: localhost:5432
- Redis: localhost:6379
- MinIO: http://localhost:9001

### 方式 B: 本地开发模式

```bash
# 安装依赖
pnpm install

# 启动数据库（仅 PostgreSQL）
docker compose up -d postgres

# 运行数据库迁移
cd apps/api
pnpm migrate

# 导入样式库数据（可选）
pnpm import-styles

# 启动 API 服务
pnpm dev
```

## 4. 验证安装

### 健康检查

```bash
curl http://localhost:3001/api/health
```

预期响应：
```json
{
  "status": "ok",
  "timestamp": "2026-06-19T..."
}
```

### 获取 API 信息

```bash
curl http://localhost:3001/api
```

### 测试样式库 API

```bash
# 获取样式列表
curl http://localhost:3001/api/styles?limit=5

# 搜索样式
curl http://localhost:3001/api/styles?search=portrait

# 样式推荐
curl -X POST http://localhost:3001/api/styles/recommend \
  -H "Content-Type: application/json" \
  -d '{"query": "modern UI design", "limit": 5}'
```

## 5. 用户注册和认证

### 注册账号

```bash
curl -X POST http://localhost:3001/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@labhaus.io",
    "password": "SecurePassword123!",
    "name": "Demo User"
  }'
```

保存返回的 `access_token`。

### 使用认证

```bash
# 使用 token 访问受保护的端点
curl http://localhost:3001/api/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 6. 运行测试

```bash
cd apps/api

# 单元测试
pnpm test

# 端到端测试（需要 API 服务运行）
pnpm test:e2e

# 类型检查
pnpm typecheck
```

## 7. 下一步

- 📖 阅读 [本地开发指南](local-development.md) 了解开发流程
- 🏗️ 查看 [API 设计文档](../../apps/api/docs/API_DESIGN.md)
- 🔐 了解 [认证系统](../../apps/api/docs/AUTHENTICATION.md)
- 📦 查看 [部署指南](../DEPLOYMENT.md)

## 常见问题

### Q: Docker 容器启动失败

**A**: 检查端口占用：
```bash
# 检查 3001 端口
lsof -i :3001

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

**A**: 运行导入脚本：
```bash
cd apps/api
pnpm import-styles
```

### Q: pnpm 安装依赖慢

**A**: 配置国内镜像：
```bash
pnpm config set registry https://registry.npmmirror.com
```

## 获取帮助

- 📋 [GitHub Issues](https://github.com/sine-io/labhaus/issues)
- 💬 [Discussions](https://github.com/sine-io/labhaus/discussions)
- 📧 Email: support@labhaus.io
