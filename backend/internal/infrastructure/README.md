# Infrastructure Layer

基础设施层承接框架、数据库、队列、存储、认证和外部 Provider 适配，实现领域层与应用层定义的接口。

## 组件

### HTTP

- Gin router：`internal/infrastructure/http/router.go`
- Auth middleware：JWT Bearer Token 校验，并把 `user_id`、`user_email`、`user_role` 写入 Gin context
- Handlers：
  - `health.go`
  - `user.go`
  - `style.go`
  - `workflow.go`
  - `image.go`

### Persistence

GORM + PostgreSQL 实现：

- `StyleRepository`
  - CRUD
  - category / tags 过滤
  - LIKE 搜索
  - tags JSON 字符串序列化
  - soft delete
- `UserRepository`
  - CRUD
  - 按 email 查询
  - email 唯一性检查
  - soft delete
- `WorkflowRepository`
  - CRUD
  - 按 user 查询
  - state 过滤
  - config/result JSONB
  - state 快速更新
  - soft delete

### Auth

- `auth.JWTService`：生成和校验 HS256 JWT。
- `persistence.BcryptHasher`：bcrypt hash/compare，默认 cost 为 10。

### Queue

- Redis list + hash 实现的轻量队列。
- 支持 pending/processing/completed/dead 状态、retry count 和 dead letter queue。
- 当前 workflow task handler 只解析 payload 并打印日志；实际视频工作流执行逻辑待接入。

### Storage

- `MinIOStorage`：通用 bucket/object 操作。
- `MinIOImageStorage`：当前批量生图接口使用的 image storage，负责上传、下载、删除、存在性检查和预签名 URL。
- API 启动时会确保 `workflows`、`images`、`videos`、`temp` buckets 存在。

### Image Provider

- `image/gptimage2`：兼容 `POST /v1/generate` 的 HTTP Provider 适配器。
- `image/mock`：测试用 Provider。
- `cmd/mock-image-provider`：本地 demo HTTP Provider，可通过 Docker Compose 启动。

## 数据库模型

### StyleModel

```text
id           varchar(36) primary key
name         varchar(100) not null, index
description  varchar(500)
prompt       text not null
category     varchar(50), index
tags         text            # JSON array string
created_at   timestamp not null
updated_at   timestamp not null
deleted_at   timestamp, index
```

### UserModel

```text
id            varchar(36) primary key
email         varchar(255) not null, unique index
password_hash varchar(255) not null
name          varchar(100) not null
role          varchar(20) not null default 'user'
created_at    timestamp not null
updated_at    timestamp not null
deleted_at    timestamp, index
```

### WorkflowModel

```text
id         varchar(36) primary key
user_id    varchar(36) not null, index
style_id   varchar(36) not null, index
state      varchar(20) not null, index
config     jsonb not null
result     jsonb null
created_at timestamp not null
updated_at timestamp not null
deleted_at timestamp, index
```

## 使用示例

```go
db, err := persistence.NewDB(persistence.DBConfig{
    Host:     "localhost",
    Port:     5432,
    User:     "labhaus",
    Password: "labhaus_dev_password",
    DBName:   "labhaus",
    SSLMode:  "disable",
})
if err != nil {
    log.Fatal(err)
}

if err := persistence.AutoMigrate(db); err != nil {
    log.Fatal(err)
}

styleRepo := persistence.NewStyleRepository(db)
userRepo := persistence.NewUserRepository(db)
workflowRepo := persistence.NewWorkflowRepository(db)
hasher := persistence.NewBcryptHasher()
```

## 集成测试

集成测试需要 PostgreSQL：

```bash
docker run -d \
  --name postgres-test \
  -e POSTGRES_USER=labhaus \
  -e POSTGRES_PASSWORD=labhaus_dev_password \
  -e POSTGRES_DB=labhaus \
  -p 5432:5432 \
  postgres:16-alpine

cd backend
go test ./tests/integration/... -v
```

MinIO storage 相关测试需要 MinIO 服务运行。
