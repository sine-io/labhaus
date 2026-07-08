# Go 后端架构

**最后更新**：2026-07-08

## 架构原则

Labhaus Go 后端采用轻量级：

- DDD Lite
- CQRS
- Clean Architecture
- 依赖倒置

目标是保持业务边界清晰，同时避免过度抽象。

## 当前目录结构

```text
backend/
├── cmd/
│   ├── api/                         # API 入口和依赖注入
│   └── mock-image-provider/         # 本地 demo Provider
├── internal/
│   ├── domain/                      # 领域层
│   │   ├── image/provider/
│   │   ├── style/
│   │   ├── user/
│   │   └── workflow/
│   ├── application/                 # 应用层
│   │   ├── command/
│   │   ├── dto/
│   │   ├── query/
│   │   └── service/image/
│   └── infrastructure/              # 基础设施层
│       ├── auth/
│       ├── config/
│       ├── http/
│       ├── image/
│       ├── logger/
│       ├── persistence/
│       ├── queue/
│       └── storage/
├── seeds/
├── tests/
├── Dockerfile
└── go.mod
```

## 分层说明

### Domain

领域层只表达业务模型、规则和接口，不依赖 Gin、GORM、Redis、MinIO。

当前领域：

- `user.Entity`
  - email/name/role/password hash
  - role: `user` / `admin`
- `style.Entity`
  - name/description/prompt/category/tags
  - 字段长度校验
- `workflow.Entity`
  - state machine
  - config/result
- `image/provider`
  - `ImageProvider` 接口
  - `ImageOptions`
  - `ImageResult`
  - Provider registry

### Application

应用层编排领域对象和 repository 接口。

当前 CQRS handlers：

- `command.UserCommandHandler`
- `query.UserQueryHandler`
- `command.StyleCommandHandler`
- `query.StyleQueryHandler`
- `command.WorkflowCommandHandler`
- `query.WorkflowQueryHandler`

当前 service：

- `service/image.BatchImageService`
  - 单图生成
  - 批量生成
  - semaphore 并发控制
  - progress channel 支持
  - 部分失败聚合为 `BatchError`

### Infrastructure

基础设施层实现外部细节：

- Gin router、middleware、handlers。
- GORM repositories。
- JWT service。
- bcrypt hasher。
- Redis queue。
- MinIO storage。
- GPT image HTTP Provider。
- Mock image provider。

## 依赖方向

```text
infrastructure -> application -> domain
```

示例：

- `domain/style.Repository` 定义接口。
- `infrastructure/persistence.StyleRepository` 用 GORM 实现接口。
- `application/query.StyleQueryHandler` 依赖 `style.Repository` 接口。
- `cmd/api/main.go` 负责把具体实现注入 handler。

## CQRS 实践

### Command

写操作：

- 注册用户。
- 更新用户信息。
- 创建/更新/删除 style。
- 创建 workflow。
- 更新 workflow state。

示例流程：

```text
HTTP request -> command handler -> domain validation -> repository -> DTO response
```

### Query

读操作：

- 认证用户。
- 获取当前用户。
- 获取 style 列表/详情。
- 获取 workflow 列表/详情。

示例流程：

```text
HTTP request -> query handler -> repository -> DTO response
```

## 状态机

Workflow 状态：

```text
DRAFT
PENDING
RUNNING
PAUSED
COMPLETED
FAILED
CANCELLED
```

允许流转：

```text
DRAFT -> PENDING | CANCELLED
PENDING -> RUNNING | CANCELLED
RUNNING -> PAUSED | COMPLETED | FAILED
PAUSED -> RUNNING | CANCELLED
COMPLETED -> terminal
FAILED -> terminal
CANCELLED -> terminal
```

当前 API 只持久化 workflow 元数据和状态，不执行完整视频工作流。

## 推荐算法边界

当前代码中存在：

- `internal/infrastructure/http/handlers.StaticStyleRecommender`
  - 当前 HTTP 运行时使用。
  - 基于 query tokens 与 style tokens 的重叠比例打分。
- `internal/domain/style/recommendation`
  - 已实现 TF-IDF 和 Cosine Similarity。
  - 后续应接入运行时，替代 handler 内的临时推荐器。

## 启动注入流程

`cmd/api/main.go`：

1. Load config。
2. Create logger。
3. Connect PostgreSQL。
4. AutoMigrate `styles`、`users`、`workflows`。
5. Connect Redis。
6. Initialize MinIO buckets。
7. Create image provider。
8. Create batch image service。
9. Create image storage。
10. Start Redis queue worker。
11. Create hasher and JWT service。
12. Create repositories。
13. Load styles and create recommender snapshot。
14. Create query/command handlers。
15. Create HTTP handlers。
16. Setup router and start server。
17. Graceful shutdown。

## 测试策略

当前仓库包含：

- 领域单元测试。
- 应用层 command/query 测试。
- HTTP handler 测试。
- Provider 和 storage 测试。
- 集成测试（需要数据库/MinIO）。

常用命令：

```bash
cd backend
go test ./...
go test ./internal/infrastructure/http/handlers -v
go test ./cmd/mock-image-provider -v
```

## 后续改进

- 接入 TF-IDF + Cosine 推荐器。
- 把 Redis queue worker 接到真实视频工作流执行器。
- 增加任务执行日志和中间产物模型。
- 增加 OpenAPI 生成。
- 增加统一错误码和错误响应格式。
- 增加限流、审计日志和 RBAC。
