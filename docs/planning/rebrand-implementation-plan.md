# Labhaus 品牌升级实施计划

**Issue**: #1  
**创建日期**: 2026-06-18  
**预计完成**: 1-2 小时

---

## 📋 任务分解

### Task 1: 更新 GitHub 仓库 (10分钟)

**操作**：
```bash
# 1. 重命名仓库（需要手动在 GitHub 网页操作）
# Settings -> Repository name -> Rename to "labhaus"

# 2. 更新本地仓库
git remote set-url origin https://github.com/sine-io/labhaus.git

# 3. 更新仓库描述
gh repo edit --description "可视化 AI 内容生产平台 - 内容实验室 | Visual AI Content Workflow Platform"

# 4. 更新主题
gh repo edit --add-topic ai --add-topic workflow --add-topic content-production --add-topic video-generation
```

---

### Task 2: 更新项目文档 (30分钟)

#### 2.1 更新 README.md

**变更**：
- 项目名称：`Labhaus` → `Labhaus`
- 副标题：`可视化 AI 内容生产平台` → `内容实验室 - Where content experiments succeed`
- 口号：添加 `"Experiment. Automate. Scale."`
- 仓库链接：更新所有 GitHub URL

**示例**：
```markdown
# Labhaus

> 内容实验室 - Where content experiments succeed

[![Status](https://img.shields.io/badge/status-MVP%20开发中-blue)](https://github.com/sine-io/labhaus)

**"Experiment. Automate. Scale."**  
**"实验一次，复制千次"**

让专业团队在内容实验室中，3小时批量生产100个视频。
```

#### 2.2 更新 docs/ 所有文档

**需要更新的文件**：
- `docs/product/PRD.md`
- `docs/planning/mvp-roadmap.md`
- `docs/research/*.md`
- `docs/user-research/*.md`

**查找替换**：
```bash
# 批量替换项目名称
find docs -name "*.md" -type f -exec sed -i 's/Labhaus/Labhaus/g' {} +
find docs -name "*.md" -type f -exec sed -i 's/ai-content-pipeline/labhaus/g' {} +
```

---

### Task 3: 设计 Logo (需要人工)

#### 3.1 GPT-Image-2 Prompt

**Prompt 1: 主 Logo**
```
A modern minimalist logo for "Labhaus" - a content workflow platform.
Design concept: A laboratory beaker merged with mechanical gear/cog elements.
The beaker should be stylized and geometric, with the gear forming part of its structure.
Color palette: Technology blue (#0EA5E9) and experimental green (#10B981).
Style: Bauhaus, modernist, clean lines, vector-ready.
Professional tech startup aesthetic, suitable for app icon.
White or transparent background, high contrast.
Aspect ratio: 1:1 (square).
```

**Prompt 2: 简化图标版**
```
Simplified icon version of Labhaus logo.
Abstract geometric shape combining beaker silhouette with gear teeth.
Minimal, recognizable at 16x16 pixels.
Single color: Technology blue (#0EA5E9).
Perfect square, center aligned, suitable for favicon.
Style: flat design, modernist.
```

**Prompt 3: 横版 Wordmark**
```
Horizontal wordmark logo for "Labhaus".
Include the beaker-gear icon on the left, followed by "LABHAUS" text.
Typography: Modern sans-serif, geometric, bold but not heavy.
Icon and text in technology blue (#0EA5E9).
Professional, clean, suitable for website header.
Aspect ratio: 16:4 (wide).
```

#### 3.2 生成后处理

**需要生成的尺寸**：
- `logo-icon.svg` - 矢量图标（如果可能）
- `logo-512.png` - 512x512 (主图标)
- `logo-256.png` - 256x256
- `logo-128.png` - 128x128
- `logo-64.png` - 64x64
- `logo-32.png` - 32x32 (favicon)
- `logo-16.png` - 16x16 (favicon)
- `logo-horizontal.png` - 横版 wordmark (1200x300)

**存放位置**：
```
.github/
  brand/
    logo-icon.png
    logo-512.png
    logo-256.png
    logo-128.png
    logo-64.png
    logo-32.png
    logo-16.png
    logo-horizontal.png
    favicon.ico
    social-card.png (1200x630 for Open Graph)
```

---

### Task 4: 创建品牌资产 (20分钟)

#### 4.1 Favicon

```bash
# 使用 ImageMagick 或在线工具转换
# logo-32.png -> favicon.ico
convert logo-32.png logo-16.png favicon.ico
```

#### 4.2 Social Card

**尺寸**: 1200x630 (Open Graph)

**内容**：
- Labhaus Logo (居中或左侧)
- 标题: "Labhaus"
- 副标题: "Experiment. Automate. Scale."
- 背景: 渐变（蓝色到绿色）

#### 4.3 品牌色板文档

创建 `.github/brand/colors.md`:
```markdown
# Labhaus 品牌色板

## 主色

### Technology Blue
- HEX: #0EA5E9
- RGB: 14, 165, 233
- 用途: Logo主色、CTA按钮、链接

### Experimental Green
- HEX: #10B981
- RGB: 16, 185, 129
- 用途: 成功状态、强调色

## 中性色

### Dark Gray
- HEX: #1E293B
- RGB: 30, 41, 59
- 用途: 文本、背景

### Light Gray
- HEX: #F1F5F9
- RGB: 241, 245, 249
- 用途: 浅色背景

## 使用指南

- 主色用于品牌识别
- 绿色用于成功/完成状态
- 避免使用纯黑#000000，使用深灰代替
```

---

### Task 5: 更新所有品牌材料 (20分钟)

#### 5.1 更新 README.md Header

```markdown
<p align="center">
  <img src=".github/brand/logo-128.png" alt="Labhaus Logo" width="128" height="128">
</p>

<h1 align="center">Labhaus</h1>

<p align="center">
  <strong>内容实验室 - Where content experiments succeed</strong>
</p>

<p align="center">
  "Experiment. Automate. Scale." | "实验一次，复制千次"
</p>

<p align="center">
  <a href="https://github.com/sine-io/labhaus/stargazers">
    <img src="https://img.shields.io/github/stars/sine-io/labhaus?style=social" alt="Stars">
  </a>
  <a href="https://github.com/sine-io/labhaus/issues">
    <img src="https://img.shields.io/github/issues/sine-io/labhaus" alt="Issues">
  </a>
  <a href="https://github.com/sine-io/labhaus/blob/master/LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
  </a>
</p>
```

#### 5.2 更新文档首页

在所有主要文档顶部添加：
```markdown
# 📚 [文档名称]

**项目**: [Labhaus](https://github.com/sine-io/labhaus)  
**定位**: 内容实验室 - Where content experiments succeed
```

#### 5.3 更新 GitHub 仓库设置

在 GitHub 网页端：
1. Settings -> General
   - About section: 添加 website (如果有)
   - Social preview: 上传 social-card.png
2. Settings -> Options
   - Features: 启用 Discussions (可选)

---

## 🔄 执行流程

### 方式 1: 使用 Codex Lane（推荐）

```bash
# 1. 加载 sine-codex-lane-implementation skill
# 2. 执行 Codex Lane
codex-lane --issue 1

# Codex 会自动：
# - 创建新分支 feature/rebrand-to-labhaus
# - 更新所有文档
# - 提交变更
# - 创建 Pull Request
```

### 方式 2: 手动执行

```bash
# 1. 创建分支
git checkout -b feature/rebrand-to-labhaus

# 2. 批量替换
find . -name "*.md" -type f -exec sed -i 's/Labhaus/Labhaus/g' {} +
find . -name "*.md" -type f -exec sed -i 's/ai-content-pipeline/labhaus/g' {} +

# 3. 手动更新特殊文件（README.md 需要重新设计结构）

# 4. 创建品牌资产目录
mkdir -p .github/brand

# 5. 提交
git add -A
git commit -m "rebrand: 更名为 Labhaus

- 更新所有文档中的项目名称
- 更新品牌口号和定位
- 添加品牌资产目录结构
- 更新 GitHub 链接"

# 6. 推送并创建 PR
git push origin feature/rebrand-to-labhaus
gh pr create --title "品牌更名：Labhaus → Labhaus" --body "Closes #1"
```

---

## ✅ 验收清单

### 代码层面
- [ ] 所有 `.md` 文件中无 "Labhaus" 残留
- [ ] 所有 GitHub URL 已更新为 `sine-io/labhaus`
- [ ] README.md 展示新品牌形象
- [ ] 品牌资产目录已创建

### 品牌层面
- [ ] Logo 已生成（至少主图标）
- [ ] 品牌色板已文档化
- [ ] 口号和定位清晰展示

### GitHub 层面
- [ ] 仓库已重命名
- [ ] 仓库描述已更新
- [ ] Topics/标签已更新
- [ ] Social preview 已设置（如果 logo 已生成）

---

## 📝 注意事项

### 仓库重命名注意
1. GitHub 会自动设置重定向（旧 URL → 新 URL）
2. 本地仓库需要更新 remote URL
3. 已 fork 的仓库需要手动更新

### Logo 生成注意
1. GPT-Image-2 可能无法生成完美的 Logo
2. 建议生成多个版本后人工筛选
3. 必要时可以找设计师优化

### 文档更新注意
1. 使用 `git grep` 检查是否有遗漏
2. 特别注意 Issue/PR 模板中的链接
3. 检查 docs/ 中的相对链接是否正常

---

## 🚀 执行

准备好后，运行：

```bash
# 使用 Codex Lane
codex-lane --issue 1

# 或手动执行上述步骤
```

---

**文档维护**：
- 创建人：项目团队
- 最后更新：2026-06-18
- Issue: #1
