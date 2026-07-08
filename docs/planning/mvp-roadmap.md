# MVP 开发路线图

**项目**：Labhaus  
**版本**：v0.1.0
**最后更新**：2026-07-08

---

## 总体目标

Labhaus 的长期目标是 AI 视频工作流实验室；当前路线先以“样式推荐 + 批量生图”验证素材生产闭环，再扩展到“文本/文章 -> 视频”。

## Phase 1：主线收敛（已完成）

- [x] 明确 `backend/` Go API 为唯一后端主线。
- [x] 明确 `apps/web/` Next.js 为唯一前端主线。
- [x] 删除早期 TypeScript API、共享包和历史执行文档。
- [x] 使用 pnpm workspace + Turborepo 管理前端。
- [x] 使用 Docker Compose 管理本地依赖服务。

## Phase 2：图像素材 MVP（已实现）

目标：用户能登录后完成“样式推荐 -> 批量生图 -> 下载素材”。

### 已完成

- [x] 用户注册和登录。
- [x] JWT Bearer Token 鉴权。
- [x] Next.js 登录/注册页面。
- [x] Next.js API 代理转发 Authorization。
- [x] 样式列表、详情、创建接口。
- [x] 样式推荐接口 `POST /api/styles/recommend`。
- [x] 样式推荐页面。
- [x] 图像 Provider 抽象和 HTTP 适配。
- [x] 本地 mock image provider。
- [x] 批量生图服务，默认最大并发 10。
- [x] 批量生图接口 `POST /api/images/generate`。
- [x] MinIO image storage 和预签名 URL。
- [x] 批量生图页面。
- [x] 本地 seed：`backend/seeds/styles.sql`。
- [x] Smoke 脚本：`scripts/mvp-smoke.sh`。

### 当前限制

- 当前 seed 为 12 条 demo 样式；500+ 样式库仍是后续导入目标。
- HTTP 运行时推荐器当前使用关键词重叠打分；`domain/style/recommendation` 中已有 TF-IDF + Cosine 实现，后续应接入。
- Workflow API 只管理元数据和状态，不执行视频工作流。
- Redis queue worker 只保留任务处理骨架。
- 尚未实现 OAuth、refresh token、RBAC、配额、限流、OpenAPI 自动生成。

### 收尾验证

- [ ] 在具备 Go toolchain 的环境运行 `go test ./...`。
- [ ] 在 Docker 环境运行 `scripts/mvp-smoke.sh`。
- [ ] 用 3-5 个种子用户验证 Web 流程。

## Phase 2.x：视频工作流（下一阶段）

目标：把当前图像素材能力扩展为完整视频生成流程。

### 里程碑 1：工作流执行内核

- [ ] 定义视频工作流状态和任务模型。
- [ ] 将 Redis queue worker 接入真实工作流执行。
- [ ] 记录每一步输入、输出、错误和耗时。
- [ ] 提供任务详情查询接口。

### 里程碑 2：文本/文章到分镜

- [ ] 文本/文章输入接口。
- [ ] LLM 剧本生成。
- [ ] 分镜结构化输出。
- [ ] 分镜 prompt 与样式推荐结合。

### 里程碑 3：媒体生成和合成

- [ ] 分镜图批量生成。
- [ ] TTS 配音。
- [ ] 字幕时间轴。
- [ ] FFmpeg 合成 MP4。
- [ ] 中间产物存储到 MinIO。

### 里程碑 4：任务监控 UI

- [ ] 任务列表。
- [ ] 任务详情。
- [ ] 进度和错误展示。
- [ ] 图片、音频、视频预览。
- [ ] 失败步骤重试。

## Phase 3：可视化工作流编辑器

目标：让用户能拖拽组合工作流，而不只使用固定流程。

- [ ] React Flow 画布。
- [ ] 节点注册系统。
- [ ] 输入节点：文本、URL、CSV。
- [ ] 处理节点：LLM、样式推荐、生图、TTS、视频合成。
- [ ] 输出节点：下载、S3/MinIO、Webhook。
- [ ] 工作流 JSON 保存和加载。
- [ ] 基础 DAG 校验和参数校验。

## Phase 4：模板、团队和商业化

- [ ] 模板保存和复用。
- [ ] 模板导入/导出。
- [ ] 团队空间。
- [ ] 配额和计费。
- [ ] 私有化部署文档。
- [ ] 模板市场。

## 工程质量路线

### 当前验证命令

```bash
pnpm --filter @labhaus/web test
pnpm --filter @labhaus/web typecheck
pnpm lint
pnpm build

cd backend
go test ./...
go build ./...
```

### 后续补强

- [ ] 后端 CI 中加入 Go toolchain、Go test 和 Go build。
- [ ] 增加 Docker Compose smoke job。
- [ ] 增加 API contract 测试。
- [ ] 为视频工作流增加集成测试。

## 成功标准

### 图像素材 MVP

- 用户能在 10 分钟内完成登录、样式推荐和批量生图。
- 生成图片可打开或下载。
- 本地 demo 不依赖真实生图服务即可跑通。

### 视频工作流阶段

- 用户能从文本输入得到可预览 MP4。
- 每一步中间产物可追踪。
- 失败步骤可重试。
- 成功工作流可保存和复用。
