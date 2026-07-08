# Labhaus 部署指南

当前主线部署对象：

- `backend/`：Go API，默认端口 `8080`。
- `apps/web/`：Next.js 前端，默认端口 `3000`。
- 依赖服务：PostgreSQL、Redis、MinIO、图像 Provider。

早期 TypeScript API 和共享包遗留代码已删除；如需参考旧部署方案，请从 git history 查看。

## 环境要求

- Linux 服务器（推荐 Ubuntu 22.04+）
- Docker 和 Docker Compose（推荐）
- Node.js 20+ 和 pnpm 9+（构建/运行前端需要）
- Go 1.25+（裸跑 Go API 需要）

## 方式 1：Docker Compose（推荐开发/演示）

```bash
git clone https://github.com/sine-io/labhaus.git
cd labhaus

# 启动依赖、mock image provider 和 Go API
docker compose up -d --build

# 导入样式种子数据
docker compose exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql

# API 启动时会加载样式快照；导入后重启 API
docker compose restart api
```

访问：

- API: http://localhost:8080
- MinIO Console: http://localhost:9001
- Mock Image Provider: http://localhost:8089

前端可单独运行：

```bash
pnpm install
pnpm --filter @labhaus/web dev
```

访问 Web: http://localhost:3000

## 方式 2：裸机运行 Go API

### 1. 准备依赖服务

```bash
docker compose up -d postgres redis minio mock-image-provider
docker compose exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql
```

### 2. 配置后端环境变量

```bash
cp backend/.env.example backend/.env
```

最小必需配置：

```bash
LABHAUS_SERVER_PORT=8080
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
LABHAUS_IMAGE_PROVIDER_BASE_URL=http://localhost:8089
LABHAUS_IMAGE_PROVIDER_API_KEY=dev-mock-key
```

### 3. 启动 Go API

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

`apps/web/.env.local` 至少需要指向 Go API：

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

## 生产部署建议

### Go API

- 使用 systemd、容器平台或进程管理器托管编译后的 Go 二进制。
- 必须显式配置 `LABHAUS_IMAGE_PROVIDER_BASE_URL` 和 `LABHAUS_IMAGE_PROVIDER_API_KEY`。
- 生产环境应使用强随机 `LABHAUS_JWT_SECRET_KEY`。

示例：

```bash
cd backend
go build -o labhaus-api ./cmd/api
LABHAUS_SERVER_ENVIRONMENT=production ./labhaus-api
```

### Next.js 前端

- 使用 `pnpm --filter @labhaus/web build` 生成生产构建。
- 使用 `pnpm --filter @labhaus/web start` 或部署到支持 Next.js 的平台。
- 通过 `NEXT_PUBLIC_API_BASE_URL` 指向公开可访问的 Go API 地址。

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

之后按你的进程管理方式重启 Go API 和 Next.js 前端。

## 安全建议

1. 更换所有默认密钥和密码。
2. 限制 PostgreSQL、Redis、MinIO 只允许可信网络访问。
3. 给 API 和前端配置 HTTPS。
4. 定期更新基础镜像和系统依赖。
5. 生产图像 Provider 使用独立 API key，并按环境隔离。
