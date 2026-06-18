# AI Content Pipeline

> 可视化 AI 内容生产平台 - 让非技术人员也能批量生产 AI 视频

[![Status](https://img.shields.io/badge/status-MVP%20开发中-blue)](https://github.com/sine-io/ai-content-pipeline)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

---

## 🎯 项目定位

**"Canva for AI Video Workflow"**

让非技术人员也能搭建自己的内容生产流水线：
- ✅ **零代码**：拖拽式工作流编辑器
- ✅ **批量生产**：一次配置，输出 10/100 个视频
- ✅ **样式复用**：500+ 工业级 GPT-Image-2 提示词库
- ✅ **模板市场**：用户共享/购买成功配方

---

## 📊 项目现状

**阶段**：需求验证 → MVP 开发准备

**时间线**：
- ✅ 市场调研完成（2026-06-18）
- ✅ 方法论分析完成（真需求、双钻模型、JTBD）
- 🔄 用户访谈进行中（目标：3-5 个种子用户）
- ⏳ MVP Phase 1 启动（预计 2 周）

---

## 🚀 核心价值

### 对标分析

| 维度 | ComfyUI | Runway | Synthesia | **我们** |
|------|---------|--------|-----------|---------|
| **易用性** | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **可扩展性** | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| **批量能力** | ⭐⭐⭐ | ⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| **样式复用** | ⭐⭐ | ⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **模板市场** | ❌ | ❌ | ❌ | ✅ |
| **私有部署** | ✅ | ❌ | 💰 | ✅ |

### 核心差异化

1. **ComfyUI 的强大** + **Canva 的易用** + **Gumroad 的社区**
2. **500+ 样式库**（来自 awesome-gpt-image-2）
3. **模板市场**（网络效应）

---

## 👥 目标用户

### 主要用户（MVP 聚焦）

#### 1. 内容创作者（B2C）
**痛点**：
- 批量做视频太慢（手动 2-4 小时/条）
- 保持风格一致性难
- 外包成本高（￥100-500/条）

**价值**：
- 3 小时批量生成 10 条
- 单条成本 < ￥50
- 风格一致性 > 90%

#### 2. 企业营销团队（B2B）
**痛点**：
- 需要大量营销素材
- 测试多种创意成本高
- 外包周期长

**价值**：
- CSV 批量导入
- A/B 测试多种风格
- 成本降低 80%

#### 3. 开发者（B2D）
**痛点**：
- 从零开发周期长（3-6 个月）
- 调用 Runway API 贵且黑盒
- 开源方案维护成本高

**价值**：
- 1 周完成集成
- API 文档完善
- 深度可定制

---

## 🏗️ 技术架构

### 整体架构

```
ai-content-pipeline/
├── packages/
│   ├── core/                    # 核心引擎
│   │   ├── workflow-engine/    # 工作流编排
│   │   ├── state-machine/      # 任务状态机
│   │   └── provider-sdk/       # Provider 抽象层
│   │
│   ├── style-library/          # 样式库服务
│   │   ├── api/                # RESTful API
│   │   ├── data/               # 500+ 案例数据
│   │   └── engine/             # 样式匹配算法
│   │
│   ├── services/               # 微服务
│   │   ├── image-gen/          # 图像生成
│   │   ├── video-render/       # 视频渲染
│   │   ├── script-gen/         # 剧本生成
│   │   └── template-market/    # 模板市场
│   │
│   └── apps/
│       ├── web/                # Web 主应用
│       │   ├── editor/         # 🔥 可视化编辑器
│       │   ├── gallery/        # 样式库展示
│       │   ├── marketplace/    # 🔥 模板市场
│       │   └── dashboard/      # 任务监控
│       │
│       └── api-gateway/        # 统一 API 网关
│
├── infra/                      # 基础设施
│   ├── docker/
│   ├── k8s/
│   └── terraform/
│
└── docs/                       # 文档
    ├── product/                # 产品文档
    ├── research/               # 调研报告
    ├── architecture/           # 架构设计
    └── planning/               # 规划文档
```

### 技术栈

**前端**
- React + TypeScript
- React Flow（节点编辑器）
- Zustand（状态管理）
- TailwindCSS + shadcn/ui

**后端**
- Python FastAPI
- Celery + Redis（异步任务）
- PostgreSQL（主数据库）
- MinIO / S3（对象存储）

**工作流引擎**
- Temporal / 自研状态机
- JSON Schema（工作流定义）

---

## 📅 MVP 路线图

### Phase 1: 基础整合（2 周）
**目标**：两个项目代码合并，建立统一架构

- [ ] 创建 Monorepo 结构
- [ ] awesome-gpt-image-2 样式库提取为服务
- [ ] ai-video-factory 状态机改造为通用引擎
- [ ] 统一 API 网关和认证

**交付物**：
- 统一项目结构
- 样式库 API（GET /styles）
- 工作流引擎 API（POST /workflows/execute）

### Phase 2: 核心工作流（4 周）🔥
**目标**：实现第一个完整的 MVP 工作流

**2.1 样式库 API 化（1 周）**
- [ ] 样式查询接口
- [ ] 样式推荐算法
- [ ] 批量生图接口

**2.2 "文章 → 视频" 工作流（2 周）**
- [ ] 剧本生成（LLM）
- [ ] 分镜设计
- [ ] 批量图像生成（GPT-Image-2 + 样式库）
- [ ] 配音生成（TTS）
- [ ] 视频合成（FFmpeg）

**2.3 任务监控面板（1 周）**
- [ ] 实时任务状态
- [ ] 中间产物预览
- [ ] 错误日志

**交付物**：
- 完整的"文章 → 视频" demo
- 任务监控面板
- API 文档

### Phase 3: 可视化编辑器（6 周）🔥🔥🔥
**目标**：核心差异化功能

**3.1 节点编辑器基础（2 周）**
- [ ] React Flow 集成
- [ ] 节点类型定义
- [ ] 数据流验证
- [ ] 保存/加载工作流

**3.2 内置节点库（2 周）**
- [ ] 输入节点（文本/URL/CSV）
- [ ] 处理节点（LLM/样式选择/生图/TTS/合成）
- [ ] 输出节点（下载/S3/Telegram）

**3.3 人工介入和预览（2 周）**
- [ ] 中间结果预览
- [ ] 人工审核节点
- [ ] 手动修改和重试

**交付物**：
- 拖拽式工作流编辑器
- 10+ 内置节点
- 工作流模板保存

### Phase 4: 模板市场（MVP+，4 周）
**目标**：商业化核心

- [ ] 模板保存/加载
- [ ] 模板分类和搜索
- [ ] 模板分享（公开/私有）
- [ ] 模板评分和评论
- [ ] 付费模板（Stripe）

---

## 💰 商业模式

### 定价

#### 免费版（Freemium）
- 基础工作流功能
- 公共样式库访问
- 单任务串行执行
- 月生成配额：10 个视频
- 社区模板使用

#### 专业版（$29/月）
- 批量并发任务（10 并发）
- 私有样式库
- 高级节点（自定义脚本）
- 月生成配额：100 个视频
- API 调用（1000 次/月）
- 优先支持

#### 企业版（定制报价）
- 私有化部署
- 无限并发
- 自定义 provider 集成
- 白标服务
- 专属技术支持
- SLA 保障

### 收入来源

1. **订阅费用**（主要）
2. **模板市场分成**（30%）
3. **API 调用费**（超出配额）
4. **专业服务**（定制开发、培训）

---

## 📚 文档索引

### 产品文档
- [产品需求文档 (PRD)](docs/product/PRD.md)
- [用户故事地图](docs/product/user-story-map.md)
- [功能规格说明](docs/product/feature-specs.md)

### 调研报告
- [竞品分析报告](docs/research/competitive-analysis.md)
- [市场调研总结](docs/research/market-research.md)
- [方法论分析报告](docs/research/methodology-analysis.md)

### 架构设计
- [系统架构设计](docs/architecture/system-design.md)
- [数据库设计](docs/architecture/database-schema.md)
- [API 接口设计](docs/architecture/api-design.md)

### 规划文档
- [MVP 开发计划](docs/planning/mvp-roadmap.md)
- [技术选型说明](docs/planning/tech-stack.md)
- [风险管理计划](docs/planning/risk-management.md)

### 用户研究
- [用户访谈计划](docs/user-research/interview-plan.md)
- [用户访谈记录](docs/user-research/interview-notes.md)
- [用户画像](docs/user-research/user-personas.md)

---

## 🤝 贡献指南

项目目前处于 MVP 阶段，暂不接受外部贡献。

待 Beta 版本发布后，我们将开放以下贡献方式：
- 🐛 Bug 报告
- 💡 功能建议
- 📝 文档改进
- 🎨 样式库贡献

---

## 📄 许可证

[MIT License](LICENSE)

---

## 📧 联系方式

- **项目主页**：https://github.com/sine-io/ai-content-pipeline
- **Issues**：https://github.com/sine-io/ai-content-pipeline/issues

---

**最后更新**：2026-06-18  
**当前版本**：0.1.0-alpha（需求验证阶段）
