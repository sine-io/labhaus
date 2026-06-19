# Infrastructure Layer

基础设施层实现了所有技术细节，包括数据库持久化、密码哈希等。

## 组件

### 1. 数据库 Repository 实现

**StyleRepository** - GORM 实现
- ✅ CRUD 操作
- ✅ 全文搜索（LIKE）
- ✅ 分类和标签过滤
- ✅ JSON 序列化 Tags
- ✅ 软删除

**UserRepository** - GORM 实现
- ✅ CRUD 操作
- ✅ 按邮箱查询
- ✅ 邮箱唯一性检查
- ✅ 软删除

**WorkflowRepository** - GORM 实现
- ✅ CRUD 操作
- ✅ 按用户查询
- ✅ 状态过滤
- ✅ JSON 序列化 Config 和 Result
- ✅ 快速状态更新
- ✅ 软删除

### 2. 密码哈希

**BcryptHasher** - Bcrypt 实现
- ✅ Hash() - 生成密码哈希
- ✅ Compare() - 验证密码
- ✅ 默认 Cost: 10

### 3. 数据库连接

**NewDB()** - PostgreSQL 连接
- ✅ GORM v2
- ✅ 连接池管理
- ✅ 日志模式配置

**AutoMigrate()** - 自动迁移
- ✅ 创建表结构
- ✅ 更新字段
- ✅ 索引管理

## 使用示例

### 初始化数据库

```go
config := persistence.DBConfig{
    Host:     "localhost",
    Port:     5432,
    User:     "postgres",
    Password: "postgres",
    DBName:   "labhaus",
    SSLMode:  "disable",
}

db, err := persistence.NewDB(config)
if err != nil {
    log.Fatal(err)
}

// 运行迁移
if err := persistence.AutoMigrate(db); err != nil {
    log.Fatal(err)
}
```

### 创建 Repository

```go
styleRepo := persistence.NewStyleRepository(db)
userRepo := persistence.NewUserRepository(db)
workflowRepo := persistence.NewWorkflowRepository(db)
```

### 创建密码哈希器

```go
hasher := persistence.NewBcryptHasher()

// 哈希密码
hash, err := hasher.Hash("mypassword")

// 验证密码
err = hasher.Compare(hash, "mypassword")
```

### 使用 Repository

```go
ctx := context.Background()

// 创建 Style
entity, _ := style.New("Anime", "desc", "prompt", "Art", []string{"anime"})
err := styleRepo.Create(ctx, entity)

// 查询 Style
found, err := styleRepo.FindByID(ctx, entity.ID)

// 搜索 Styles
results, err := styleRepo.Search(ctx, "anime", 10)
```

## 数据库模型

### StyleModel
```
id           VARCHAR(36)  PRIMARY KEY
name         VARCHAR(100) NOT NULL, INDEX
description  VARCHAR(500)
prompt       TEXT         NOT NULL
category     VARCHAR(50)  INDEX
tags         TEXT         (JSON array)
created_at   TIMESTAMP    NOT NULL
updated_at   TIMESTAMP    NOT NULL
deleted_at   TIMESTAMP    INDEX (soft delete)
```

### UserModel
```
id            VARCHAR(36)  PRIMARY KEY
email         VARCHAR(255) NOT NULL, UNIQUE INDEX
password_hash VARCHAR(255) NOT NULL
name          VARCHAR(100) NOT NULL
role          VARCHAR(20)  NOT NULL, DEFAULT 'user'
created_at    TIMESTAMP    NOT NULL
updated_at    TIMESTAMP    NOT NULL
deleted_at    TIMESTAMP    INDEX (soft delete)
```

### WorkflowModel
```
id         VARCHAR(36)  PRIMARY KEY
user_id    VARCHAR(36)  NOT NULL, INDEX
style_id   VARCHAR(36)  NOT NULL, INDEX
state      VARCHAR(20)  NOT NULL, INDEX
config     JSONB        NOT NULL
result     JSONB        NULL
created_at TIMESTAMP    NOT NULL
updated_at TIMESTAMP    NOT NULL
deleted_at TIMESTAMP    INDEX (soft delete)
```

## 集成测试

集成测试需要 PostgreSQL 数据库：

```bash
# 启动 PostgreSQL (Docker)
docker run -d \
  --name postgres-test \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=labhaus_test \
  -p 5432:5432 \
  postgres:14

# 运行集成测试
go test ./tests/integration/... -v
```

## 依赖

- `gorm.io/gorm` - ORM 框架
- `gorm.io/driver/postgres` - PostgreSQL 驱动
- `github.com/google/uuid` - UUID 生成
- `golang.org/x/crypto/bcrypt` - 密码哈希

## 架构优势

✅ **依赖倒置（DIP）**: 实现领域层定义的接口  
✅ **技术隔离**: 领域层不知道 GORM 存在  
✅ **易于替换**: 可切换到其他数据库或 ORM  
✅ **测试友好**: 集成测试与单元测试分离  
✅ **软删除**: 保留历史数据，支持恢复  
