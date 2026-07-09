# Labhaus Web

`frontend` 是当前 Labhaus 前端主线，基于 Next.js App Router 实现 MVP 页面：

- `/`：MVP 入口页
- `/auth`：注册、登录、保存/清除 Bearer Token
- `/styles/recommend`：认证后的样式推荐
- `/images/generate`：认证后的批量生图

前端通过 App Router route handlers 代理请求到 Go API，并转发浏览器侧保存的 `Authorization` header。

## 环境变量

复制模板：

```bash
cp .env.example .env.local
```

当前使用的变量：

```bash
BACKEND_URL=http://localhost:8080
```

`BACKEND_URL` 是 Go API 地址。浏览器请求仍访问 Next.js 的相对路径，例如 `/api/users/login`，由 Next.js 服务端代理到 `BACKEND_URL`。

## 开发

从 `frontend/` 目录运行：

```bash
cd frontend
pnpm install
pnpm dev
```

访问 http://localhost:3000。

## 测试和检查

```bash
# 前端代理契约和 token helper 测试
pnpm test

# TypeScript 类型检查
pnpm typecheck

# ESLint
pnpm lint

# 生产构建
pnpm build
```

## API 代理契约

代理 helper 位于 `app/api/_lib/backend-contract.mjs`：

- `buildBackendHeaders()`：向 Go API 转发 `Content-Type` 和可选 `Authorization`。
- `normalizeGenerateImageResponse()`：把 Go 后端的 `results` 同步暴露为前端兼容的 `images`。
- `normalizeRecommendStyleRequest()`：把旧字段 `prompt`/`top_k` 兼容转换为 Go 后端当前使用的 `query`/`limit`。

Bearer Token helper 位于 `app/lib/auth-token.mjs`，浏览器使用 `localStorage` 键 `labhaus_bearer_token` 保存 token。
