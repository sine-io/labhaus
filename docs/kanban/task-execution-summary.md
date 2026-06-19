# Labhaus 看板任务拆分与执行计划

生成时间: 2026-06-19
状态: ✅ 已完成拆分和标签分类

---

## 执行结果总结

### ✅ 第一步：任务原子化拆分

将 6 个复杂任务拆分为 **13 个原子任务**，每个任务：
- 单一职责（接口定义 / 算法实现 / API 封装）
- 独立验证（明确的测试命令）
- 3-5 小时完成时间
- 清晰的依赖关系

详细拆分分析见：`docs/kanban/task-breakdown-analysis.md`

---

### ✅ 第二步：Codex Lane 筛选与标签添加

根据 `sine-codex-lane-implementation` skill 标准，筛选出适合 Codex 实现的任务。

#### Codex Lane 任务（9个，69%）

标记为 `codex-lane-allowed` 的任务：

1. **#31** - C1: 样式推荐算法（TF-IDF + Cosine）
2. **#32** - D1: OpenAI GPT-4 剧本生成
3. **#33** - D3: Edge-TTS 配音生成
4. **#34** - A2: Mock ImageProvider 实现
5. **#35** - A3: GPT-Image-2 Provider 实现
6. **#37** - D2: 分镜设计与场景生成
7. **#38** - B1: 批量生图 Service 层
8. **#40** - B3: 批量生图 HTTP API
9. **#42** - D5: 文章到视频完整工作流 API

**筛选理由**：
- 外部 API 集成（OpenAI, GPT-Image-2, Edge-TTS）
- 复杂业务逻辑（算法实现、并发控制）
- 标准化 REST API 实现
- 重复性高的步骤实现

#### Hermes 直接实现任务（4个，31%）

**不标记** `codex-lane-allowed` 的任务：

1. **#30** - A1: ImageProvider 接口定义 ✅ `status:ready-for-dev`
2. **#36** - C2: 样式推荐 HTTP API
3. **#39** - B2: MinIO 图像存储集成
4. **#41** - D4: 图像并发生成编排

**筛选理由**：
- 基础架构设计（接口定义、架构框架）
- 扩展现有实现（MinIO 存储层）
- 简单编排逻辑
- 小范围 API 封装

---

### ✅ 第三步：严格遵守 Kanban 流程

#### 单一 Ready 规则 ✅

只有 **#30 (A1: ImageProvider 接口定义)** 标记为 `status:ready-for-dev`。

所有其他任务标记为 `status:inbox`，等待前置任务完成后逐个提升。

#### 标签完整性 ✅

每个任务都包含：
- ✅ `mvp:phase2` - 里程碑
- ✅ `backend:go` / `frontend:nextjs` - 技术栈
- ✅ `type:implementation` - 任务类型
- ✅ `priority:high` / `priority:medium` - 优先级
- ✅ `status:inbox` / `status:ready-for-dev` - 当前状态
- ✅ `codex-lane-allowed` - 实现方式（仅适用任务）

---

## 任务依赖图

```
第一批（可并行）
├── #30 A1: ImageProvider 接口 ✅ ready-for-dev
├── #31 C1: 样式推荐算法 🤖 codex
├── #32 D1: GPT-4 剧本生成 🤖 codex
└── #33 D3: Edge-TTS 配音 🤖 codex

第二批（依赖第一批）
├── #34 A2: Mock Provider 🤖 codex (依赖 #30)
├── #35 A3: GPT-Image-2 Provider 🤖 codex (依赖 #30)
├── #36 C2: 样式推荐 API (依赖 #31)
└── #37 D2: 分镜设计 🤖 codex (依赖 #32, #31)

第三批
├── #38 B1: 批量生图 Service 🤖 codex (依赖 #30, #34)
└── #39 B2: MinIO 存储 (依赖 #30)

第四批
├── #40 B3: 批量生图 API 🤖 codex (依赖 #38, #39)
└── #41 D4: 图像并发编排 (依赖 #38)

第五批
└── #42 D5: 文章到视频 API 🤖 codex (依赖 #32, #37, #33, #41)
```

---

## 实施策略

### 立即执行：#30 A1: ImageProvider 接口定义

**方式**: Hermes 直接 TDD 实现

**理由**:
- 基础架构定义，需要精确设计
- 影响所有后续 Provider 实现
- 接口设计不适合探索式开发

**目标**:
- 清晰的接口定义
- Registry 注册机制
- 单元测试验证接口契约

### 后续任务流程

对于每个任务：

1. **Ready Gate** - 运行 `sine-kanban-ready-gate` 验证
2. **实施**:
   - Codex Lane 任务：运行 `sine-codex-lane-implementation`
   - Hermes 任务：直接 TDD 实现
3. **Review Gate** - 运行 `sine-hermes-review-gate`
4. **Close & Sync** - 运行 `sine-issue-close-and-doc-sync`

---

## 统计数据

| 指标 | 数值 |
|------|------|
| 原始复杂任务 | 6 |
| 拆分原子任务 | 13 |
| Codex Lane 任务 | 9 (69%) |
| Hermes 直接实现 | 4 (31%) |
| 当前 Ready 任务 | 1 (#30) |
| 预计总工时 | 40-50 小时 |
| 最大并行度 | 4 任务 |

---

## 原任务映射

| 原任务 | 状态 | 拆分为 |
|--------|------|--------|
| #15 - 样式推荐与批量生图 | ready-for-dev | #31, #35, #36, #38, #39, #40 |
| #19 - Provider 接口抽象 | inbox | #30, #34, #35 |
| #20 - 批量生图并发控制 | inbox | #38, #39, #40 |
| #16 - 文章到视频工作流 | inbox | #32, #33, #37, #41, #42 |
| #29 - Workflow 执行引擎 | inbox | 暂缓（等待 #26/#27/#28） |
| #17 - 任务监控面板 | inbox | 暂缓（等待后端稳定） |

---

## 下一步行动

✅ 进入第三步：在 `/home/ubuntu/labhaus` 目录开始开发

**当前任务**: #30 A1: ImageProvider 接口定义与类型系统

**执行命令**:
```bash
cd /home/ubuntu/labhaus
# Hermes 直接 TDD 实现
```

---

## 流程遵守确认

✅ 严格遵守 `sine-requirements-to-kanban` skill  
✅ 严格遵守 `sine-codex-lane-implementation` skill  
✅ 单一 ready-for-dev 任务规则  
✅ 明确的依赖关系和验收标准  
✅ Codex Lane 标签正确标记  
✅ 所有任务包含完整标签集

