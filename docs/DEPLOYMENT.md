# Labhaus 部署指南

当前可部署对象：

- `backend/`：Go API，默认端口 `8080`。
- `apps/web/`：Next.js 前端，默认端口 `3000`。
- 依赖服务：PostgreSQL、Redis、MinIO/S3、图像 Provider。

Docker Compose 当前只编排后端依赖、mock image provider 和 Go API；Web 前端需要单独运行或部署到支持 Next.js 的平台。

## 环境要求

- Docker 和 Docker Compose
- Node.js 20+、pnpm 11（构建/运行前端）
- Go 1.25+（裸机编译/运行 Go API）

## 方式 1：Docker Compose 本地演示

```bash
git clone https://github.com/sine-io/labhaus.git
cd labhaus

docker compose up -d --build
docker compose exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql
docker compose restart api
```

访问：

- API: http://localhost:8080
- MinIO Console: http://localhost:9001
- Mock Image Provider: http://localhost:8089

启动 Web：

```bash
pnpm install
cp apps/web/.env.example apps/web/.env.local
pnpm --filter @labhaus/web dev
```

访问 Web: http://localhost:3000

## 方式 2：裸机运行 Go API

### 1. 准备依赖

```bash
docker compose up -d postgres redis minio mock-image-provider
docker compose exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql
```

### 2. 配置环境变量

最小本地配置：

```bash
LABHAUS_SERVER_PORT=8080
LABHAUS_SERVER_ENVIRONMENT=development
LABHAUS_DATABASE_HOST=localhost
LABHAUS_DATABASE_PORT=5432
LABHAUS_DATABASE_USER=labhaus
LABHAUS_DATABASE_PASSWORD=labhaus_dev_password
LABHAUS_DATABASE_DBNAME=labhaus
LABHAUS_DATABASE_SSLMODE=disable
LABHAUS_REDIS_HOST=localhost
LABHAUS_REDIS_PORT=6379
LABHAUS_MINIO_ENDPOINT=localhost:9000
LABHAUS_MINIO_ACCESS_KEY=minioadmin
LABHAUS_MINIO_SECRET_KEY=minioadmin
LABHAUS_MINIO_USE_SSL=false
LABHAUS_JWT_SECRET_KEY=change-me
LABHAUS_JWT_TOKEN_DURATION=24
LABHAUS_IMAGE_PROVIDER_BASE_URL=http://localhost:8089
LABHAUS_IMAGE_PROVIDER_API_KEY=dev-mock-key
```

`LABHAUS_IMAGE_PROVIDER_BASE_URL` 和 `LABHAUS_IMAGE_PROVIDER_API_KEY` 是必填项；缺失时 API 会启动失败。

### 3. 启动

```bash
cd backend
go mod download
go run cmd/api/main.go
```

健康检查：

```bash
curl http://localhost:8080/api/health
```

## 方式 3：前端生产构建

```bash
pnpm install --frozen-lockfile
cp apps/web/.env.example apps/web/.env.local
pnpm --filter @labhaus/web build
pnpm --filter @labhaus/web start
```

`apps/web/.env.local`：

```bash
BACKEND_URL=http://localhost:8080
```

`BACKEND_URL` 是 Next.js 服务端代理访问的 Go API 地址。

## 生产部署建议

### Go API

- 使用容器平台、systemd 或进程管理器托管。
- 使用 PostgreSQL 16+。
- 使用 Redis 7+。
- 使用生产级 MinIO/S3。
- 使用真实图像 Provider，并按环境隔离 API key。
- 使用强随机 `LABHAUS_JWT_SECRET_KEY`。

编译示例：

```bash
cd backend
go build -o labhaus-api ./cmd/api
LABHAUS_SERVER_ENVIRONMENT=production ./labhaus-api
```

### Next.js 前端

- 使用 `pnpm --filter @labhaus/web build` 构建。
- 运行时设置 `BACKEND_URL` 指向 Go API 的内网或公网地址。
- 如果前端和 API 分域部署，浏览器仍请求 Next.js 自身 `/api/*` 代理，不需要暴露 Go API 给浏览器。

### 图像 Provider

生产 Provider 需兼容：

```http
POST /v1/generate
Authorization: Bearer <api-key>
Content-Type: application/json
```

返回：

```json
{
  "image_url": "https://provider.example/images/xxx.png",
  "created_at": "2026-07-08T00:00:00Z"
}
```

API 会下载 `image_url` 并上传到 MinIO/S3。

### Nginx 反向代理示例

```nginx
server {
    listen 80;
    server_name api.labhaus.example;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 更新部署

```bash
git pull
pnpm install --frozen-lockfile
pnpm --filter @labhaus/web build

cd backend
go build -o labhaus-api ./cmd/api
```

之后按进程管理方式重启 Go API 和 Next.js 前端。

## 安全建议

1. 更换默认数据库、MinIO、JWT 和 Provider 密钥。
2. PostgreSQL、Redis、MinIO 不要直接暴露到公网。
3. API 和 Web 使用 HTTPS。
4. Provider API key 按环境隔离。
5. 对生产环境补充限流、审计日志、备份和监控。
6. 定期更新基础镜像和系统依赖。
