# Labhaus 技术栈文档

## 总览

Labhaus 采用 **TypeScript 全栈**架构，基于现代化的 Web 技术栈构建。

## 后端技术栈

### 核心框架
- **语言**: TypeScript 5.7+
- **运行时**: Node.js 20+
- **Web 框架**: Hono 4.7+ (轻量级高性能 Web 框架)
  - 为什么选择 Hono？
    - 比 Express 快 3-4 倍
    - TypeScript 原生支持
    - 中间件生态完善
    - 边缘计算友好

### 数据层
- **主数据库**: PostgreSQL 14+
  - 全文搜索 (FTS)
  - 数组类型支持
  - 事务保证
- **缓存**: Redis 7+ (规划中)
- **对象存储**: MinIO / Amazon S3
  - 图片、视频存储
  - CDN 集成

### 认证与安全
- **认证**: JWT (jsonwebtoken)
  - Access token (1小时过期)
  - Refresh token (7天过期)
- **密码加密**: bcrypt (10 rounds)
- **类型验证**: Zod
  - 运行时类型校验
  - 自动生成 TypeScript 类型

### NLP / 推荐
- **自然语言处理**: natural
  - TF-IDF 算法
  - 分词器
  - 余弦相似度计算

### 测试
- **测试框架**: Vitest
  - 快速的单元测试
  - 与 Vite 生态集成
- **E2E 测试**: 基于 Fetch API 的集成测试

### 开发工具
- **Monorepo**: Turborepo
  - 并行构建
  - 增量构建
  - 缓存优化
- **包管理器**: pnpm 9+
  - 磁盘空间节省
  - 严格的依赖管理
- **代码规范**: 
  - ESLint (代码检查)
  - Prettier (代码格式化)
  - TypeScript strict mode

## 前端技术栈（规划中）

### 核心框架
- **框架**: React 18
- **路由**: Next.js 14 (App Router)
  - 服务端渲染 (SSR)
  - 静态生成 (SSG)
  - API Routes

### UI 库
- **CSS 框架**: TailwindCSS 3+
  - 原子化 CSS
  - 响应式设计
  - 暗色模式支持
- **组件库**: shadcn/ui
  - 可定制组件
  - Radix UI 基础
  - Tailwind 风格

### 可视化
- **工作流编辑器**: React Flow
  - 拖拽节点
  - 连线交互
  - 自定义节点
- **图表**: Recharts / Chart.js

### 状态管理
- **全局状态**: Zustand
  - 轻量级
  - TypeScript 友好
  - 中间件支持

### 实时通信
- **服务器推送**: Server-Sent Events (SSE)
  - 任务进度更新
  - 实时通知

## AI / 媒体处理（规划中）

### 图像生成
- **OpenAI DALL-E 3** (API 集成)
- **Stable Diffusion** (可选)
- **Midjourney** (via Discord bot, 可选)

### 文本生成
- **OpenAI GPT-4** (剧本生成)
- **Prompt 工程**: Few-shot learning

### 语音合成
- **Edge-TTS** (免费、高质量)
- **OpenAI TTS** (备选)

### 视频处理
- **FFmpeg**
  - 图片合成
  - 音频混合
  - 字幕渲染
  - 格式转换

## 基础设施

### 容器化
- **Docker**
- **Docker Compose** (开发环境)
- **Kubernetes** (生产环境规划)

### CI/CD
- **GitHub Actions**
  - 自动测试
  - 类型检查
  - 代码规范检查
  - 自动部署 (规划中)

### 监控与日志（规划中）
- **日志**: Winston / Pino
- **APM**: Prometheus + Grafana
- **错误追踪**: Sentry

## 完整技术栈图

```
┌─────────────────────────────────────────────┐
│             前端层（规划中）                 │
│  React 18 + Next.js 14 + TailwindCSS        │
│  React Flow + Zustand + shadcn/ui           │
└─────────────────────────────────────────────┘
                     ↓ REST API
┌─────────────────────────────────────────────┐
│              API Gateway                     │
│  Hono + 中间件 (日志/认证/CORS/限流)        │
└─────────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────┐
│              业务服务层                      │
│  TypeScript + Zod + natural (NLP)           │
│  JWT + bcrypt (认证)                        │
└─────────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────┐
│             AI / 媒体处理                    │
│  OpenAI (GPT-4/DALL-E) + Edge-TTS + FFmpeg  │
└─────────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────┐
│              数据存储层                      │
│  PostgreSQL + Redis + MinIO/S3              │
└─────────────────────────────────────────────┘
```

## 项目结构

```
labhaus/
├── apps/
│   ├── api/              # 后端 API (Hono + TypeScript)
│   └── web/              # 前端 (Next.js, 规划中)
├── packages/
│   ├── types/            # 共享类型定义
│   ├── workflow/         # 工作流引擎
│   └── ui/               # UI 组件库 (规划中)
├── docs/                 # 文档
├── .github/              # GitHub Actions
└── docker-compose.yml    # 开发环境
```

## 依赖关系

### 后端核心依赖
```json
{
  "hono": "^4.7.9",
  "pg": "^8.14.0",
  "jsonwebtoken": "^9.0.2",
  "bcrypt": "^5.1.1",
  "zod": "^3.24.1",
  "natural": "^8.0.1"
}
```

### 开发依赖
```json
{
  "typescript": "^5.7.2",
  "vitest": "^2.1.8",
  "tsx": "^4.20.0",
  "turbo": "^2.3.3"
}
```

## 为什么选择这个技术栈？

### TypeScript 全栈
✅ **类型安全**: 编译时捕获错误  
✅ **开发效率**: 智能提示和重构  
✅ **代码质量**: 自文档化  
✅ **团队协作**: 接口即契约

### Hono over Express
✅ **性能**: 3-4倍更快  
✅ **现代化**: 原生 TypeScript  
✅ **轻量级**: 核心 < 13KB  
✅ **边缘友好**: 支持 Cloudflare Workers

### PostgreSQL
✅ **功能强大**: 全文搜索、数组类型  
✅ **可靠性**: ACID 事务保证  
✅ **扩展性**: 丰富的插件生态  
✅ **开源**: 无供应商锁定

### Turborepo + pnpm
✅ **构建速度**: 并行 + 增量构建  
✅ **磁盘节省**: pnpm 共享依赖  
✅ **Monorepo**: 统一管理多包  
✅ **缓存**: 本地 + 远程缓存

## 技术债务

### 当前限制
- ⚠️ 无 Redis 缓存（内存限制）
- ⚠️ 图像生成仅 Mock Provider
- ⚠️ 无分布式任务队列
- ⚠️ 无性能监控

### 后续优化
- [ ] 引入 Redis 缓存层
- [ ] 集成真实图像生成 API
- [ ] 使用 BullMQ 任务队列
- [ ] 添加 Prometheus 监控
- [ ] 实现 API 请求去重

## 更多文档

- [系统设计](./system-design.md)
- [API 设计](./api-design.md)
- [部署指南](../DEPLOYMENT.md)
