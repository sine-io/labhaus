# Labhaus API 设计文档

## API 基础

**Base URL**: `http://localhost:3001/api`  
**协议**: REST + JSON  
**认证**: JWT Bearer Token

## 通用规范

### 请求头

```
Content-Type: application/json
Authorization: Bearer <access_token>  # 需要认证的端点
```

### 响应格式

**成功响应** (2xx):
```json
{
  "data": { ... },
  "pagination": {    // 列表接口才有
    "page": 1,
    "limit": 20,
    "total": 100,
    "totalPages": 5
  }
}
```

**错误响应** (4xx/5xx):
```json
{
  "error": "ERROR_CODE",
  "message": "Human-readable error message",
  "details": { ... }  // 可选，仅验证错误
}
```

### 错误码

| 状态码 | 错误码 | 说明 |
|--------|--------|------|
| 400 | BAD_REQUEST | 请求参数错误 |
| 400 | VALIDATION_ERROR | 数据验证失败 |
| 401 | UNAUTHORIZED | 未认证或 token 无效 |
| 403 | FORBIDDEN | 无权访问 |
| 404 | NOT_FOUND | 资源不存在 |
| 409 | CONFLICT | 资源冲突（如重复注册）|
| 429 | RATE_LIMIT_EXCEEDED | 超过速率限制 |
| 500 | INTERNAL_ERROR | 服务器内部错误 |

## 端点列表

### 1. 健康检查

#### GET /api/health

**描述**: 服务健康检查

**响应**:
```json
{
  "status": "ok",
  "timestamp": "2026-06-19T12:00:00Z"
}
```

### 2. API 信息

#### GET /api

**描述**: 获取 API 版本和端点列表

**响应**:
```json
{
  "name": "Labhaus API",
  "version": "0.1.0",
  "endpoints": {
    "health": "/api/health",
    "styles": "/api/styles",
    "auth": "/api/auth"
  }
}
```

---

## 样式库 API

### 1. 获取样式列表

#### GET /api/styles

**描述**: 查询样式库，支持筛选、搜索和分页

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| category | string | 否 | 按分类筛选 |
| style | string | 否 | 按风格筛选 |
| scene | string | 否 | 按场景筛选 |
| featured | boolean | 否 | 是否精选 |
| search | string | 否 | 全文搜索关键词 |
| page | integer | 否 | 页码（默认 1）|
| limit | integer | 否 | 每页数量（默认 20，最大 100）|

**示例**:
```bash
GET /api/styles?category=UI%20%26%20Interfaces&limit=10
GET /api/styles?search=portrait&page=2
```

**响应**:
```json
{
  "styles": [
    {
      "id": "uuid",
      "case_id": 505,
      "title": "夜间手机光沙发肖像",
      "prompt": "A young adult woman...",
      "prompt_preview": "A young adult woman...",
      "category": "Photography & Realism",
      "styles": ["Realistic"],
      "scenes": ["Tech", "Commerce"],
      "image_url": "/images/case505.jpg",
      "featured": false,
      "created_at": "2026-06-19T...",
      "updated_at": "2026-06-19T..."
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 503,
    "totalPages": 51
  }
}
```

### 2. 获取样式详情

#### GET /api/styles/:id

**描述**: 根据 ID 获取单个样式详情

**路径参数**:
- `id` (uuid): 样式 ID

**示例**:
```bash
GET /api/styles/550e8400-e29b-41d4-a716-446655440000
```

**响应**:
```json
{
  "id": "uuid",
  "case_id": 505,
  "title": "...",
  "prompt": "...",
  ...
}
```

### 3. 样式推荐

#### POST /api/styles/recommend

**描述**: 基于查询文本推荐相关样式（TF-IDF + 余弦相似度）

**请求体**:
```json
{
  "query": "modern minimalist UI design",
  "limit": 10
}
```

**响应**:
```json
{
  "query": "modern minimalist UI design",
  "recommendations": [
    {
      "style": { ...完整样式对象... },
      "score": 0.856
    }
  ],
  "total": 10
}
```

### 4. 相似样式

#### GET /api/styles/:id/similar

**描述**: 查找与指定样式相似的其他样式

**查询参数**:
- `limit` (integer): 返回数量（默认 10，最大 50）

**示例**:
```bash
GET /api/styles/550e8400-e29b-41d4-a716-446655440000/similar?limit=5
```

**响应**:
```json
{
  "style_id": "uuid",
  "recommendations": [
    {
      "style": { ... },
      "score": 0.742
    }
  ],
  "total": 5
}
```

---

## 认证 API

### 1. 用户注册

#### POST /api/auth/register

**描述**: 注册新用户账号

**请求体**:
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "name": "John Doe"  // 可选
}
```

**响应** (201):
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "email_verified": false,
    "created_at": "2026-06-19T..."
  },
  "tokens": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

### 2. 用户登录

#### POST /api/auth/login

**描述**: 使用邮箱和密码登录

**请求体**:
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**响应** (200):
```json
{
  "user": { ... },
  "tokens": { ... }
}
```

### 3. 刷新 Token

#### POST /api/auth/refresh

**描述**: 使用 refresh token 获取新的 access token

**请求体**:
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**响应** (200):
```json
{
  "tokens": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

### 4. 获取当前用户

#### GET /api/auth/me

**描述**: 获取当前认证用户信息

**Headers**:
```
Authorization: Bearer <access_token>
```

**响应** (200):
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "avatar_url": null,
    "email_verified": false,
    "created_at": "2026-06-19T..."
  }
}
```

---

## 速率限制

| 环境 | 限制 |
|------|------|
| 开发环境 | 无限制 |
| 生产环境 | 100 请求/分钟 per IP |

**速率限制响应头**:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1624291200
```

---

## 完整示例

### 完整认证流程

```bash
# 1. 注册账号
curl -X POST http://localhost:3001/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@labhaus.io","password":"Demo123!","name":"Demo"}'

# 2. 保存返回的 access_token

# 3. 使用 token 查询样式
curl http://localhost:3001/api/styles?limit=5 \
  -H "Authorization: Bearer eyJhbGc..."

# 4. 样式推荐
curl -X POST http://localhost:3001/api/styles/recommend \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGc..." \
  -d '{"query":"modern UI design","limit":10}'

# 5. Token 过期后刷新
curl -X POST http://localhost:3001/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"eyJhbGc..."}'
```

## 更多文档

- [认证详细文档](../../apps/api/docs/AUTHENTICATION.md)
- [API 路由设计](../../apps/api/docs/API_DESIGN.md)
- [部署指南](../DEPLOYMENT.md)
