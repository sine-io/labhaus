# 贡献指南

感谢你对 Labhaus 项目的关注！

## 贡献方式

- 🐛 报告 Bug
- 💡 提出功能建议
- 📝 改进文档
- 🛠️ 提交代码

## 开发流程

### 1. Fork 和克隆

```bash
# Fork 项目后克隆你的 fork
git clone https://github.com/YOUR_USERNAME/labhaus.git
cd labhaus

# 添加上游仓库
git remote add upstream https://github.com/sine-io/labhaus.git
```

### 2. 创建分支

```bash
# 更新 main 分支
git checkout main
git pull upstream main

# 创建功能分支
git checkout -b feature/your-feature-name
```

### 3. 开发

参考 [本地开发指南](docs/guides/local-development.md) 搭建环境

```bash
pnpm install
docker compose up -d
pnpm dev
```

### 4. 提交

```bash
# 暂存修改
git add .

# 提交（使用清晰的 commit message）
git commit -m "feat: add new feature"

# 推送到你的 fork
git push origin feature/your-feature-name
```

### 5. 提交 Pull Request

1. 前往 https://github.com/sine-io/labhaus/pulls
2. 点击 "New Pull Request"
3. 选择你的分支
4. 填写 PR 描述
5. 等待 Review

## Commit Message 规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/)：

```
<type>: <description>

[optional body]

[optional footer]
```

### Type 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具配置

### 示例

```
feat: add style library API

- Implement GET /api/styles endpoint
- Add pagination and filtering
- Write unit tests

Closes #10
```

## Code Review

所有 PR 需要至少 1 人 Review 后才能合并。

Review 重点：
- ✅ 代码质量和可读性
- ✅ 测试覆盖
- ✅ 文档完整性
- ✅ 性能影响

## 测试要求

- 新功能必须包含单元测试
- 测试覆盖率不低于 80%
- 所有测试必须通过

```bash
# 运行测试
pnpm test

# 查看覆盖率
pnpm test:coverage
```

## 问题反馈

报告 Bug 或提出建议，请访问：
https://github.com/sine-io/labhaus/issues

提供以下信息：
- 问题描述
- 复现步骤
- 期望行为
- 实际行为
- 环境信息（OS、Node 版本等）

## 行为准则

- 保持友善和尊重
- 欢迎建设性的讨论
- 尊重维护者的决定

## 许可证

贡献的代码将采用 [AGPL-3.0](LICENSE) 许可证。
