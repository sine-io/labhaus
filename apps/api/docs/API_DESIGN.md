# API 路由设计

## 路由结构

```
/                           → Redirect to /api
/api/                       → API info & endpoint list
/api/health                 → Health check

/api/styles                 → Style library
  GET  /                    → List/search styles
  GET  /:id                 → Get style by ID

/api/workflows (future)     → Workflow management
  POST   /                  → Create workflow
  GET    /:id               → Get workflow
  PUT    /:id               → Update workflow
  DELETE /:id               → Delete workflow

/api/executions (future)    → Workflow execution
  POST   /                  → Start execution
  GET    /:id               → Get execution status
  POST   /:id/pause         → Pause execution
  POST   /:id/resume        → Resume execution
  POST   /:id/cancel        → Cancel execution
```

## 错误响应格式

所有错误遵循统一格式：

```json
{
  "error": "ERROR_CODE",
  "message": "Human-readable message",
  "details": {}  // Optional, only for validation errors
}
```

### 常见错误码

- `BAD_REQUEST` - 400
- `UNAUTHORIZED` - 401
- `FORBIDDEN` - 403
- `NOT_FOUND` - 404
- `CONFLICT` - 409
- `VALIDATION_ERROR` - 400
- `RATE_LIMIT_EXCEEDED` - 429
- `INTERNAL_ERROR` - 500

## 中间件

### 1. Request Logger
记录所有请求：method, URL, status, duration

### 2. Security Headers
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`

### 3. CORS
- Origin: 可配置（默认 localhost:3000）
- Credentials: true
- Methods: GET, POST, PUT, DELETE, PATCH, OPTIONS

### 4. Rate Limiting (生产环境)
- 100 requests / minute per IP
- Headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset

## 环境变量

```bash
API_PORT=3001
API_HOST=0.0.0.0
NODE_ENV=development|production
CORS_ORIGIN=http://localhost:3000,https://labhaus.io
DATABASE_URL=postgresql://...
```
