# Labhaus 技术栈文档

## 总览

Labhaus 采用 **Go 后端 + TypeScript 前端**架构，基于现代化的技术栈构建。

## 后端技术栈 (Go)

### 核心框架
- **语言**: Go 1.21+
- **Web 框架**: Gin
  - 为什么选择 Gin？
    - 性能优秀 (40x faster than Martini)
    - 中间件生态成熟
    - 社区活跃 (Star 75k+)
    - 文档完善

### 数据层
- **主数据库**: PostgreSQL 14+
  - 全文搜索 (FTS)
  - 数组类型支持
  - 事务保证
- **ORM**: GORM (常规 CRUD)
- **SQL Builder**: sqlc (复杂查询，类型安全)
- **缓存**: Redis 7+
- **对象存储**: MinIO / Amazon S3

### 认证与安全
- **认证**: golang-jwt/jwt/v5
  - Access token (1小时过期)
  - Refresh token (7天过期)
- **密码加密**: golang.org/x/crypto/bcrypt
- **数据验证**: go-playground/validator/v10

### 并发与任务
- **任务队列**: Asynq (Redis-based)
  - 持久化任务
  - 失败重试
  - 定时任务
  - Dashboard 监控
- **并发控制**: Goroutine + Channel

### 配置与日志
- **配置管理**: Viper
  - 环境变量
  - 配置文件
  - 热重载
- **日志**: zerolog
  - 零分配，性能最佳
  - 结构化日志
  - 友好的 API

### HTTP 客户端
- **HTTP**: resty (类似 axios)

### 测试
- **测试框架**: testing (标准库)
- **断言**: testify

## 前端技术栈（规划中）

### 核心框架
- **框架**: React 18
- **路由**: Next.js 14 (App Router)
  - 服务端渲染 (SSR)
  - 静态生成 (SSG)
  - API Routes

### UI 库
- **CSS 框架**: TailwindCSS 3+
- **组件库**: shadcn/ui

### 可视化
- **工作流编辑器**: React Flow

### 状态管理
- **全局状态**: Zustand

### 实时通信
- **服务器推送**: Server-Sent Events (SSE)

## AI / 媒体处理（规划中）

### 图像生成
- **OpenAI DALL-E 3** (API 集成)
- **Stable Diffusion** (可选)

### 文本生成
- **OpenAI GPT-4** (剧本生成)

### 语音合成
- **Edge-TTS** (免费、高质量)

### 视频处理
- **FFmpeg**
  - 图片合成
  - 音频混合
  - 字幕渲染

## 基础设施

### 容器化
- **Docker**
- **Docker Compose** (开发环境)
- **Kubernetes** (生产环境规划)

### CI/CD
- **GitHub Actions**
  - 自动测试
  - 代码检查
  - 自动部署 (规划中)

### 监控与日志（规划中）
- **APM**: Prometheus + Grafana
- **错误追踪**: Sentry

## 完整技术栈图

```
┌─────────────────────────────────────────────┐
│             前端层（规划中）                 │
│  React 18 + Next.js 14 + TailwindCSS        │
│  React Flow + Zustand + shadcn/ui           │
└─────────────────────────────────────────────┘
                     ↓ REST API
┌─────────────────────────────────────────────┐
│              API Gateway (Go)                │
│  Gin + 中间件 (日志/认证/CORS/限流)         │
└─────────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────┐
│            业务服务层 (Go)                   │
│  GORM + sqlc + validator + JWT              │
└─────────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────┐
│           任务队列 (Asynq)                   │
│  批量任务 + 失败重试 + 定时任务              │
└─────────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────┐
│         AI / 媒体处理（规划中）              │
│  OpenAI (GPT-4/DALL-E) + Edge-TTS + FFmpeg  │
└─────────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────┐
│              数据存储层                      │
│  PostgreSQL + Redis + MinIO/S3              │
└─────────────────────────────────────────────┘
```

## 项目结构

```
labhaus/
├── backend/                      # Go 后端
│   ├── cmd/
│   │   └── api/
│   │       └── main.go          # 入口
│   ├── internal/
│   │   ├── api/                 # HTTP 层
│   │   │   ├── middleware/      # 中间件
│   │   │   ├── routes/          # 路由
│   │   │   └── handlers/        # 处理器
│   │   ├── service/             # 业务逻辑
│   │   │   ├── auth/
│   │   │   ├── style/
│   │   │   └── recommendation/
│   │   ├── repository/          # 数据访问
│   │   │   └── postgres/
│   │   ├── model/               # 数据模型
│   │   ├── config/              # 配置
│   │   └── pkg/                 # 工具包
│   ├── migrations/              # 数据库迁移
│   ├── tests/                   # 测试
│   └── go.mod
├── apps/
│   └── web/                     # Next.js 前端 (规划中)
├── packages/                    # TypeScript 共享包 (遗留)
└── docs/                        # 文档
```

## Go 依赖 (go.mod)

```go
module github.com/sine-io/labhaus

go 1.21

require (
    github.com/gin-gonic/gin v1.10.0           // Web 框架
    gorm.io/gorm v1.25.5                       // ORM
    gorm.io/driver/postgres v1.5.4             // PostgreSQL 驱动
    github.com/google/uuid v1.5.0              // UUID
    github.com/golang-jwt/jwt/v5 v5.2.0        // JWT
    golang.org/x/crypto v0.17.0                // bcrypt
    github.com/go-playground/validator/v10     // 验证
    github.com/spf13/viper v1.18.2             // 配置
    github.com/redis/go-redis/v9 v9.4.0        // Redis
    github.com/hibiken/asynq v0.24.1           // 任务队列
    github.com/minio/minio-go/v7 v7.0.66       // 对象存储
    github.com/rs/zerolog v1.31.0              // 日志
    github.com/go-resty/resty/v2 v2.11.0       // HTTP 客户端
    github.com/stretchr/testify v1.8.4         // 测试
)
```

## 为什么选择 Go？

### 性能优势
✅ **并发性能**: Goroutine 比 Node.js async/await 高效 10x  
✅ **内存占用**: 约 Node.js 的 1/3  
✅ **API 吞吐量**: 2-3x TypeScript  
✅ **启动时间**: < 1s (vs Node.js 2-3s)

### 工程优势
✅ **类型安全**: 编译时错误检测  
✅ **并发模型**: Goroutine + Channel 天然支持  
✅ **部署简单**: 单一二进制文件  
✅ **生态成熟**: 云原生工具首选语言

### 适合场景
✅ 批量任务处理（图像生成、视频合成）  
✅ 高并发 API 服务  
✅ 长连接、实时通信  
✅ 微服务架构

## 技术债务（TypeScript 遗留）

### 已废弃
- ❌ apps/api (TypeScript Hono)
- ❌ packages/workflow (TypeScript)
- ❌ packages/types (TypeScript)

### 保留用途
- 📦 仅作为 Phase 1 参考实现
- 📦 前端开发时可复用类型定义

## 迁移计划

### Phase 1: Go 基础框架 (Week 1)
- ✅ 项目结构
- ✅ Gin + 中间件
- ✅ PostgreSQL + GORM
- ✅ 配置 + 日志

### Phase 2: 核心功能迁移 (Week 2)
- ✅ 认证系统
- ✅ 样式库 API
- ✅ 推荐算法
- ✅ 测试覆盖

### Phase 3: 高级功能 (Week 3)
- ✅ Redis 缓存
- ✅ Asynq 任务队列
- ✅ MinIO 集成
- ✅ E2E 测试

### Phase 4: 部署与文档 (Week 4)
- ✅ Docker 配置
- ✅ 性能测试
- ✅ 文档更新
- ✅ 上线准备

## 更多文档

- [系统设计](./architecture/system-design.md)
- [API 设计](./architecture/api-design.md)
- [部署指南](./DEPLOYMENT.md)
- [Go 迁移指南](./GO_MIGRATION.md) (待创建)
