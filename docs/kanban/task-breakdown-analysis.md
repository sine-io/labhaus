# Labhaus 看板任务原子化拆分分析

生成时间: 2026-06-19
目的: 将 Phase 2 的 6 个复杂任务拆分为可独立实现的原子任务

---

## 原任务清单

### 高优先级任务
1. **#15** - Phase 2.1: 样式推荐算法与批量生图接口 (status:ready-for-dev)
2. **#19** - Phase 2.1.2: 图像生成 Provider 接口抽象 (status:inbox)
3. **#20** - Phase 2.1.3: 批量生图接口与并发控制 (status:inbox, codex-lane-allowed)
4. **#16** - Phase 2.2: 文章到视频工作流实现 (status:inbox, codex-lane-allowed)

### 中优先级任务
5. **#29** - Go 后端：Workflow 执行引擎实现 (status:inbox, 依赖 #26/#27/#28)
6. **#17** - Phase 2.3: 任务监控面板（前端） (status:inbox)

---

## 拆分策略

### 原则
- **最小可控单元**: 每个任务聚焦单一领域（接口定义 / 实现 / 测试）
- **独立可验证**: 有明确的验收标准和测试命令
- **依赖清晰**: 显式标注前置依赖
- **3-5 小时完成**: 单个任务适合 Codex 或 Hermes 一次性完成

---

## 拆分结果

### 组 A: Provider 接口层（基础设施）

#### A1: ImageProvider 接口定义与类型系统
**父任务**: #19  
**范围**:
- 定义 Go 版本的 `ImageProvider` 接口
- 定义 `ImageOptions` 和 `ImageResult` 结构体
- Provider 注册机制（registry pattern）
- 基础单元测试

**不做**:
- 不实现任何具体 Provider
- 不涉及 HTTP API

**验收**:
```bash
go test ./internal/domain/image/provider -v
go build ./internal/domain/image/...
```

**预计时间**: 2-3h  
**适合**: Hermes 直接实现（基础接口定义）

---

#### A2: Mock ImageProvider 实现
**父任务**: #19  
**依赖**: A1  
**范围**:
- 实现 `MockImageProvider`（生成占位图或返回固定 base64）
- 支持延迟模拟（用于测试超时）
- 支持错误模拟（用于测试重试）
- 单元测试覆盖所有模式

**验收**:
```bash
go test ./internal/infrastructure/image/mock -v
```

**预计时间**: 2-3h  
**适合**: Codex Lane（测试用例编写密集）

---

#### A3: GPT-Image-2 Provider 实现
**父任务**: #15 (部分), #19  
**依赖**: A1  
**范围**:
- HTTP 客户端封装
- 错误处理和重试逻辑
- 超时控制
- 集成测试（mock HTTP server）

**验收**:
```bash
go test ./internal/infrastructure/image/gptimage2 -v
```

**预计时间**: 3-4h  
**适合**: Codex Lane（HTTP 集成复杂）

---

### 组 B: 批量生图与并发控制

#### B1: 批量生图 Service 层实现
**父任务**: #20, #15  
**依赖**: A1, A2  
**范围**:
- `BatchImageService` 实现
- 并发控制（10 并发 semaphore）
- 任务队列管理（内存队列即可）
- 进度追踪接口
- 单元测试（使用 MockProvider）

**不做**:
- 不涉及 MinIO 存储（B2 负责）
- 不涉及 HTTP API（B3 负责）

**验收**:
```bash
go test ./internal/application/service/image -v
```

**预计时间**: 4-5h  
**适合**: Codex Lane（并发控制逻辑复杂）

---

#### B2: MinIO 图像存储集成
**父任务**: #20, #15  
**依赖**: A1  
**范围**:
- 扩展现有 `storage.Storage` 接口支持图像元数据
- 图像上传/下载
- URL 签名生成
- 集成测试（真实 MinIO 或 testcontainer）

**验收**:
```bash
go test ./internal/infrastructure/storage -tags=integration -v
```

**预计时间**: 3h  
**适合**: Hermes 直接实现（存储层已有基础）

---

#### B3: 批量生图 HTTP API
**父任务**: #20, #15  
**依赖**: B1, B2  
**范围**:
- `POST /api/images/generate` 接口
- `GET /api/images/:id` 查询接口
- `GET /api/images/:id/progress` 进度查询
- 请求验证
- 端到端测试

**验收**:
```bash
go test ./internal/infrastructure/http/handler/image -v
curl -X POST http://localhost:8080/api/images/generate -d @test.json
```

**预计时间**: 3h  
**适合**: Codex Lane（API 实现标准化）

---

### 组 C: 样式推荐算法

#### C1: 样式推荐算法实现
**父任务**: #15  
**依赖**: 无（独立）  
**范围**:
- TF-IDF 关键词提取
- Cosine Similarity 计算
- 推荐排序逻辑
- 单元测试（已知数据集验证准确率）

**验收**:
```bash
go test ./internal/domain/style/recommendation -v
# 验证准确率 > 70%
```

**预计时间**: 4-5h  
**适合**: Codex Lane（算法实现）

---

#### C2: 样式推荐 HTTP API
**父任务**: #15  
**依赖**: C1  
**范围**:
- `POST /api/styles/recommend` 接口
- 集成现有样式库数据
- 端到端测试

**验收**:
```bash
go test ./internal/infrastructure/http/handler/style -v
curl -X POST http://localhost:8080/api/styles/recommend -d '{"keywords": ["modern", "office"]}'
```

**预计时间**: 2h  
**适合**: Hermes 直接实现（简单 API 封装）

---

### 组 D: 文章到视频工作流

#### D1: OpenAI GPT-4 剧本生成 Service
**父任务**: #16  
**依赖**: 无  
**范围**:
- OpenAI API 客户端
- Prompt 模板管理
- JSON 结构化输出解析
- 错误处理和重试
- 单元测试 + 集成测试

**验收**:
```bash
go test ./internal/application/service/script -v
```

**预计时间**: 4h  
**适合**: Codex Lane（外部 API 集成）

---

#### D2: 分镜设计与场景生成
**父任务**: #16  
**依赖**: D1, C1  
**范围**:
- 剧本解析逻辑
- 场景描述生成
- 样式库匹配
- 分镜预览数据结构
- 单元测试

**验收**:
```bash
go test ./internal/domain/storyboard -v
```

**预计时间**: 3-4h  
**适合**: Codex Lane（业务逻辑复杂）

---

#### D3: Edge-TTS 配音生成 Service
**父任务**: #16  
**依赖**: 无  
**范围**:
- Edge-TTS CLI 封装
- 配音参数配置
- 字幕生成（SRT/VTT）
- 集成测试

**验收**:
```bash
go test ./internal/infrastructure/tts -v
```

**预计时间**: 3h  
**适合**: Codex Lane（外部工具集成）

---

#### D4: 图像并发生成编排
**父任务**: #16  
**依赖**: B1  
**范围**:
- 调用批量生图 Service
- 进度追踪
- 失败重试策略
- 结果质量检查（可选）

**验收**:
```bash
go test ./internal/application/service/workflow/image_generation -v
```

**预计时间**: 2-3h  
**适合**: Hermes 直接实现（编排逻辑）

---

#### D5: 文章到视频完整工作流 API
**父任务**: #16  
**依赖**: D1, D2, D3, D4  
**范围**:
- `POST /api/workflows/article-to-video` 接口
- 异步任务提交
- 状态查询接口
- 端到端测试

**验收**:
```bash
go test ./internal/infrastructure/http/handler/workflow -v
curl -X POST http://localhost:8080/api/workflows/article-to-video
```

**预计时间**: 3h  
**适合**: Codex Lane（集成测试复杂）

---

### 组 E: Workflow 执行引擎

#### E1: Workflow Executor 基础框架
**父任务**: #29  
**依赖**: #26, #27, #28（Redis/MinIO 基础设施）  
**范围**:
- WorkflowExecutor 接口定义
- 步骤编排逻辑
- 状态更新机制
- 错误处理和回滚钩子
- 单元测试

**验收**:
```bash
go test ./internal/application/service/workflow/executor -v
```

**预计时间**: 4h  
**适合**: Hermes 直接实现（架构设计）

---

#### E2: Workflow 步骤实现
**父任务**: #29  
**依赖**: E1, D1-D4  
**范围**:
- 文本处理步骤
- 图像生成步骤
- 视频合成步骤（FFmpeg）
- 上传结果步骤
- 集成测试

**验收**:
```bash
go test ./internal/application/service/workflow/steps -v
```

**预计时间**: 4-5h  
**适合**: Codex Lane（步骤实现重复性高）

---

### 组 F: 前端监控面板

#### F1: Next.js + React 项目初始化
**父任务**: #17  
**依赖**: 无  
**范围**:
- Next.js App Router 初始化
- TailwindCSS 配置
- API 客户端封装
- 基础 Layout 组件

**验收**:
```bash
npm run build
npm run dev
```

**预计时间**: 2h  
**适合**: Hermes 直接实现（脚手架搭建）

---

#### F2: 任务列表与实时状态组件
**父任务**: #17  
**依赖**: F1  
**范围**:
- 任务列表页面
- 分页和筛选
- SSE/WebSocket 集成
- 响应式设计

**验收**:
```bash
npm run test
npm run build
```

**预计时间**: 5h  
**适合**: Codex Lane（前端组件实现）

---

#### F3: 任务详情与进度可视化
**父任务**: #17  
**依赖**: F1, F2  
**范围**:
- 进度条组件
- 中间产物预览
- 错误日志展示
- 视频播放器集成

**验收**:
```bash
npm run test
npm run build
```

**预计时间**: 4h  
**适合**: Codex Lane（UI 组件复杂）

---

## 任务优先级排序（依赖拓扑）

### 第一批（可并行）
1. **A1** - ImageProvider 接口定义（基础）
2. **C1** - 样式推荐算法（独立）
3. **D1** - GPT-4 剧本生成（独立）
4. **D3** - Edge-TTS 配音（独立）

### 第二批（依赖第一批）
5. **A2** - Mock Provider（依赖 A1）
6. **A3** - GPT-Image-2 Provider（依赖 A1）
7. **C2** - 样式推荐 API（依赖 C1）
8. **D2** - 分镜设计（依赖 D1, C1）

### 第三批
9. **B1** - 批量生图 Service（依赖 A1, A2）
10. **B2** - MinIO 存储（依赖 A1）

### 第四批
11. **B3** - 批量生图 API（依赖 B1, B2）
12. **D4** - 图像并发编排（依赖 B1）

### 第五批
13. **D5** - 文章到视频 API（依赖 D1-D4）

### 第六批（等待基础设施）
14. **E1** - Workflow Executor 框架（依赖 #26/#27/#28）
15. **E2** - Workflow 步骤实现（依赖 E1, D1-D4）

### 第七批（前端可并行）
16. **F1** - Next.js 初始化（独立）
17. **F2** - 任务列表组件（依赖 F1）
18. **F3** - 详情页组件（依赖 F1, F2）

---

## Codex Lane 筛选标准

基于 `sine-codex-lane-implementation` skill，标记适合 Codex 的任务：

### ✅ 必须使用 Codex Lane
- **A3** - GPT-Image-2 Provider（HTTP 集成 + 复杂错误处理）
- **B1** - 批量生图 Service（并发控制 + 队列管理）
- **B3** - 批量生图 API（标准 REST API）
- **C1** - 样式推荐算法（TF-IDF + Cosine Similarity）
- **D1** - GPT-4 剧本生成（外部 API + 结构化输出）
- **D2** - 分镜设计（复杂业务逻辑）
- **D3** - Edge-TTS 配音（外部工具集成）
- **D5** - 文章到视频 API（多步骤集成）
- **E2** - Workflow 步骤实现（重复性步骤实现）
- **F2** - 任务列表组件（前端 SSE/WebSocket）
- **F3** - 详情页组件（复杂 UI 组件）

### ⚙️ Hermes 直接实现
- **A1** - ImageProvider 接口（基础接口定义）
- **A2** - Mock Provider（简单实现）
- **B2** - MinIO 存储（扩展现有实现）
- **C2** - 样式推荐 API（简单 API 封装）
- **D4** - 图像并发编排（编排逻辑）
- **E1** - Workflow Executor 框架（架构设计）
- **F1** - Next.js 初始化（脚手架）

---

## 下一步行动

1. ✅ 生成拆分后的原子任务列表
2. ⏳ 在 GitHub 创建对应的 Issues
3. ⏳ 添加 `codex-lane-allowed` 标签
4. ⏳ 只有第一个任务标记为 `status:ready-for-dev`
5. ⏳ 开始实现

---

## 总结

- **总任务数**: 18 个原子任务（拆分自 6 个复杂任务）
- **Codex Lane 任务**: 11 个（61%）
- **Hermes 直接实现**: 7 个（39%）
- **预计总时间**: 60-70 小时
- **并行度**: 最高 4 个任务可并行（第一批）

