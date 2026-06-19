# Go 后端架构设计文档

## 架构理念

Labhaus Go 后端采用 **DDD Lite + CQRS + Clean Architecture + DIP + TDD** 的混合架构模式。

### 核心原则

1. **DDD Lite** (领域驱动设计-轻量版)
   - 聚焦核心领域模型
   - 避免过度工程化
   - 保持 Go 的简洁性

2. **CQRS** (命令查询职责分离)
   - 读写分离
   - 查询优化（缓存、只读副本）
   - 命令验证和事件溯源

3. **Clean Architecture** (整洁架构)
   - 依赖倒置：内层不依赖外层
   - 业务逻辑独立于框架
   - 可测试性优先

4. **DIP** (依赖倒置原则)
   - 高层模块不依赖低层模块
   - 都依赖于抽象（接口）
   - 接口由消费者定义

5. **TDD** (测试驱动开发)
   - 先写测试，后写实现
   - Red-Green-Refactor 循环
   - 保持高测试覆盖率

## 项目结构

```
backend/
├── cmd/
│   └── api/
│       └── main.go                    # 应用入口
├── internal/
│   ├── domain/                        # 领域层 (最内层)
│   │   ├── style/
│   │   │   ├── style.go              # 领域实体
│   │   │   ├── repository.go         # 仓储接口（由领域定义）
│   │   │   └── service.go            # 领域服务
│   │   ├── user/
│   │   │   ├── user.go
│   │   │   ├── repository.go
│   │   │   └── auth_service.go
│   │   └── workflow/
│   │       ├── workflow.go
│   │       └── executor.go
│   ├── application/                   # 应用层
│   │   ├── command/                  # CQRS - 命令
│   │   │   ├── register_user.go
│   │   │   ├── create_style.go
│   │   │   └── handler.go
│   │   ├── query/                    # CQRS - 查询
│   │   │   ├── get_styles.go
│   │   │   ├── recommend_styles.go
│   │   │   └── handler.go
│   │   └── dto/                      # 数据传输对象
│   │       ├── style_dto.go
│   │       └── user_dto.go
│   ├── infrastructure/                # 基础设施层 (最外层)
│   │   ├── persistence/              # 数据持久化
│   │   │   ├── postgres/
│   │   │   │   ├── style_repository.go    # 实现 domain.StyleRepository
│   │   │   │   ├── user_repository.go
│   │   │   │   └── gorm.go
│   │   │   └── redis/
│   │   │       └── cache.go
│   │   ├── http/                     # HTTP 适配器
│   │   │   ├── server.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── logging.go
│   │   │   │   └── recovery.go
│   │   │   └── handlers/
│   │   │       ├── style_handler.go
│   │   │       ├── auth_handler.go
│   │   │       └── health_handler.go
│   │   ├── queue/                    # 任务队列
│   │   │   └── asynq.go
│   │   └── external/                 # 外部服务
│   │       ├── openai/
│   │       └── minio/
│   └── pkg/                          # 共享工具包
│       ├── errors/
│       │   └── errors.go
│       ├── validator/
│       │   └── validator.go
│       ├── jwt/
│       │   └── jwt.go
│       └── logger/
│           └── logger.go
├── migrations/                        # 数据库迁移
│   ├── 001_create_users.sql
│   └── 002_create_styles.sql
├── tests/                            # 测试
│   ├── unit/                         # 单元测试
│   ├── integration/                  # 集成测试
│   └── e2e/                          # 端到端测试
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

## 分层详解

### 1. Domain Layer (领域层) - 核心

**职责**：业务规则和领域逻辑

**原则**：
- 不依赖任何外层
- 只包含纯业务逻辑
- 定义接口，不实现基础设施

**示例**：
```go
// internal/domain/style/style.go
package style

type Style struct {
    ID       string
    Title    string
    Prompt   string
    Category string
    Featured bool
}

// 领域方法
func (s *Style) MarkAsFeatured() {
    s.Featured = true
}

func (s *Style) Validate() error {
    if s.Title == "" {
        return errors.New("title is required")
    }
    return nil
}

// internal/domain/style/repository.go
// 接口由领域层定义（DIP）
package style

type Repository interface {
    Save(ctx context.Context, style *Style) error
    FindByID(ctx context.Context, id string) (*Style, error)
    FindAll(ctx context.Context, filter Filter) ([]*Style, error)
}
```

### 2. Application Layer (应用层) - 用例

**职责**：协调领域对象完成用例

**原则**：
- 不包含业务规则
- 编排领域对象
- CQRS 分离

**Command 示例** (写操作)：
```go
// internal/application/command/register_user.go
package command

type RegisterUserCommand struct {
    Email    string
    Password string
    Name     string
}

type RegisterUserHandler struct {
    userRepo user.Repository
}

func (h *RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUserCommand) error {
    // 1. 验证
    if err := validate(cmd); err != nil {
        return err
    }
    
    // 2. 创建领域对象
    user := user.NewUser(cmd.Email, cmd.Password, cmd.Name)
    
    // 3. 持久化
    return h.userRepo.Save(ctx, user)
}
```

**Query 示例** (读操作)：
```go
// internal/application/query/get_styles.go
package query

type GetStylesQuery struct {
    Category string
    Page     int
    Limit    int
}

type GetStylesHandler struct {
    styleRepo style.Repository
    cache     Cache
}

func (h *GetStylesHandler) Handle(ctx context.Context, q GetStylesQuery) ([]*dto.StyleDTO, error) {
    // 查询可以直接访问缓存、读副本等优化
    styles, err := h.styleRepo.FindAll(ctx, style.Filter{
        Category: q.Category,
        Page:     q.Page,
        Limit:    q.Limit,
    })
    
    return toDTO(styles), nil
}
```

### 3. Infrastructure Layer (基础设施层) - 实现

**职责**：实现接口，对接外部系统

**原则**：
- 实现领域层定义的接口
- 包含框架、数据库、HTTP 等
- 可替换

**Repository 实现**：
```go
// internal/infrastructure/persistence/postgres/style_repository.go
package postgres

type StyleRepository struct {
    db *gorm.DB
}

// 实现 domain.StyleRepository 接口
func (r *StyleRepository) Save(ctx context.Context, style *style.Style) error {
    model := toModel(style)
    return r.db.WithContext(ctx).Create(model).Error
}

func (r *StyleRepository) FindByID(ctx context.Context, id string) (*style.Style, error) {
    var model StyleModel
    err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return toDomain(&model), nil
}
```

**HTTP Handler**：
```go
// internal/infrastructure/http/handlers/style_handler.go
package handlers

type StyleHandler struct {
    getStylesQuery *query.GetStylesHandler
    createStyleCmd *command.CreateStyleHandler
}

func (h *StyleHandler) GetStyles(c *gin.Context) {
    var req GetStylesRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 调用应用层
    result, err := h.getStylesQuery.Handle(c.Request.Context(), query.GetStylesQuery{
        Category: req.Category,
        Page:     req.Page,
        Limit:    req.Limit,
    })
    
    c.JSON(200, result)
}
```

## CQRS 实践

### 命令 (Command) - 写操作

```go
// 特点：
// - 修改状态
// - 返回错误或成功
// - 可能触发事件
// - 需要事务保证

type CreateStyleCommand struct {
    Title    string
    Prompt   string
    Category string
}

func (h *CreateStyleHandler) Handle(ctx context.Context, cmd CreateStyleCommand) error {
    style := style.New(cmd.Title, cmd.Prompt, cmd.Category)
    
    if err := style.Validate(); err != nil {
        return err
    }
    
    return h.styleRepo.Save(ctx, style)
}
```

### 查询 (Query) - 读操作

```go
// 特点：
// - 不修改状态
// - 返回数据
// - 可以使用缓存
// - 可以直接查询读模型

type GetStylesQuery struct {
    Category string
    Search   string
    Page     int
    Limit    int
}

func (h *GetStylesHandler) Handle(ctx context.Context, q GetStylesQuery) ([]*dto.StyleDTO, error) {
    // 优先从缓存读取
    if cached := h.cache.Get(q); cached != nil {
        return cached, nil
    }
    
    // 查询数据库
    styles, err := h.styleRepo.FindAll(ctx, toFilter(q))
    if err != nil {
        return nil, err
    }
    
    result := toDTO(styles)
    h.cache.Set(q, result)
    
    return result, nil
}
```

## 依赖注入 (DI)

使用构造函数注入，保持简单：

```go
// cmd/api/main.go
func main() {
    // 1. 初始化基础设施
    db := initDB()
    cache := initRedis()
    logger := initLogger()
    
    // 2. 初始化仓储（实现接口）
    styleRepo := postgres.NewStyleRepository(db)
    userRepo := postgres.NewUserRepository(db)
    
    // 3. 初始化应用服务（注入依赖）
    getStylesQuery := query.NewGetStylesHandler(styleRepo, cache)
    createStyleCmd := command.NewCreateStyleHandler(styleRepo)
    registerUserCmd := command.NewRegisterUserHandler(userRepo)
    
    // 4. 初始化 HTTP handlers
    styleHandler := handlers.NewStyleHandler(getStylesQuery, createStyleCmd)
    authHandler := handlers.NewAuthHandler(registerUserCmd)
    
    // 5. 启动服务器
    router := gin.Default()
    router.GET("/api/styles", styleHandler.GetStyles)
    router.POST("/api/styles", styleHandler.CreateStyle)
    router.Run(":8080")
}
```

## TDD 工作流

### Red-Green-Refactor 循环

1. **Red** - 写一个失败的测试
2. **Green** - 写最简单的实现让测试通过
3. **Refactor** - 重构代码，保持测试通过

### 测试金字塔

```
        /\
       /E2E\         ← 少量（慢、集成度高）
      /------\
     /Integration\   ← 适量（中速、部分集成）
    /------------\
   /    Unit      \  ← 大量（快、隔离）
  /----------------\
```

### 示例：TDD 开发新功能

```go
// 1. RED - 写测试（失败）
// internal/domain/style/style_test.go
func TestStyle_MarkAsFeatured(t *testing.T) {
    style := style.New("Test", "Prompt", "Category")
    
    style.MarkAsFeatured()
    
    assert.True(t, style.Featured)
}

// 2. GREEN - 最简实现（通过）
func (s *Style) MarkAsFeatured() {
    s.Featured = true
}

// 3. REFACTOR - 优化（如需要）
func (s *Style) MarkAsFeatured() error {
    if s.Featured {
        return errors.New("already featured")
    }
    s.Featured = true
    return nil
}
```

### 测试策略

**Unit Tests** (单元测试)：
```go
// 测试领域逻辑
func TestStyleValidation(t *testing.T) {
    tests := []struct {
        name    string
        style   *Style
        wantErr bool
    }{
        {"valid", NewStyle("Title", "Prompt", "Cat"), false},
        {"empty title", NewStyle("", "Prompt", "Cat"), true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.style.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("want error: %v, got: %v", tt.wantErr, err)
            }
        })
    }
}
```

**Integration Tests** (集成测试)：
```go
// 测试仓储实现
func TestStyleRepository_Save(t *testing.T) {
    // 使用测试数据库
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := postgres.NewStyleRepository(db)
    style := style.New("Test", "Prompt", "Category")
    
    err := repo.Save(context.Background(), style)
    assert.NoError(t, err)
    
    // 验证保存成功
    found, err := repo.FindByID(context.Background(), style.ID)
    assert.NoError(t, err)
    assert.Equal(t, style.Title, found.Title)
}
```

**E2E Tests** (端到端测试)：
```go
// 测试完整 HTTP 流程
func TestCreateStyle_E2E(t *testing.T) {
    server := setupTestServer(t)
    defer server.Close()
    
    resp, err := http.Post(server.URL+"/api/styles", "application/json", 
        strings.NewReader(`{"title":"Test","prompt":"Prompt","category":"Cat"}`))
    
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
}
```

## 最佳实践

### 1. 保持简单（Go Way）

```go
// ❌ 过度抽象
type StyleService interface {
    CreateStyle(StyleDTO) error
    UpdateStyle(StyleDTO) error
    DeleteStyle(string) error
}

// ✅ 简洁实用
type StyleRepository interface {
    Save(context.Context, *Style) error
    FindByID(context.Context, string) (*Style, error)
}
```

### 2. 接口隔离

```go
// ❌ 胖接口
type Repository interface {
    Save(...) error
    Update(...) error
    Delete(...) error
    FindByID(...) error
    FindAll(...) error
    Count(...) int
}

// ✅ 小接口组合
type Saver interface {
    Save(context.Context, *Style) error
}

type Finder interface {
    FindByID(context.Context, string) (*Style, error)
}

type Repository interface {
    Saver
    Finder
}
```

### 3. 表驱动测试

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "invalid", true},
        {"empty", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试逻辑
        })
    }
}
```

### 4. 错误处理

```go
// internal/pkg/errors/errors.go
package errors

type DomainError struct {
    Code    string
    Message string
    Err     error
}

func (e *DomainError) Error() string {
    return e.Message
}

// 预定义错误
var (
    ErrNotFound = &DomainError{Code: "NOT_FOUND", Message: "resource not found"}
    ErrInvalidInput = &DomainError{Code: "INVALID_INPUT", Message: "invalid input"}
)
```

## 参考资料

- [Three Dots Labs - DDD Lite in Go](https://threedots.tech/post/ddd-lite-in-go-introduction/)
- [Three Dots Labs - Basic CQRS in Go](https://threedots.tech/post/basic-cqrs-in-go/)
- [Clean Architecture in Golang](https://pkritiotis.io/clean-architecture-in-golang/)
- [Dependency Inversion in Go](https://www.ompluscator.com/article/golang/practical-solid-dependency-inversion/)
- [TDD in Go](https://threedots.tech/post/introducing-clean-architecture/)
