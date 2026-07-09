# Labhaus 技术栈

**最后更新**：2026-07-08

## 当前主线

- `backend/`：Go API 服务。
- `frontend/`：Next.js App Router 前端。
- `infra/docker-compose.yml`：本地 PostgreSQL、Redis、MinIO、mock image provider 和 Go API。

当前 MVP 是“认证后的样式推荐 + 批量生图”。视频生成、TTS、FFmpeg、可视化编辑器和模板市场属于后续阶段。

## 后端

### 语言和框架

- Go 1.25
- Gin
- GORM
- Viper
- Zerolog
- testify

### 数据和基础设施

- PostgreSQL 16
- Redis 7
- MinIO / S3
- Docker / Docker Compose

### 认证

- `github.com/golang-jwt/jwt/v5`
- JWT Bearer Token
- bcrypt 密码哈希

当前没有 refresh token、OAuth、RBAC、配额和限流。

### 图像生成

后端通过 Provider 合约调用外部服务：

```http
POST /v1/generate
Authorization: Bearer <api-key>
```

本地 demo 使用 `backend/cmd/mock-image-provider`。

### 队列

当前实现为 Redis list/hash 轻量队列：

- pending / processing / completed / dead 状态
- retry count
- dead letter queue

当前 workflow worker 只保留执行骨架；完整视频工作流执行待接入。

## 前端

- Next.js 16
- React 19
- TypeScript
- Tailwind CSS 4
- lucide-react
- Node built-in test runner

当前没有 shadcn/ui、Zustand、React Flow 或 SSE；这些可以在后续视频和可视化阶段引入。

## 当前运行架构

```text
Browser
  |
  v
frontend (Next.js)
  |  route handlers forward Authorization
  v
backend Go API (Gin)
  |-- PostgreSQL: users/styles/workflows
  |-- Redis: queue skeleton
  |-- MinIO: generated images and future media assets
  |-- Image Provider: mock or real HTTP provider
```

## 项目结构

```text
labhaus/
├── frontend/
├── backend/
│   ├── cmd/
│   │   ├── api/
│   │   └── mock-image-provider/
│   ├── internal/
│   │   ├── application/
│   │   ├── domain/
│   │   └── infrastructure/
│   ├── seeds/
│   └── tests/
├── docs/
└── infra/
    ├── docker-compose.yml
    └── scripts/
```

## Go 依赖摘要

实际版本以 `backend/go.mod` 为准。核心直接依赖包括：

- `github.com/gin-gonic/gin`
- `github.com/golang-jwt/jwt/v5`
- `github.com/google/uuid`
- `github.com/redis/go-redis/v9`
- `github.com/rs/zerolog`
- `github.com/spf13/viper`
- `github.com/stretchr/testify`
- `golang.org/x/crypto`
- `gorm.io/driver/postgres`
- `gorm.io/gorm`

## Node 依赖摘要

实际版本以 `frontend/package.json` 和 `frontend/pnpm-lock.yaml` 为准。

`frontend`：

- Next.js
- React / React DOM
- Tailwind CSS
- lucide-react
- TypeScript
- ESLint
- Prettier

## 后续技术方向

### 视频工作流

- LLM 剧本和分镜生成
- TTS 配音
- FFmpeg 合成
- 任务监控和中间产物预览
- 更完整的 Redis queue worker

### 推荐算法

当前 HTTP 运行时使用关键词重叠打分；后续应接入 `internal/domain/style/recommendation` 中已有的 TF-IDF + Cosine 实现，并结合 500+ 样式库数据评估推荐质量。

### 可视化编辑器

- React Flow
- 工作流 JSON schema
- 节点注册表
- 输入、LLM、生图、TTS、视频合成、输出节点

## 参考文档

- `docs/architecture/system-design.md`
- `docs/architecture/api-design.md`
- `docs/architecture/GO_DDD_ARCHITECTURE.md`
- `docs/guides/quick-start.md`
- `docs/DEPLOYMENT.md`
