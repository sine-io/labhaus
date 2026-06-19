# 开发指南

## 环境要求

- Node.js >= 20.0.0
- pnpm >= 9.0.0
- Docker & Docker Compose
- Git

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/sine-io/labhaus.git
cd labhaus
```

### 2. 安装依赖

```bash
pnpm install
```

### 3. 启动开发环境

```bash
# 启动 PostgreSQL、Redis、MinIO
docker compose up -d

# 等待服务就绪
docker compose ps

# 复制环境变量
cp .env.example .env.local
```

### 4. 开发模式

```bash
# 启动所有应用的开发服务器
pnpm dev

# 或单独启动某个应用
pnpm --filter @labhaus/api dev
```

## 项目结构

```
labhaus/
├── apps/                    # 应用
│   ├── api/                 # 后端 API 服务
│   └── web/                 # 前端 Web 应用
├── packages/                # 共享包
│   ├── types/               # TypeScript 类型定义
│   ├── config/              # 共享配置
│   └── utils/               # 工具函数
├── docs/                    # 文档
├── .github/                 # GitHub Actions
├── docker-compose.yml       # Docker 配置
├── turbo.json               # Turborepo 配置
├── pnpm-workspace.yaml      # pnpm workspace 配置
└── package.json             # 根 package.json
```

## 常用命令

```bash
# 安装依赖
pnpm install

# 开发模式
pnpm dev

# 构建所有包
pnpm build

# 代码检查
pnpm lint

# 格式化代码
pnpm format

# 运行测试
pnpm test

# 清理构建产物
pnpm clean
```

## Docker 服务

### PostgreSQL
- 端口: 5432
- 数据库: labhaus
- 用户名: labhaus
- 密码: labhaus_dev_password

### Redis
- 端口: 6379

### MinIO
- API 端口: 9000
- Console 端口: 9001
- 用户名: minioadmin
- 密码: minioadmin

访问 MinIO Console: http://localhost:9001

## 开发工作流

1. 从 `main` 分支创建功能分支
2. 开发并提交代码
3. 运行 `pnpm lint` 和 `pnpm test` 确保通过
4. 提交 Pull Request
5. 等待 CI 通过和 Code Review
6. 合并到 `main`

## 代码规范

- 使用 ESLint 进行代码检查
- 使用 Prettier 进行代码格式化
- 提交前运行 `pnpm lint` 和 `pnpm format`
- 遵循 TypeScript 严格模式

## 故障排查

### Docker 权限问题

如果遇到 "permission denied" 错误：

```bash
# 将当前用户添加到 docker 组
sudo usermod -aG docker $USER

# 重新登录或运行
newgrp docker
```

### pnpm 安装失败

```bash
# 清理缓存
pnpm store prune

# 删除 node_modules 重新安装
rm -rf node_modules
pnpm install
```

## 贡献

参考 [CONTRIBUTING.md](CONTRIBUTING.md)
