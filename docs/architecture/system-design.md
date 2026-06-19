# Labhaus 系统设计

## 系统架构

### 总体架构

```
┌─────────────────────────────────────────────────────────────┐
│                         前端层                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ 可视化编辑器 │  │  样式库UI    │  │  任务监控    │      │
│  │ (React Flow) │  │  (Gallery)   │  │ (Dashboard)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ↓ REST API
┌─────────────────────────────────────────────────────────────┐
│                      API Gateway 层                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Hono + 中间件链                                      │   │
│  │  (日志、认证、CORS、限流、错误处理)                   │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                       业务服务层                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  样式库服务  │  │  工作流引擎  │  │  认证服务    │      │
│  │  (Styles)    │  │  (Workflow)  │  │  (Auth)      │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  推荐算法    │  │  图像生成    │  │  视频合成    │      │
│  │  (TF-IDF)    │  │  (Providers) │  │  (FFmpeg)    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                       数据存储层                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  PostgreSQL  │  │    Redis     │  │    MinIO     │      │
│  │  (主数据库)  │  │   (缓存)     │  │ (对象存储)   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## 核心模块

### 1. 样式库服务

**功能**:
- 样式查询、筛选、搜索
- 基于 TF-IDF 的智能推荐
- 全文搜索（PostgreSQL FTS）

**技术**:
- PostgreSQL + 全文搜索索引
- natural (NLP 库)
- 余弦相似度算法

### 2. 工作流引擎

**功能**:
- 状态机管理（7 种状态）
- 节点执行器
- DAG 验证（循环检测）

**设计**:
```typescript
WorkflowDefinition
  ├── nodes[]          // 节点定义
  ├── edges[]          // 边连接
  └── version          // 版本号

WorkflowExecution
  ├── status           // 执行状态
  ├── current_node     // 当前节点
  ├── context          // 上下文数据
  └── error            // 错误信息
```

**状态转换**:
```
DRAFT → PENDING → RUNNING → COMPLETED
                     ↓
                  PAUSED → RUNNING
                     ↓
                  FAILED → PENDING (retry)
```

### 3. 认证系统

**功能**:
- JWT 认证
- Refresh token 机制
- bcrypt 密码加密

**流程**:
```
1. 注册/登录 → 获取 access_token + refresh_token
2. 请求 API → Authorization: Bearer <token>
3. Token 过期 → 用 refresh_token 获取新 token
```

### 4. 图像生成服务（规划中）

**架构**:
```typescript
interface ImageProvider {
  generate(prompt, options): Promise<ImageResult>
  batchGenerate(prompts[], options): Promise<ImageResult[]>
}

// Provider 实现
- MockProvider (测试)
- OpenAIProvider (DALL-E)
- StableDiffusionProvider (Stable Diffusion)
```

## 数据模型

### 样式库 (styles)

```sql
CREATE TABLE styles (
  id UUID PRIMARY KEY,
  case_id INTEGER UNIQUE,
  title TEXT NOT NULL,
  prompt TEXT NOT NULL,
  category TEXT NOT NULL,
  styles TEXT[],
  scenes TEXT[],
  image_url TEXT,
  featured BOOLEAN,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- 索引
CREATE INDEX idx_styles_category ON styles(category);
CREATE INDEX idx_styles_featured ON styles(featured);
CREATE INDEX idx_styles_title_search ON styles USING gin(to_tsvector('english', title));
```

### 用户 (users)

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT,
  name TEXT,
  google_id TEXT UNIQUE,
  email_verified BOOLEAN,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

### 刷新令牌 (refresh_tokens)

```sql
CREATE TABLE refresh_tokens (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  token TEXT UNIQUE NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP
);
```

## 技术栈

### 当前实现 (Phase 1)

- **语言**: TypeScript
- **后端框架**: Hono (轻量级高性能)
- **数据库**: PostgreSQL 14+
- **缓存**: Redis 7+ (规划中)
- **对象存储**: MinIO / S3
- **认证**: JWT + bcrypt
- **类型验证**: Zod
- **测试**: Vitest
- **Monorepo**: Turborepo + pnpm

### 规划中 (Phase 2-3)

- **前端框架**: React 18 + Next.js 14
- **UI 库**: TailwindCSS + shadcn/ui
- **可视化**: React Flow (工作流编辑器)
- **状态管理**: Zustand
- **图像生成**: OpenAI DALL-E / Stable Diffusion
- **视频合成**: FFmpeg
- **TTS**: Edge-TTS

## API 设计

### RESTful 规范

```
GET    /api/resource         # 列表
GET    /api/resource/:id     # 详情
POST   /api/resource         # 创建
PUT    /api/resource/:id     # 更新
DELETE /api/resource/:id     # 删除
```

### 统一响应格式

**成功响应**:
```json
{
  "data": { ... },
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

**错误响应**:
```json
{
  "error": "ERROR_CODE",
  "message": "Human-readable message",
  "details": { ... }
}
```

## 安全设计

### 1. 认证与授权

- JWT access token (1小时过期)
- JWT refresh token (7天过期，数据库存储)
- bcrypt 密码加密 (10 rounds)

### 2. API 安全

- CORS 配置
- 速率限制 (100 req/min per IP)
- 安全头 (X-Frame-Options, CSP 等)
- 输入验证 (Zod schema)

### 3. 数据安全

- 密码哈希存储
- 敏感数据不记录日志
- 生产环境错误信息脱敏

## 性能优化

### 1. 数据库

- 索引优化 (category, email, FTS)
- 查询分页
- 连接池管理

### 2. 缓存策略（规划）

- Redis 缓存热点数据
- 样式库查询结果缓存
- CDN 缓存静态资源

### 3. 并发控制

- 批量任务并发限制 (10 并发)
- 任务队列管理
- 失败重试机制

## 可扩展性

### 水平扩展

```
              Load Balancer
                    ↓
    ┌───────────────┼───────────────┐
    ↓               ↓               ↓
  API 1           API 2           API 3
    ↓               ↓               ↓
        PostgreSQL (Primary-Replica)
              Redis Cluster
```

### 模块化设计

- Provider 接口抽象
- 插件化节点系统
- 微服务架构预留

## 监控与日志

### 日志级别

- INFO: 正常请求日志
- WARN: 业务警告
- ERROR: 错误和异常

### 监控指标（规划）

- API 响应时间
- 错误率
- 数据库连接池
- 任务成功率

## 部署架构

### Docker Compose (开发/小规模)

```yaml
services:
  api:
    image: labhaus-api
    ports: ["3001:3001"]
  postgres:
    image: postgres:14
  redis:
    image: redis:7
  minio:
    image: minio/minio
```

### Kubernetes (生产/大规模)

- API Deployment (多副本)
- PostgreSQL StatefulSet
- Redis Cluster
- Ingress (负载均衡 + SSL)

## 下一步规划

1. **Phase 2**: 图像生成服务和批量任务管理
2. **Phase 3**: 可视化工作流编辑器
3. **Phase 4**: 模板市场和社区功能
