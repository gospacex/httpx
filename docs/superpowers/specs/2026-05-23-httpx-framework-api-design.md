# httpx 框架级 API 设计

> **面向 AI 代理的工作者：** 必需子技能：使用 superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 逐任务实现此计划。步骤使用复选框（`- [ ]`）语法来跟踪进度。

**目标：** 设计一套高表现力的框架级 API，让用户通过 `httpx.new()` + `app.run("./config.yaml")` 即可启动服务，无需关心底层 adapter。

**架构：** 核心层（App/Router/Route）+ Adapter 层（gin/hertz/nethttp）+ 配置层。App 继承 Router 的所有方法，Adapter 根据 config.yaml 自动选择，Graceful Shutdown 默认实现。

**技术栈：** Go, httpx 核心类型（HandlerFunc, MiddlewareFunc, HandlerContext）, viper（配置文件解析）

---

## 文件结构

```
httpx/
├── app.go              # App 类型定义
├── router.go           # 现有 router.go 扩展
├── types.go            # 现有 types.go 扩展
├── config.go           # 配置解析
├── adapter/
│   ├── gin/
│   ├── hertz/
│   └── nethttp/
└── adapter_registry.go # Adapter 注册表

examples/
├── main.go
├── route.go
└── config.yaml
```

---

## 核心类型设计

### App

```go
type App struct {
    router     *routerImpl
    middlewares []MiddlewareFunc
    adapterName string
    configPath  string
}

func new() *App {
    return &App{
        router: newRouter(),
    }
}

func (a *App) GET(path string, handlers ...HandlerFunc) Route {
    return a.router.GET(path, handlers...)
}
func (a *App) POST(path string, handlers ...HandlerFunc) Route { /* 同上 */ }
func (a *App) PUT(path string, handlers ...HandlerFunc) Route { /* 同上 */ }
func (a *App) DELETE(path string, handlers ...HandlerFunc) Route { /* 同上 */ }
func (a *App) PATCH(path string, handlers ...HandlerFunc) Route { /* 同上 */ }
func (a *App) WS(path string, handlers ...HandlerFunc) Route { /* 同上 */ }
func (a *App) GROUP(prefix string, mw ...MiddlewareFunc) *RouterGroup { /* 同上 */ }

func (a *App) use(mw ...MiddlewareFunc) *App {
    a.middlewares = append(a.middlewares, mw...)
    return a
}

func (a *App) run(configPath string) error {
    a.configPath = configPath
    // 解析配置，选择 adapter，启动服务
}
```

### Route

```go
type Route struct {
    path     string
    handlers []HandlerFunc
    middlewares []MiddlewareFunc
    router     *routerImpl
}

func (r *Route) use(mw ...MiddlewareFunc) *Route {
    r.middlewares = append(r.middlewares, mw...)
    return r
}
```

### RouterGroup

```go
type RouterGroup struct {
    Prefix      string
    middlewares []MiddlewareFunc
    router      *routerImpl
}

func (g *RouterGroup) GET(path string, handlers ...HandlerFunc) *Route {
    return g.router.GET(g.Prefix+path, handlers...)
}
func (g *RouterGroup) POST(path string, handlers ...HandlerFunc) *Route { /* ... */ }
func (g *RouterGroup) PUT(path string, handlers ...HandlerFunc) *Route { /* ... */ }
func (g *RouterGroup) DELETE(path string, handlers ...HandlerFunc) *Route { /* ... */ }
func (g *RouterGroup) PATCH(path string, handlers ...HandlerFunc) *Route { /* ... */ }
```

---

## 配置解析

### config.yaml

```yaml
adapter: gin
port: 8080
```

### config.go

```go
type Config struct {
    Adapter string `yaml:"adapter"`
    Port    int    `yaml:"port"`
}

func LoadConfig(path string) (*Config, error) {
    viper.SetConfigFile(path)
    viper.SetConfigType("yaml")
    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    return &cfg, nil
}
```

---

## Adapter 注册机制

### 内置 + 可扩展

httpx 内置 gin、hertz、nethttp，用户无需 import。配置选择哪个就用哪个。自定义 adapter 可通过 `RegisterAdapter` 注册。

```go
var adapterRegistry = map[string]AdapterFactory{}

func RegisterAdapter(name string, factory AdapterFactory) {
    adapterRegistry[name] = factory
}

func getAdapter(name string) (AdapterFactory, error) {
    factory, ok := adapterRegistry[name]
    if !ok {
        return nil, fmt.Errorf("adapter %q not found", name)
    }
    return factory, nil
}

// 内置注册
func init() {
    RegisterAdapter("gin", gin.NewServer)
    RegisterAdapter("hertz", hertz.NewServer)
    RegisterAdapter("nethttp", nethttp.NewServer)
}
```

---

## 实现计划

### 任务 1：创建 app.go

**文件：**
- 创建：`httpx/app.go`

- [ ] **步骤 1：定义 App 结构体**

```go
type App struct {
    router     *routerImpl
    middlewares []MiddlewareFunc
    adapterName string
    configPath  string
}

func new() *App {
    return &App{
        router: newRouter(),
    }
}
```

- [ ] **步骤 2：实现 App.Use() 方法**

```go
func (a *App) Use(mw ...MiddlewareFunc) *App {
    a.middlewares = append(a.middlewares, mw...)
    return a
}
```

- [ ] **步骤 3：实现 App.Run() 方法**

```go
func (a *App) Run(configPath string) error {
    cfg, err := LoadConfig(configPath)
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }

    factory, err := getAdapter(cfg.Adapter)
    if err != nil {
        return err
    }

    srv := factory()
    for _, mw := range a.middlewares {
        srv.Use(mw)
    }

    router := srv.Router()
    a.router.setupToRouter(router)

    return srv.StartWithGraceful()
}
```

- [ ] **步骤 4：Commit**

---

### 任务 2：创建 router.go 扩展

**文件：**
- 修改：`httpx/router.go`
- 创建：`httpx/route.go`

- [ ] **步骤 1：添加 routerImpl 到 Router 的转换**

routerImpl 需要能够将注册路由转换到具体 adapter 的 Router 上。

- [ ] **步骤 2：定义 Route 类型**

```go
type Route struct {
    path       string
    method     string
    handlers   []HandlerFunc
    middlewares []MiddlewareFunc
    router     *routerImpl
}

func (r *Route) Use(mw ...MiddlewareFunc) *Route {
    r.middlewares = append(r.middlewares, mw...)
    return r
}
```

- [ ] **步骤 3：实现 Route.Use() 链式调用**

- [ ] **步骤 4：Commit**

---

### 任务 3：创建 config.go

**文件：**
- 创建：`httpx/config.go`

- [ ] **步骤 1：使用 viper 解析 YAML 配置**

- [ ] **步骤 2：Commit**

---

### 任务 4：创建 adapter_registry.go

**文件：**
- 创建：`httpx/adapter_registry.go`

- [ ] **步骤 1：定义 AdapterFactory 接口**

```go
type AdapterFactory func() Server
```

- [ ] **步骤 2：实现注册表和内置注册**

- [ ] **步骤 3：Commit**

---

### 任务 5：适配现有 adapter

**文件：**
- 修改：`httpx/adapter/gin/server.go`
- 修改：`httpx/adapter/hertz/server.go`
- 修改：`httpx/adapter/nethttp/server.go`

- [ ] **步骤 1：确认各 adapter 实现 httpx.Server 接口**

- [ ] **步骤 2：在 adapter_registry.go 中注册**

- [ ] **步骤 3：Commit**

---

### 任务 6：创建示例

**文件：**
- 创建：`examples/demo/main.go`
- 创建：`examples/demo/route.go`
- 创建：`examples/demo/config.yaml`
- 创建：`examples/demo/config_hertz.yaml`

```go
// main.go
func main() {
    app := httpx.new()
    setupRoutes(app)
    app.Run("./config.yaml")
}
```

```go
// route.go
func setupRoutes(app *httpx.App) {
    app.Use(recoveryMiddleware, loggerMiddleware)
    app.GET("/health", healthHandler)
    app.GET("/users", listUsers).POST("/users", createUser)
    app.group("/articles", authMiddleware).GET("/", listArticles)
}
```

```yaml
# config.yaml (gin)
adapter: gin
port: 8080
```

```yaml
# config_hertz.yaml (hertz)
adapter: hertz
port: 8080
```

- [ ] **步骤 1：实现示例代码（gin 配置）**

- [ ] **步骤 2：实现示例代码（hertz 配置）**

- [ ] **步骤 3：运行 gin 测试**

- [ ] **步骤 4：运行 hertz 测试**

- [ ] **步骤 5：Commit**

---

### 任务 7：全量编译验证

- [ ] **步骤 1：运行 `go build ./...` 确认无编译错误**

- [ ] **步骤 2：运行 `go vet ./...` 确认无警告**

- [ ] **步骤 3：运行 `go test ./...` 确认测试通过**

- [ ] **步骤 4：Commit**

---

## 测试覆盖要求

测试必须覆盖 **gin** 和 **hertz** 两种 adapter 场景：

1. **App 创建和路由注册** — 使用 gin 和 hertz 分别创建 App，注册相同路由，验证行为一致
2. **全局中间件** — 分别在 gin 和 hertz 上测试 `app.Use()` 效果
3. **单路由中间件** — 分别在 gin 和 hertz 上测试 `app.GET().Use()` 链式调用
4. **分组路由** — 分别在 gin 和 hertz 上测试 `app.Group()` 创建分组
5. **启动和停止** — 分别在 gin 和 hertz 上测试 `app.Run()` 和 graceful shutdown

测试代码示例：
```go
func TestGinApp(t *testing.T) {
    app := httpx.new()
    app.GET("/health", healthHandler)
    // 验证 gin adapter 行为
}

func TestHertzApp(t *testing.T) {
    app := httpx.new()
    app.GET("/health", healthHandler)
    // 验证 hertz adapter 行为
}
```

---

## 使用示例

### main.go

```go
package main

import (
    "context"
    "github.com/gospacex/httpx"
)

func main() {
    app := httpx.new()
    setupRoutes(app)
    if err := app.Run("./config.yaml"); err != nil {
        panic(err)
    }
}
```

### route.go

```go
func setupRoutes(app *httpx.App) {
    // 全局中间件
    app.Use(recoveryMiddleware, loggerMiddleware)

    // 简单路由直接注册
    app.GET("/health", func(ctx context.Context, hc httpx.HandlerContext) {
        hc.AbortJSON(200, map[string]string{"status": "ok"})
    })

    // 单路由中间件
    app.GET("/admin", adminHandler).Use(adminMiddleware)

    // 分组路由
    usersGroup := app.Group("/users", authMiddleware)
    usersGroup.GET("/", listUsersHandler)
    usersGroup.POST("/", createUserHandler)
    usersGroup.GET("/:id", getUserHandler)
    usersGroup.PUT("/:id", updateUserHandler)
    usersGroup.DELETE("/:id", deleteUserHandler)

    // 嵌套分组
    articlesGroup := app.Group("/articles", authMiddleware)
    articlesGroup.GET("/", listArticlesHandler)
    articlesGroup.POST("/", createArticleHandler)
}
```

### config.yaml

```yaml
adapter: gin
port: 8080
```