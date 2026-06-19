# Labhaus Backend

Go DDD (Domain-Driven Design) 后端服务，实现依赖注入和应用启动流程。

## 技术栈

- **Go 1.25+**
- **Gin** - HTTP 框架
- **Viper** - 配置管理
- **Zerolog** - 日志系统
- **GORM** - ORM
- **PostgreSQL** - 数据库
- **Redis** - 缓存
- **Docker** - 容器化

## 项目结构

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # 应用入口，依赖注入
├── internal/
│   ├── application/             # 应用层（CQRS）
│   │   ├── command/            # 命令处理器（写操作）
│   │   ├── query/              # 查询处理器（读操作）
│   │   └── dto/                # 数据传输对象
│   ├── domain/                 # 领域层
│   │   ├── style/              # Style 聚合根
│   │   ├── user/               # User 聚合根
│   │   └── workflow/           # Workflow 聚合根
│   └── infrastructure/         # 基础设施层
│       ├── config/             # 配置管理
│       ├── logger/             # 日志系统
│       ├── http/               # HTTP 路由和处理器
│       └── persistence/        # 数据持久化
└── tests/
    ├── integration/            # 集成测试
    └── unit/                   # 单元测试
```

## 快速开始

### 1. 使用 Docker Compose（推荐）

```bash
# 启动所有服务（PostgreSQL, Redis, MinIO, API）
docker compose up -d --build

# 查看日志
docker compose logs -f api

# 停止服务
docker compose down
```

### 2. 本地开发

**前置条件：**
- Go 1.25+
- PostgreSQL 16+
- Redis 7+

**配置环境变量：**

```bash
cp backend/.env.example backend/.env
# 编辑 .env 文件，配置数据库连接等
```

**运行：**

```bash
cd backend
go mod download
go run cmd/api/main.go
```

## 配置

通过环境变量配置（前缀 `LABHAUS_`）：

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `LABHAUS_SERVER_PORT` | 8080 | 服务端口 |
| `LABHAUS_SERVER_ENVIRONMENT` | development | 运行环境 |
| `LABHAUS_DATABASE_HOST` | localhost | 数据库主机 |
| `LABHAUS_DATABASE_PORT` | 5432 | 数据库端口 |
| `LABHAUS_DATABASE_USER` | postgres | 数据库用户 |
| `LABHAUS_DATABASE_PASSWORD` | postgres | 数据库密码 |
| `LABHAUS_DATABASE_DBNAME` | labhaus | 数据库名称 |
| `LABHAUS_LOG_LEVEL` | info | 日志级别（debug/info/warn/error） |
| `LABHAUS_LOG_FORMAT` | json | 日志格式（json/console） |

## API 端点

### 健康检查

```bash
GET /api/health
```

**响应：**
```json
{
  "status": "healthy",
  "version": "0.1.0"
}
```

### Styles

**列出所有 styles：**
```bash
GET /api/styles?category=video&limit=20&offset=0
```

**获取单个 style：**
```bash
GET /api/styles/:id
```

**创建 style：**
```bash
POST /api/styles
Content-Type: application/json

{
  "name": "Cinematic",
  "description": "Hollywood-style cinematic look",
  "prompt": "cinematic lighting, dramatic shadows, film grain",
  "category": "video",
  "tags": ["cinematic", "dramatic", "film"]
}
```

## 测试

```bash
# 单元测试
cd backend
go test ./internal/...

# 集成测试
go test ./tests/integration/...

# 测试 API
curl http://localhost:8080/api/health
curl http://localhost:8080/api/styles
```

## 依赖注入流程

应用启动时的依赖注入顺序（`cmd/api/main.go`）：

```go
1. 加载配置（Viper）
2. 初始化日志（Zerolog）
3. 连接数据库（GORM + PostgreSQL）
4. 运行数据库迁移
5. 初始化仓储层（实现 Repository 接口）
6. 初始化应用服务（Query/Command Handlers）
7. 初始化 HTTP Handlers
8. 配置路由（Gin）
9. 启动 HTTP 服务器
10. 监听信号，优雅关闭
```

## Graceful Shutdown

应用支持优雅关闭：

- 接收 `SIGINT` 或 `SIGTERM` 信号
- 停止接受新请求
- 等待现有请求完成（默认超时 30 秒）
- 关闭数据库连接
- 退出

```bash
# 发送关闭信号
docker compose stop api
# 或
kill -TERM <pid>
```

## 开发

### 添加新的聚合根

1. 在 `internal/domain/` 创建新目录
2. 定义 Entity 和 Repository 接口
3. 在 `internal/infrastructure/persistence/` 实现 Repository
4. 在 `internal/application/` 添加 Command/Query Handlers
5. 在 `internal/infrastructure/http/handlers/` 添加 HTTP Handler
6. 在 `router.go` 注册路由
7. 在 `main.go` 注入依赖

### 日志

```go
log.Info("message", "key1", value1, "key2", value2)
log.Error("error occurred", err, "context", "value")
log.Debug("debug info", "data", data)
```

## 故障排查

**数据库连接失败：**
```bash
# 检查 PostgreSQL 是否运行
docker compose ps postgres

# 查看数据库日志
docker compose logs postgres
```

**端口被占用：**
```bash
# 检查端口占用
lsof -i :8080

# 修改端口
export LABHAUS_SERVER_PORT=8081
```

## License

MIT
