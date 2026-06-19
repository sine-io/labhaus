# Labhaus API

后端 API 服务，提供样式库查询、工作流管理等功能。

## 技术栈

- **框架**: Hono (轻量级 Web 框架)
- **数据库**: PostgreSQL
- **运行时**: Node.js + tsx
- **类型系统**: TypeScript + Zod

## 快速开始

### 1. 启动数据库

```bash
# 在项目根目录
docker compose up -d postgres
```

### 2. 运行数据库迁移

```bash
cd apps/api
pnpm migrate
```

### 3. 导入样式库数据

从 awesome-gpt-image-2 项目导入 500+ 样式案例：

```bash
pnpm import-styles
```

### 4. 启动开发服务器

```bash
pnpm dev
```

API 将运行在 http://localhost:3001

## API 端点

### GET /api/styles

查询样式库，支持分类、风格、场景筛选和全文搜索。

**查询参数：**

- `category` - 分类筛选（如 "UI & Interfaces"）
- `style` - 风格筛选（如 "Realistic"）
- `scene` - 场景筛选（如 "Tech"）
- `featured` - 是否精选（true/false）
- `search` - 全文搜索关键词
- `page` - 页码（默认 1）
- `limit` - 每页数量（默认 20，最大 100）

**示例：**

```bash
# 获取所有样式
curl http://localhost:3001/api/styles

# 筛选 UI 类样式
curl http://localhost:3001/api/styles?category=UI%20%26%20Interfaces

# 全文搜索
curl http://localhost:3001/api/styles?search=portrait

# 组合筛选
curl "http://localhost:3001/api/styles?style=Realistic&scene=Tech&page=1&limit=10"
```

**响应：**

```json
{
  "styles": [
    {
      "id": "uuid",
      "case_id": 505,
      "title": "夜间手机光沙发肖像",
      "prompt": "A young adult woman...",
      "prompt_preview": "A young adult woman...",
      "category": "Photography & Realism",
      "styles": ["Realistic"],
      "scenes": ["Tech", "Commerce", "Social"],
      "image_url": "/images/case505.jpg",
      "source_label": "@iamaiistudio",
      "source_url": "https://x.com/...",
      "github_url": "https://github.com/...",
      "featured": false,
      "created_at": "2026-06-19T...",
      "updated_at": "2026-06-19T..."
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 503,
    "totalPages": 26
  }
}
```

### GET /api/styles/:id

获取单个样式详情。

**示例：**

```bash
curl http://localhost:3001/api/styles/{uuid}
```

### GET /health

健康检查端点。

```bash
curl http://localhost:3001/health
```

## 数据模型

### Style

```typescript
{
  id: string;              // UUID
  case_id: number;         // 原始案例 ID
  title: string;           // 标题
  prompt: string;          // 完整提示词
  prompt_preview: string;  // 提示词预览
  category: string;        // 分类
  styles: string[];        // 风格标签
  scenes: string[];        // 场景标签
  image_url: string;       // 预览图 URL
  source_label: string;    // 来源标签
  source_url: string;      // 来源 URL
  github_url: string;      // GitHub URL
  featured: boolean;       // 是否精选
  created_at: string;      // 创建时间
  updated_at: string;      // 更新时间
}
```

## 开发

```bash
# 类型检查
pnpm typecheck

# 运行测试
pnpm test

# 构建
pnpm build

# 生产启动
pnpm start
```

## 环境变量

在项目根目录的 `.env.local` 中配置：

```bash
DATABASE_URL=postgresql://labhaus:labhaus_dev_password@localhost:5432/labhaus
API_PORT=3001
NODE_ENV=development
```
