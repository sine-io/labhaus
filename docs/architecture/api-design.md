# Labhaus API 设计

**Base URL**：`http://localhost:8080/api`
**当前实现**：Go Gin API
**最后更新**：2026-07-08

## 通用规范

### 请求头

```http
Content-Type: application/json
Authorization: Bearer <token>
```

`Authorization` 只在受保护接口中需要。

### 错误响应

当前 Go API 使用简单错误格式：

```json
{
  "error": "human readable error"
}
```

### 鉴权

- 注册和登录是公开接口。
- `/api/users/me`、`/api/styles/*`、`/api/workflows/*`、`/api/images/*` 均需要 Bearer Token。
- 登录返回单个 JWT：`{ "token": "...", "user": {...} }`。
- 当前没有 refresh token、OAuth、RBAC 和配额。

## Health

### GET `/health`

响应：

```json
{
  "status": "healthy",
  "version": "0.1.0"
}
```

## Users

### POST `/users/register`

公开接口。

请求：

```json
{
  "email": "demo@labhaus.io",
  "password": "SecurePassword123!",
  "name": "Demo User"
}
```

约束：

- `email` 必填且必须为邮箱格式。
- `password` 必填，最少 8 字符。
- `name` 必填。

成功响应：`201 Created`

```json
{
  "id": "uuid",
  "email": "demo@labhaus.io",
  "name": "Demo User",
  "role": "user",
  "created_at": "2026-07-08T00:00:00Z",
  "updated_at": "2026-07-08T00:00:00Z"
}
```

重复邮箱：`409 Conflict`

```json
{
  "error": "email already exists"
}
```

### POST `/users/login`

公开接口。

请求：

```json
{
  "email": "demo@labhaus.io",
  "password": "SecurePassword123!"
}
```

成功响应：`200 OK`

```json
{
  "token": "jwt-token",
  "user": {
    "id": "uuid",
    "email": "demo@labhaus.io",
    "name": "Demo User",
    "role": "user",
    "created_at": "2026-07-08T00:00:00Z",
    "updated_at": "2026-07-08T00:00:00Z"
  }
}
```

### GET `/users/me`

需要 Bearer Token。

成功响应：`200 OK`

```json
{
  "id": "uuid",
  "email": "demo@labhaus.io",
  "name": "Demo User",
  "role": "user",
  "created_at": "2026-07-08T00:00:00Z",
  "updated_at": "2026-07-08T00:00:00Z"
}
```

### PATCH `/users/me`

需要 Bearer Token。

请求：

```json
{
  "email": "new-email@labhaus.io",
  "name": "New Name"
}
```

成功响应：更新后的 user DTO。

## Styles

所有 styles 接口均需要 Bearer Token。

### GET `/styles`

查询参数：

| 参数       | 类型         | 说明                                             |
| ---------- | ------------ | ------------------------------------------------ |
| `category` | string       | 按分类过滤                                       |
| `tags`     | string array | 按 tags 过滤，当前 repository 使用 LIKE contains |
| `limit`    | integer      | 默认 20                                          |
| `offset`   | integer      | 默认 0                                           |

成功响应：

```json
{
  "styles": [
    {
      "id": "style-id",
      "name": "Minimal UI Dashboard",
      "description": "Clean modern dashboard style.",
      "prompt": "modern minimal UI dashboard...",
      "category": "ui",
      "tags": ["ui", "dashboard"],
      "created_at": "2026-07-08T00:00:00Z",
      "updated_at": "2026-07-08T00:00:00Z"
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

### GET `/styles/:id`

成功响应：style DTO。

不存在：`404`

```json
{
  "error": "style not found"
}
```

### POST `/styles`

请求：

```json
{
  "name": "Cinematic Product",
  "description": "Premium product campaign style",
  "prompt": "cinematic product photography...",
  "category": "product",
  "tags": ["product", "cinematic"]
}
```

约束：

- `name` 必填，最多 100 字符。
- `description` 最多 500 字符。
- `prompt` 必填，最多 2000 字符。

成功响应：`201 Created`，返回 style DTO。

### POST `/styles/recommend`

请求：

```json
{
  "query": "modern UI dashboard for SaaS",
  "limit": 5
}
```

约束：

- `query` 必填，1-500 字符。
- `limit` 可选，1-50；默认 10。

成功响应：

```json
{
  "query": "modern UI dashboard for SaaS",
  "recommendations": [
    {
      "id": "style-id",
      "name": "Minimal UI Dashboard",
      "prompt": "modern minimal UI dashboard...",
      "category": "ui",
      "description": "Clean modern dashboard style.",
      "tags": ["ui", "dashboard", "saas"],
      "score": 0.75
    }
  ],
  "total": 1
}
```

实现说明：

- 当前 HTTP 运行时使用 handler 内的关键词重叠打分。
- `internal/domain/style/recommendation` 中已有 TF-IDF + Cosine 实现，后续应接入运行时。

## Workflows

所有 workflow 接口均需要 Bearer Token。

当前 Workflow API 管理元数据和状态，不执行完整视频工作流。

### POST `/workflows`

请求：

```json
{
  "style_id": "style-id",
  "config": {
    "image_count": 4,
    "width": 1024,
    "height": 1024,
    "steps": 30,
    "seed": 12345
  }
}
```

约束：

- `style_id` 必填。
- `config.image_count` 必填，1-10。
- `config.width` 和 `config.height` 必填且大于 0。
- `config.steps` 可选，1-100。

成功响应：`201 Created`

```json
{
  "id": "workflow-id",
  "user_id": "user-id",
  "style_id": "style-id",
  "state": "DRAFT",
  "config": {
    "image_count": 4,
    "width": 1024,
    "height": 1024,
    "steps": 30,
    "seed": 12345
  },
  "created_at": "2026-07-08T00:00:00Z",
  "updated_at": "2026-07-08T00:00:00Z"
}
```

### GET `/workflows`

查询参数：

| 参数     | 类型    | 说明         |
| -------- | ------- | ------------ |
| `state`  | string  | 可选状态过滤 |
| `limit`  | integer | 默认 20      |
| `offset` | integer | 默认 0       |

成功响应：

```json
{
  "workflows": [],
  "total": 0,
  "limit": 20,
  "offset": 0
}
```

只返回当前用户自己的 workflow。

### GET `/workflows/:id`

需要 workflow 属于当前用户，否则返回 `403`。

### PATCH `/workflows/:id/status`

请求：

```json
{
  "state": "PENDING"
}
```

支持状态：

- `DRAFT`
- `PENDING`
- `RUNNING`
- `PAUSED`
- `COMPLETED`
- `FAILED`
- `CANCELLED`

当前领域状态机允许：

```text
DRAFT -> PENDING | CANCELLED
PENDING -> RUNNING | CANCELLED
RUNNING -> PAUSED | COMPLETED | FAILED
PAUSED -> RUNNING | CANCELLED
COMPLETED -> terminal
FAILED -> terminal
CANCELLED -> terminal
```

非法状态或非法流转返回 `400`。

## Images

所有 image 接口均需要 Bearer Token。

### POST `/images/generate`

请求：

```json
{
  "prompts": ["modern dashboard hero", "minimal product card"],
  "width": 512,
  "height": 512,
  "quality": "standard",
  "style": "optional style prompt"
}
```

约束：

- `prompts` 必填且不能为空。
- `width` 和 `height` 必填，范围 256-2048。
- `quality` 必须为 `standard` 或 `hd`。
- `style` 可选。

成功响应：

```json
{
  "results": [
    {
      "id": "uuid.png",
      "url": "http://localhost:9000/images/...",
      "prompt": "modern dashboard hero",
      "created_at": "2026-07-08T00:00:00Z"
    }
  ],
  "total": 2,
  "success": 1,
  "failed": 1
}
```

说明：

- 后端对每条 prompt 调用图像 Provider。
- Provider 返回 `image_url` 后，后端下载图片并上传 MinIO `images` bucket。
- `url` 是 24 小时预签名 URL。
- 如果部分 prompt 失败但至少有成功结果，接口仍返回 `200` 并通过 `failed` 计数体现。

### GET `/images/:id`

检查图片是否存在，并返回新的预签名 URL：

```json
{
  "id": "uuid.png",
  "url": "http://localhost:9000/images/..."
}
```

### GET `/images/:id/progress`

当前为同步完成占位响应：

```json
{
  "id": "uuid.png",
  "status": "completed",
  "progress": 100
}
```

## Next.js 代理兼容

`frontend` 提供以下代理：

- `POST /api/users/register`
- `POST /api/users/login`
- `POST /api/styles/recommend`
- `POST /api/images/generate`

代理行为：

- 使用 `BACKEND_URL` 指向 Go API。
- 转发浏览器请求里的 `Authorization`。
- `images/generate` 会把 Go API 的 `results` 同步暴露为 `images`，兼容前端页面。
- `styles/recommend` 会兼容旧字段 `prompt`/`top_k`，转换为 `query`/`limit`。
