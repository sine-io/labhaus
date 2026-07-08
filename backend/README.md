# Labhaus Backend

Go 后端服务，负责认证、样式、工作流元数据、批量生图、Redis 队列和 MinIO 存储。

## 技术栈

- Go 1.25+
- Gin
- Viper
- Zerolog
- GORM + PostgreSQL
- Redis
- MinIO / S3
- JWT + bcrypt
- Docker

## 项目结构

```text
backend/
├── cmd/
│   ├── api/                         # Go API 入口
│   └── mock-image-provider/         # 本地 demo 图像 Provider
├── internal/
│   ├── application/                 # CQRS handlers、DTO、批量生图服务
│   ├── domain/                      # User、Style、Workflow、Image Provider 领域模型
│   └── infrastructure/
│       ├── auth/                    # JWT
│       ├── config/                  # Viper 配置
│       ├── http/                    # Gin router、middleware、handlers
│       ├── image/                   # Provider 适配
│       ├── persistence/             # GORM repositories
│       ├── queue/                   # Redis queue
│       └── storage/                 # MinIO storage
├── seeds/styles.sql                 # 本地 demo 样式数据
├── tests/                           # 集成测试
├── Dockerfile
├── Dockerfile.mock-image-provider
└── go.mod
```

## 本地启动

### Docker Compose（推荐演示）

从仓库根目录：

```bash
docker compose up -d --build
docker compose exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql
docker compose restart api
docker compose logs -f api
```

Compose 会启动：

- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`
- MinIO: `localhost:9000`，Console: `localhost:9001`
- Mock Image Provider: `localhost:8089`
- API: `localhost:8080`

### 裸跑 Go API

先启动依赖服务：

```bash
docker compose up -d postgres redis minio mock-image-provider
docker compose exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql
```

配置环境变量：

```bash
cp backend/.env.example backend/.env
```

然后运行：

```bash
cd backend
go mod download
go run cmd/api/main.go
```

健康检查：

```bash
curl http://localhost:8080/api/health
```

## 配置

配置通过 `LABHAUS_` 前缀环境变量读取。

| 环境变量                          | 默认值                                 | 说明                                               |
| --------------------------------- | -------------------------------------- | -------------------------------------------------- |
| `LABHAUS_SERVER_PORT`             | `8080`                                 | API 端口                                           |
| `LABHAUS_SERVER_ENVIRONMENT`      | `development`                          | 运行环境                                           |
| `LABHAUS_SERVER_SHUTDOWN_TIMEOUT` | `30`                                   | 优雅关闭超时，单位秒                               |
| `LABHAUS_DATABASE_HOST`           | `localhost`                            | PostgreSQL 主机                                    |
| `LABHAUS_DATABASE_PORT`           | `5432`                                 | PostgreSQL 端口                                    |
| `LABHAUS_DATABASE_USER`           | `postgres`                             | PostgreSQL 用户                                    |
| `LABHAUS_DATABASE_PASSWORD`       | `postgres`                             | PostgreSQL 密码                                    |
| `LABHAUS_DATABASE_DBNAME`         | `labhaus`                              | PostgreSQL 数据库                                  |
| `LABHAUS_DATABASE_SSLMODE`        | `disable`                              | PostgreSQL SSL 模式                                |
| `LABHAUS_REDIS_HOST`              | `localhost`                            | Redis 主机                                         |
| `LABHAUS_REDIS_PORT`              | `6379`                                 | Redis 端口                                         |
| `LABHAUS_REDIS_PASSWORD`          | 空                                     | Redis 密码                                         |
| `LABHAUS_REDIS_DB`                | `0`                                    | Redis DB                                           |
| `LABHAUS_MINIO_ENDPOINT`          | `localhost:9000`                       | MinIO/S3 endpoint                                  |
| `LABHAUS_MINIO_ACCESS_KEY`        | `minioadmin`                           | MinIO access key                                   |
| `LABHAUS_MINIO_SECRET_KEY`        | `minioadmin`                           | MinIO secret key                                   |
| `LABHAUS_MINIO_USE_SSL`           | `false`                                | 是否使用 HTTPS                                     |
| `LABHAUS_IMAGE_PROVIDER_BASE_URL` | 必填                                   | 图像 Provider base URL，需实现 `POST /v1/generate` |
| `LABHAUS_IMAGE_PROVIDER_API_KEY`  | 必填                                   | 图像 Provider API key                              |
| `LABHAUS_LOG_LEVEL`               | `info`                                 | 日志级别                                           |
| `LABHAUS_LOG_FORMAT`              | `json`                                 | `json` 或 `console`                                |
| `LABHAUS_JWT_SECRET_KEY`          | `your-secret-key-change-in-production` | JWT 签名密钥                                       |
| `LABHAUS_JWT_TOKEN_DURATION`      | `24`                                   | JWT 有效期，单位小时                               |

`LABHAUS_IMAGE_PROVIDER_BASE_URL` 和 `LABHAUS_IMAGE_PROVIDER_API_KEY` 缺失时，API 会在启动阶段失败。Docker Compose 会为本地 demo 自动注入 mock provider 默认值。

## API 端点

完整契约见 `docs/architecture/api-design.md`。

### Public

```http
GET  /api/health
POST /api/users/register
POST /api/users/login
```

登录成功响应：

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

### Authenticated

```http
GET   /api/users/me
PATCH /api/users/me

GET  /api/styles
GET  /api/styles/:id
POST /api/styles
POST /api/styles/recommend

POST  /api/workflows
GET   /api/workflows
GET   /api/workflows/:id
PATCH /api/workflows/:id/status

POST /api/images/generate
GET  /api/images/:id
GET  /api/images/:id/progress
```

受保护接口需要：

```http
Authorization: Bearer <token>
Content-Type: application/json
```

## 样式 seed

`backend/seeds/styles.sql` 是本地 demo seed，目前包含 12 条代表性样式，并使用固定 ID 和 `ON CONFLICT (id) DO UPDATE`，可重复执行。

当前产品目标仍是导入和维护 500+ GPT-Image-2 样式库；当前代码仓库中的 seed 不等同于完整生产样式库。

导入后需要重启 API，因为运行时推荐器在启动时加载样式快照：

```bash
docker compose exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql
docker compose restart api
```

## 图像 Provider 合约

后端当前通过 `gptimage2` 适配器调用兼容 HTTP Provider：

```http
POST /v1/generate
Authorization: Bearer <api-key>
Content-Type: application/json
```

请求：

```json
{
  "prompt": "modern dashboard",
  "width": 1024,
  "height": 1024,
  "quality": "standard",
  "style": "optional style prompt"
}
```

响应：

```json
{
  "image_url": "http://provider/images/abc.png",
  "created_at": "2026-07-08T00:00:00Z"
}
```

后端会下载 `image_url` 内容，上传到 MinIO `images` bucket，并返回 MinIO 预签名 URL。

## 测试

```bash
cd backend

# 全部 Go 测试
go test ./...

# 单包示例
go test ./internal/infrastructure/http/handlers -v
go test ./cmd/mock-image-provider -v
```

部分存储或集成测试需要正在运行的 PostgreSQL / MinIO 服务。

## 依赖注入流程

`cmd/api/main.go` 启动顺序：

1. 加载配置。
2. 初始化日志。
3. 连接 PostgreSQL 并执行 GORM AutoMigrate。
4. 连接 Redis。
5. 初始化 MinIO 并确保 `workflows`、`images`、`videos`、`temp` buckets 存在。
6. 初始化图像 Provider、批量生图服务和 MinIO image storage。
7. 初始化 Redis queue 并启动 workflow task worker。
8. 初始化 bcrypt hasher 和 JWT service。
9. 初始化 repositories、query/command handlers、HTTP handlers。
10. 注册 Gin routes 并启动 HTTP server。
11. 接收 `SIGINT` / `SIGTERM` 后优雅关闭。

## 当前限制

- 样式推荐 HTTP 运行时当前使用 handler 内的关键词重叠打分；`domain/style/recommendation` 中已有 TF-IDF + Cosine 实现，后续应接入运行时。
- Workflow API 当前管理工作流元数据和状态，不执行完整视频工作流。
- Redis queue 的 workflow worker 当前只解析并打印任务，实际工作流执行逻辑待实现。
- API 尚未实现 refresh token、OAuth、RBAC、速率限制、OpenAPI 自动生成。

## License

GNU Affero General Public License v3.0。详见仓库根目录 `LICENSE`。
