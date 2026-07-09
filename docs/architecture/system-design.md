# Labhaus 系统设计

**最后更新**：2026-07-08

## 设计目标

Labhaus 的长期目标是 AI 视频内容工作流实验室。当前系统先实现图像素材 MVP：

```text
认证 -> 样式推荐 -> 批量生图 -> MinIO 下载链接
```

这个 MVP 是后续视频工作流的基础。

## 当前架构

```text
Browser
  |
  | HTTP
  v
frontend (Next.js App Router)
  |
  | Route handlers forward Authorization
  v
Go API (Gin)
  |
  |-- PostgreSQL: users, styles, workflows
  |-- Redis: queue skeleton
  |-- MinIO: generated images and future media
  |-- Image Provider: mock or real HTTP provider
```

## 模块

### 1. Web 前端

位置：`frontend`

职责：

- 注册/登录页面。
- 保存 Bearer Token。
- 样式推荐页面。
- 批量生图页面。
- 用 Next.js route handlers 代理 Go API 请求。

当前页面：

- `/`
- `/auth`
- `/styles/recommend`
- `/images/generate`

### 2. API 层

位置：`backend/internal/infrastructure/http`

职责：

- Gin 路由。
- 请求日志和 recovery。
- JWT 鉴权 middleware。
- HTTP handlers。

当前路由：

- `/api/health`
- `/api/users/*`
- `/api/styles/*`
- `/api/workflows/*`
- `/api/images/*`

### 3. 应用层

位置：`backend/internal/application`

职责：

- CQRS command/query handlers。
- DTO 转换。
- 批量生图服务。

批量生图服务通过 semaphore 限制并发，默认最大并发为 10。

### 4. 领域层

位置：`backend/internal/domain`

当前聚合和模型：

- `user`
- `style`
- `workflow`
- `image/provider`

状态机：

```text
DRAFT -> PENDING | CANCELLED
PENDING -> RUNNING | CANCELLED
RUNNING -> PAUSED | COMPLETED | FAILED
PAUSED -> RUNNING | CANCELLED
COMPLETED / FAILED / CANCELLED -> terminal
```

### 5. 推荐模块

当前存在两套实现：

1. HTTP 运行时：`handlers.StaticStyleRecommender`，使用关键词重叠打分。
2. 领域模块：`internal/domain/style/recommendation`，已有 TF-IDF + Cosine 实现。

后续应将 HTTP 运行时推荐器替换为领域层 TF-IDF + Cosine，并基于 500+ 样式库评估推荐质量。

### 6. 图像 Provider

位置：

- `backend/internal/domain/image/provider`
- `backend/internal/infrastructure/image/gptimage2`
- `backend/cmd/mock-image-provider`

Provider 合约：

```http
POST /v1/generate
Authorization: Bearer <api-key>
Content-Type: application/json
```

Go API 下载 Provider 返回的 `image_url`，再上传到 MinIO `images` bucket。

### 7. 存储

位置：`backend/internal/infrastructure/storage`

- `MinIOStorage`：通用 bucket/object 操作。
- `MinIOImageStorage`：图片上传、下载、存在性检查、预签名 URL。

API 启动时确保 buckets：

- `workflows`
- `images`
- `videos`
- `temp`

### 8. 队列

位置：`backend/internal/infrastructure/queue`

当前 Redis queue 支持：

- enqueue/dequeue
- worker loop
- retry
- dead letter queue
- task status

当前 worker 仅保留 workflow task handler 骨架；视频工作流执行逻辑待实现。

## 数据模型

### styles

```text
id           varchar(36) primary key
name         varchar(100) not null
description  varchar(500)
prompt       text not null
category     varchar(50)
tags         text          # JSON array string
created_at   timestamptz not null
updated_at   timestamptz not null
deleted_at   timestamptz
```

### users

```text
id            varchar(36) primary key
email         varchar(255) not null unique
password_hash varchar(255) not null
name          varchar(100) not null
role          varchar(20) not null default 'user'
created_at    timestamptz not null
updated_at    timestamptz not null
deleted_at    timestamptz
```

### workflows

```text
id         varchar(36) primary key
user_id    varchar(36) not null
style_id   varchar(36) not null
state      varchar(20) not null
config     jsonb not null
result     jsonb
created_at timestamptz not null
updated_at timestamptz not null
deleted_at timestamptz
```

## 当前 MVP 数据流

### 登录

```text
Browser -> Next /api/users/login -> Go /api/users/login
Go -> PostgreSQL users
Go -> JWT token
Next -> Browser
Browser -> localStorage
```

### 样式推荐

```text
Browser -> Next /api/styles/recommend
Next forwards Authorization
Go AuthMiddleware validates JWT
Go recommender scores loaded style snapshot
Go returns recommendations
```

### 批量生图

```text
Browser -> Next /api/images/generate
Next forwards Authorization
Go AuthMiddleware validates JWT
Go BatchImageService runs prompts concurrently
Go -> Image Provider /v1/generate
Go downloads provider image_url
Go uploads bytes to MinIO images bucket
Go returns presigned URLs
```

## 配置边界

关键环境变量：

- `LABHAUS_DATABASE_*`
- `LABHAUS_REDIS_*`
- `LABHAUS_MINIO_*`
- `LABHAUS_JWT_SECRET_KEY`
- `LABHAUS_JWT_TOKEN_DURATION`
- `LABHAUS_IMAGE_PROVIDER_BASE_URL`
- `LABHAUS_IMAGE_PROVIDER_API_KEY`

图像 Provider 配置缺失时 API 启动失败，避免误连示例地址。

## 当前限制

- 只有本地 demo seed，尚未导入 500+ 样式库。
- 样式推荐运行时未接入 TF-IDF + Cosine。
- Workflow API 尚未执行真实工作流。
- 图片 progress 接口当前固定返回 completed。
- 没有 refresh token、OAuth、RBAC、限流、OpenAPI 自动生成。
- Docker Compose 未包含 Web 前端服务。

## 后续架构演进

### 视频工作流

```text
Text/Article
  -> Script generation
  -> Storyboard
  -> Style recommendation
  -> Batch image generation
  -> TTS
  -> FFmpeg render
  -> MP4 in MinIO
```

需要补充：

- 任务持久化和执行日志。
- Redis queue worker 接入真实执行器。
- 中间产物存储。
- 任务进度 API。
- 任务监控 UI。

### 可视化工作流

后续可引入：

- React Flow。
- 工作流 JSON schema。
- 节点注册表。
- DAG 校验。
- 输入/处理/输出节点。

### 生产化

后续需要：

- OpenAPI 文档。
- API 限流。
- 审计日志。
- 监控和告警。
- 备份策略。
- 多环境配置管理。
