# 贡献指南

感谢你对 Labhaus 的关注。

## 当前主线

- 后端：`backend/` Go API
- 前端：`apps/web/` Next.js
- 本地依赖：PostgreSQL、Redis、MinIO、mock image provider

当前 MVP 是“认证后的样式推荐 + 批量生图”。视频工作流、可视化编辑器和模板市场是后续阶段。

## 开发环境

```bash
git clone https://github.com/sine-io/labhaus.git
cd labhaus
pnpm install
cp apps/web/.env.example apps/web/.env.local
cp backend/.env.example backend/.env
docker compose up -d --build
docker compose exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql
docker compose restart api
```

启动前端：

```bash
pnpm --filter @labhaus/web dev
```

裸跑后端：

```bash
cd backend
go run cmd/api/main.go
```

## 分支和提交

```bash
git checkout main
git pull
git checkout -b feature/your-feature-name
```

提交信息遵循 Conventional Commits：

```text
feat: add style recommendation runtime adapter
fix: forward authorization in web proxy
docs: update API contract
test: cover image generation response normalization
chore: update tooling
```

常用类型：

- `feat`
- `fix`
- `docs`
- `test`
- `refactor`
- `chore`

## 测试要求

提交前根据修改范围运行：

```bash
# 前端 helper 测试
pnpm --filter @labhaus/web test

# 前端类型检查
pnpm --filter @labhaus/web typecheck

# 前端 lint
pnpm --filter @labhaus/web lint

# 全 workspace
pnpm lint
pnpm test
pnpm build

# Go 后端
cd backend
go test ./...
go build ./...
```

本地 demo 验证：

```bash
scripts/mvp-smoke.sh
```

部分 Go 测试需要 PostgreSQL 或 MinIO 正在运行。

## 文档要求

如果修改 API、启动方式、配置、产品范围或目录结构，请同步更新：

- `README.md`
- `docs/architecture/api-design.md`
- `docs/guides/quick-start.md`
- `docs/guides/local-development.md`
- `docs/DEPLOYMENT.md`
- 相关产品或架构文档

当前仓库不保留历史执行报告类文档；过期内容请删除，不要继续引用旧 TypeScript API 或旧目录。

## Pull Request

PR 描述建议包含：

- 变更目的
- 主要文件
- 验证命令和结果
- 需要 reviewer 重点看的部分

Review 重点：

- 代码是否符合当前 Go + Next 主线
- API 契约是否清晰
- 鉴权、配置和错误处理是否合理
- 测试和文档是否同步

## 问题反馈

提交 issue 时请提供：

- 问题描述
- 复现步骤
- 期望行为
- 实际行为
- 环境信息（OS、Node、pnpm、Go、Docker 版本）
- 相关日志或截图

## 许可证

贡献的代码和文档将采用 [GNU Affero General Public License v3.0](LICENSE)。
