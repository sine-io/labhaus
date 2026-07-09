# 产品需求文档（PRD）

**项目名称**：Labhaus  
**版本**：v0.1.0  
**当前阶段**：图像素材 MVP 已实现，下一阶段扩展视频工作流
**最后更新**：2026-07-08

---

## 1. 产品定位

Labhaus 是一个 **可视化 AI 内容生产平台**，长期目标是让专业团队用工作流方式批量生产 AI 视频。

当前 MVP 先验证更小的闭环：

> 用户登录后，用样式推荐选择视觉方向，再批量生成可下载的 AI 图像素材。

这个闭环服务于后续视频能力：图像素材会成为“文章/文本 -> 分镜 -> 生图 -> TTS -> FFmpeg 合成 -> 视频”的中间产物。

## 2. 目标用户

### 内容创作者

- 需要稳定风格和高频素材输出。
- 当前痛点：手动调 prompt、反复生成、整理素材耗时。

### 营销和内容团队

- 需要围绕多个文案、产品或活动批量生成视觉素材。
- 当前痛点：多风格测试成本高，成功样式难复制。

### 开发者和自动化团队

- 需要可自部署、可扩展、可 API 化的内容生产能力。
- 当前痛点：从零搭建认证、任务、存储、Provider 适配成本高。

## 3. 当前 MVP 范围

### 3.1 已实现能力

| 能力               | 状态   | 说明                                                                          |
| ------------------ | ------ | ----------------------------------------------------------------------------- |
| 用户注册/登录      | 已实现 | 邮箱、密码、名称注册；JWT Bearer Token 登录                                   |
| 认证保护 API       | 已实现 | `/api/styles/*`、`/api/images/*`、`/api/workflows/*` 均需 Bearer Token        |
| 样式列表和详情     | 已实现 | PostgreSQL 存储，支持 category/tags/limit/offset 过滤                         |
| 样式创建           | 已实现 | 创建 name/description/prompt/category/tags                                    |
| 样式推荐           | 已实现 | 当前 HTTP 运行时使用关键词重叠打分；领域层已有 TF-IDF + Cosine 实现，后续接入 |
| 批量生图           | 已实现 | 多 prompt 并发调用图像 Provider，默认最大并发 10                              |
| 图像存储           | 已实现 | 下载 Provider 返回图片并上传 MinIO `images` bucket                            |
| 图片下载链接       | 已实现 | 返回 24 小时 MinIO 预签名 URL                                                 |
| 本地 mock Provider | 已实现 | Docker Compose 默认启动 mock image provider                                   |
| Smoke 脚本         | 已实现 | `infra/scripts/mvp-smoke.sh` 验证注册/登录/推荐/生图                                |

### 3.2 当前不包含

- 文章解析
- LLM 剧本生成
- TTS 配音
- FFmpeg 视频合成
- React Flow 可视化编辑器
- 模板市场
- OAuth / Refresh Token / RBAC / 配额计费
- OpenAPI 自动生成和生产级限流

这些是后续阶段目标，不应描述为当前已上线能力。

## 4. 核心用户流程

### 4.1 Web 首次使用

1. 打开 `http://localhost:3000`。
2. 进入 `/auth` 注册或登录。
3. 前端把登录返回的 token 保存到 `localStorage`。
4. 进入 `/styles/recommend`，输入创意描述，获取推荐样式。
5. 复制或参考推荐样式 prompt。
6. 进入 `/images/generate`，每行输入一个 prompt。
7. 前端带上 Bearer Token 调用 Next.js API 代理。
8. Go 后端调用图像 Provider 并写入 MinIO。
9. 用户在页面查看并下载预签名链接图片。

### 4.2 API 使用

1. `POST /api/users/register`
2. `POST /api/users/login` 获取 `token`
3. `POST /api/styles/recommend`
4. `POST /api/images/generate`

完整 API 见 `docs/architecture/api-design.md`。

## 5. 样式库策略

当前仓库内 `backend/seeds/styles.sql` 提供 12 条本地 demo 样式，覆盖 UI、retro、nature、cyberpunk、art、luxury 等类别。

产品目标仍是接入和维护 500+ GPT-Image-2 样式库，但当前 seed 不能表述为完整 500+ 样式库。

推荐策略分两层：

- **当前运行时**：启动时读取样式快照，HTTP handler 使用关键词重叠打分返回推荐。
- **后续目标**：接入 `internal/domain/style/recommendation` 中的 TF-IDF + Cosine 推荐实现，并补充更完整的样式数据。

## 6. 图像 Provider 策略

后端不内置真实生图模型，只通过 Provider 合约调用外部或本地服务：

```http
POST /v1/generate
Authorization: Bearer <api-key>
```

请求字段：

- `prompt`
- `width`
- `height`
- `quality`
- `style`

响应字段：

- `image_url`
- `created_at`

本地开发默认使用 `cmd/mock-image-provider`，生产环境应配置真实 Provider：

- `LABHAUS_IMAGE_PROVIDER_BASE_URL`
- `LABHAUS_IMAGE_PROVIDER_API_KEY`

## 7. 非功能需求

### 当前 MVP 目标

| 指标              | 目标         |
| ----------------- | ------------ |
| 本地 smoke 成功率 | 100%         |
| 单次批量 prompts  | 1-10         |
| 生图并发          | 10           |
| API 鉴权          | Bearer Token |
| 图片链接有效期    | 24 小时      |

### 后续生产化目标

| 指标                     | 目标                            |
| ------------------------ | ------------------------------- |
| 批量任务成功率           | > 90%                           |
| 单批 5-10 张素材完成时间 | < 10 分钟                       |
| 样式推荐质量             | 接入 TF-IDF + Cosine 后持续评估 |
| 可用性                   | MVP 99%，生产 99.9%             |

## 8. 后续路线

### Phase 2.x：视频工作流

- 文本/文章输入
- LLM 剧本和分镜生成
- 样式推荐辅助分镜生图
- TTS 配音
- FFmpeg 合成
- 任务状态和中间产物预览

### Phase 3：可视化工作流编辑器

- React Flow 画布
- 输入、LLM、生图、TTS、视频合成、输出节点
- 工作流保存和复用

### Phase 4：模板和团队能力

- 模板保存/分享
- 团队协作
- 配额、计费和私有化部署

## 9. 成功标准

当前 MVP 成功标准：

- 10 个种子用户能独立完成注册/登录、样式推荐和批量生图。
- 从创意描述到 5-10 张可下载图像素材小于 10 分钟。
- 至少 2 个种子用户确认愿意继续试用视频工作流版本。

---

**相关文档**：

- `docs/product/user-story-map.md`
- `docs/planning/mvp-roadmap.md`
- `docs/architecture/api-design.md`
- `docs/guides/quick-start.md`
