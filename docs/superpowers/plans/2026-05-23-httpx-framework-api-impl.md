# httpx 框架级 API 实现计划

> **面向 AI 代理的工作者：** 必需子技能：使用 superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 逐任务实现此计划。步骤使用复选框（`- [ ]`）语法来跟踪进度。

**目标：** 实现一套高表现力的框架级 API，让用户通过 `httpx.new()` + `app.Run("./config.yaml")` 即可启动服务，无需关心底层 adapter。

**架构：** 核心层（App/Router/Route）+ Adapter 层（gin/hertz/nethttp）+ 配置层。App 继承 Router 的所有方法，Adapter 根据 config.yaml 自动选择，Graceful Shutdown 默认实现。

**技术栈：** Go, viper（配置文件解析）

---

## 文件结构

```
httpx/
├── app.go              # App 类型定义（新增）
├── router.go           # Router/RouterGroup 定义（修改现有）
├── route.go            # Route 类型（新增）
├── config.go           # 配置解析（新增）
├── adapter_registry.go # Adapter 注册表（新增）
├── adapter/
│   ├── gin/
│   │   └── server.go   # GinServer（修改：实现 AdapterFactory）
│   ├── hertz/
│   │   └── server.go   # HertzServer（修改：实现 AdapterFactory）
│   └── nethttp/
│       └── server.go   # NetHTTPServer（修改：实现 AdapterFactory）

examples/demo/
├── main.go             # 入口（新增）
├── route.go            # 路由注册（新增）
└── config.yaml         # 配置文件（新增）
```

---

## 任务 1：创建 config.go（配置解析）

**文件：**
- 创建：`httpx/config.go`

- [ ] **步骤 1：编写 config.go**

```go
package httpx

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Adapter string `mapstructure:"adapter"`
	Port    int    `mapstructure:"port"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.Adapter == "" {
		return nil, fmt.Errorf("adapter is required in config")
	}
	if cfg.Port <= 0 {
		cfg.Port = 8080
	}

	return &cfg, nil
}
```

- [ ] **步骤 2：运行测试确认无编译错误**

运行：`go build ./httpx/`
预期：PASS

- [ ] **步骤 3：Commit**

```bash
git add httpx/config.go
git commit -m "feat: add config loading with viper"
```

---

## 任务 2：创建 adapter_registry.go（Adapter 注册表）

**文件：**
- 创建：`httpx/adapter_registry.go`

- [ ] **步骤 1：编写 adapter_registry.go**

```go
package httpx

import "fmt"

type AdapterFactory func() Server

var adapterRegistry = make(map[string]AdapterFactory)

func RegisterAdapter(name string, factory AdapterFactory) {
	adapterRegistry[name] = factory
}

func getAdapter(name string) (AdapterFactory, error) {
	factory, ok := adapterRegistry[name]
	if !ok {
		return nil, fmt.Errorf("adapter %q not found, available adapters: gin, hertz, nethttp", name)
	}
	return factory, nil
}
```

- [ ] **步骤 2：编写 adapter_registry_test.go**

```go
package httpx

import (
	"testing"
)

func TestGetAdapter_Gin(t *testing.T) {
	RegisterAdapter("gin", func() Server { return nil })
	factory, err := getAdapter("gin")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if factory == nil {
		t.Fatal("expected factory, got nil")
	}
}

func TestGetAdapter_NotFound(t *testing.T) {
	_, err := getAdapter("nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
```

- [ ] **步骤 3：运行测试验证失败**

运行：`go test ./httpx/ -run TestGetAdapter -v`
预期：PASS

- [ ] **步骤 4：Commit**

```bash
git add httpx/adapter_registry.go httpx/adapter_registry_test.go
git commit -m "feat: add adapter registry with RegisterAdapter/getAdapter"
```

---

## 任务 3：修改现有 router.go（扩展 RouterGroup）

**文件：**
- 修改：`httpx/router.go`

现有 router.go 内容：
```go
package httpx

type Router interface {
	GET(path string, handlers ...HandlerFunc) Router
	POST(path string, handlers ...HandlerFunc) Router
	PUT(path string, handlers ...HandlerFunc) Router
	DELETE(path string, handlers ...HandlerFunc) Router
	PATCH(path string, handlers ...HandlerFunc) Router
	GROUP(prefix string, mw ...MiddlewareFunc) *RouterGroup
	WS(path string, handlers ...HandlerFunc) Router
}

type RouterGroup struct {
	Router      Router
	Prefix      string
	middlewares []MiddlewareFunc
}
```

- [ ] **步骤 1：添加 routerImpl 内部实现**

在 router.go 末尾添加：

```go
type routerImpl struct {
	routes    []routeRecord
	groups    []*RouterGroup
}

type routeRecord struct {
	method     string
	path       string
	handlers   []HandlerFunc
	middlewares []MiddlewareFunc
}

func newRouter() *routerImpl {
	return &routerImpl{}
}

func (r *routerImpl) GET(path string, handlers ...HandlerFunc) *Route {
	return r.addRoute("GET", path, handlers)
}

func (r *routerImpl) POST(path string, handlers ...HandlerFunc) *Route {
	return r.addRoute("POST", path, handlers)
}

func (r *routerImpl) PUT(path string, handlers ...HandlerFunc) *Route {
	return r.addRoute("PUT", path, handlers)
}

func (r *routerImpl) DELETE(path string, handlers ...HandlerFunc) *Route {
	return r.addRoute("DELETE", path, handlers)
}

func (r *routerImpl) PATCH(path string, handlers ...HandlerFunc) *Route {
	return r.addRoute("PATCH", path, handlers)
}

func (r *routerImpl) WS(path string, handlers ...HandlerFunc) *Route {
	return r.addRoute("WS", path, handlers)
}

func (r *routerImpl) GROUP(prefix string, mw ...MiddlewareFunc) *RouterGroup {
	g := &RouterGroup{
		Prefix:      prefix,
		Router:      &groupRouter{impl: r, prefix: prefix, middlewares: mw},
		middlewares: mw,
	}
	r.groups = append(r.groups, g)
	return g
}

func (r *routerImpl) addRoute(method, path string, handlers []HandlerFunc) *Route {
	r.routes = append(r.routes, routeRecord{
		method:     method,
		path:       path,
		handlers:   handlers,
		middlewares: nil,
	})
	return &Route{
		path:       path,
		method:     method,
		handlers:   handlers,
		middlewares: nil,
		routerImpl: r,
	}
}

func (r *routerImpl) setupToRouter(router Router) {
	for _, rec := range r.routes {
		var allMW []MiddlewareFunc
		allMW = append(allMW, rec.middlewares...)
		wrapped := make([]HandlerFunc, len(allMW)+len(rec.handlers))

		for i, mw := range allMW {
			wrapped[i] = mw
		}
		copy(wrapped[len(allMW):], rec.handlers)

		switch rec.method {
		case "GET":
			router.GET(rec.path, wrapped...)
		case "POST":
			router.POST(rec.path, wrapped...)
		case "PUT":
			router.PUT(rec.path, wrapped...)
		case "DELETE":
			router.DELETE(rec.path, wrapped...)
		case "PATCH":
			router.PATCH(rec.path, wrapped...)
		case "WS":
			router.WS(rec.path, wrapped...)
		}
	}

	for _, g := range r.groups {
		grouter := g.Router.(*groupRouter)
		for _, rec := range grouter.routes {
			var allMW []MiddlewareFunc
			allMW = append(allMW, g.middlewares...)
			allMW = append(allMW, rec.middlewares...)
			wrapped := make([]HandlerFunc, len(allMW)+len(rec.handlers))

			for i, mw := range allMW {
				wrapped[i] = mw
			}
			copy(wrapped[len(allMW):], rec.handlers)

			fullPath := grouter.prefix + rec.path
			switch rec.method {
			case "GET":
				router.GET(fullPath, wrapped...)
			case "POST":
				router.POST(fullPath, wrapped...)
			case "PUT":
				router.PUT(fullPath, wrapped...)
			case "DELETE":
				router.DELETE(fullPath, wrapped...)
			case "PATCH":
				router.PATCH(fullPath, wrapped...)
			case "WS":
				router.WS(fullPath, wrapped...)
			}
		}
	}
}

type groupRouter struct {
	impl      *routerImpl
	prefix    string
	middlewares []MiddlewareFunc
	routes    []routeRecord
}

func (g *groupRouter) GET(path string, handlers ...HandlerFunc) *Route {
	g.routes = append(g.routes, routeRecord{
		method:     "GET",
		path:       path,
		handlers:   handlers,
		middlewares: g.middlewares,
	})
	return &Route{
		path:       path,
		method:     "GET",
		handlers:   handlers,
		middlewares: g.middlewares,
		routerImpl: g.impl,
	}
}

func (g *groupRouter) POST(path string, handlers ...HandlerFunc) *Route {
	g.routes = append(g.routes, routeRecord{
		method:     "POST",
		path:       path,
		handlers:   handlers,
		middlewares: g.middlewares,
	})
	return &Route{
		path:       path,
		method:     "POST",
		handlers:   handlers,
		middlewares: g.middlewares,
		routerImpl: g.impl,
	}
}

func (g *groupRouter) PUT(path string, handlers ...HandlerFunc) *Route {
	g.routes = append(g.routes, routeRecord{
		method:     "PUT",
		path:       path,
		handlers:   handlers,
		middlewares: g.middlewares,
	})
	return &Route{
		path:       path,
		method:     "PUT",
		handlers:   handlers,
		middlewares: g.middlewares,
		routerImpl: g.impl,
	}
}

func (g *groupRouter) DELETE(path string, handlers ...HandlerFunc) *Route {
	g.routes = append(g.routes, routeRecord{
		method:     "DELETE",
		path:       path,
		handlers:   handlers,
		middlewares: g.middlewares,
	})
	return &Route{
		path:       path,
		method:     "DELETE",
		handlers:   handlers,
		middlewares: g.middlewares,
		routerImpl: g.impl,
	}
}

func (g *groupRouter) PATCH(path string, handlers ...HandlerFunc) *Route {
	g.routes = append(g.routes, routeRecord{
		method:     "PATCH",
		path:       path,
		handlers:   handlers,
		middlewares: g.middlewares,
	})
	return &Route{
		path:       path,
		method:     "PATCH",
		handlers:   handlers,
		middlewares: g.middlewares,
		routerImpl: g.impl,
	}
}
```

- [ ] **步骤 2：Commit**

```bash
git add httpx/router.go
git commit -m "feat: add routerImpl and groupRouter for internal routing"
```

---

## 任务 4：创建 route.go（Route 类型）

**文件：**
- 创建：`httpx/route.go`

- [ ] **步骤 1：编写 route.go**

```go
package httpx

type Route struct {
	path       string
	method     string
	handlers   []HandlerFunc
	middlewares []MiddlewareFunc
	routerImpl *routerImpl
}

func (r *Route) Use(mw ...MiddlewareFunc) *Route {
	r.middlewares = append(r.middlewares, mw...)
	return r
}
```

- [ ] **步骤 2：Commit**

```bash
git add httpx/route.go
git commit -m "feat: add Route type with Use() for single-route middleware"
```

---

## 任务 5：创建 app.go（App 类型）

**文件：**
- 创建：`httpx/app.go`

- [ ] **步骤 1：编写 app.go**

```go
package httpx

type App struct {
	router      *routerImpl
	middlewares []MiddlewareFunc
}

func new() *App {
	return &App{
		router: newRouter(),
	}
}

func (a *App) GET(path string, handlers ...HandlerFunc) *Route {
	return a.router.GET(path, handlers...)
}

func (a *App) POST(path string, handlers ...HandlerFunc) *Route {
	return a.router.POST(path, handlers...)
}

func (a *App) PUT(path string, handlers ...HandlerFunc) *Route {
	return a.router.PUT(path, handlers...)
}

func (a *App) DELETE(path string, handlers ...HandlerFunc) *Route {
	return a.router.DELETE(path, handlers...)
}

func (a *App) PATCH(path string, handlers ...HandlerFunc) *Route {
	return a.router.PATCH(path, handlers...)
}

func (a *App) WS(path string, handlers ...HandlerFunc) *Route {
	return a.router.WS(path, handlers...)
}

func (a *App) Group(prefix string, mw ...MiddlewareFunc) *RouterGroup {
	return a.router.GROUP(prefix, mw...)
}

func (a *App) Use(mw ...MiddlewareFunc) *App {
	a.middlewares = append(a.middlewares, mw...)
	return a
}

func (a *App) Run(configPath string) error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
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

	addr := ":" + string(rune(cfg.Port))
	return srv.StartWithGraceful()
}
```

- [ ] **步骤 2：运行测试确认无编译错误**

运行：`go build ./httpx/`
预期：PASS（可能有类型错误，需要调整）

- [ ] **步骤 3：修复编译错误（如果有）**

预期问题：端口转换需要用 fmt.Sprintf

```go
addr := fmt.Sprintf(":%d", cfg.Port)
```

- [ ] **步骤 4：Commit**

```bash
git add httpx/app.go
git commit -m "feat: add App type with Run() for config-driven startup"
```

---

## 任务 6：修改 Gin adapter（实现 AdapterFactory）

**文件：**
- 修改：`httpx/adapter/gin/server.go`

- [ ] **步骤 1：在 gin server.go 末尾添加 AdapterFactory 实现**

```go
func NewServer(opts ...ServerOption) *GinServer {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	srv := &GinServer{
		engine:  engine,
		router:  newGinRouter(engine),
		httpSrv: &http.Server{Addr: "", Handler: engine},
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

func init() {
	RegisterAdapter("gin", func() httpx.Server {
		return NewServer()
	})
}
```

注意：检查是否已有 NewServer 函数。如果已有，确认函数签名。

- [ ] **步骤 2：Commit**

```bash
git add httpx/adapter/gin/server.go
git commit -m "feat: register gin adapter in adapter_registry"
```

---

## 任务 7：修改 Hertz adapter（实现 AdapterFactory）

**文件：**
- 修改：`httpx/adapter/hertz/server.go`

- [ ] **步骤 1：在 hertz server.go 末尾添加 AdapterFactory 实现**

```go
func init() {
	RegisterAdapter("hertz", func() httpx.Server {
		return NewServer()
	})
}
```

- [ ] **步骤 2：Commit**

```bash
git add httpx/adapter/hertz/server.go
git commit -m "feat: register hertz adapter in adapter_registry"
```

---

## 任务 8：修改 nethttp adapter（实现 AdapterFactory）

**文件：**
- 修改：`httpx/adapter/nethttp/server.go`

- [ ] **步骤 1：在 nethttp server.go 末尾添加 AdapterFactory 实现**

```go
func init() {
	RegisterAdapter("nethttp", func() httpx.Server {
		return NewServer()
	})
}
```

- [ ] **步骤 2：Commit**

```bash
git add httpx/adapter/nethttp/server.go
git commit -m "feat: register nethttp adapter in adapter_registry"
```

---

## 任务 9：创建 demo 示例

**文件：**
- 创建：`examples/demo/main.go`
- 创建：`examples/demo/route.go`
- 创建：`examples/demo/config.yaml`
- 创建：`examples/demo/config_hertz.yaml`

- [ ] **步骤 1：创建 config.yaml**

```yaml
adapter: gin
port: 8080
```

- [ ] **步骤 2：创建 config_hertz.yaml**

```yaml
adapter: hertz
port: 8080
```

- [ ] **步骤 3：创建 route.go**

```go
package main

import (
	"context"
	"fmt"

	"github.com/gospacex/httpx"
)

func setupRoutes(app *httpx.App) {
	app.Use(recoveryMiddleware, loggerMiddleware)

	app.GET("/health", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{"status": "ok"})
	})

	app.GET("/hello", func(ctx context.Context, hc httpx.HandlerContext) {
		name := hc.Query("name")
		if name == "" {
			name = "World"
		}
		hc.AbortJSON(200, map[string]string{"message": fmt.Sprintf("Hello, %s!", name)})
	})
}

func recoveryMiddleware(next httpx.HandlerFunc) httpx.HandlerFunc {
	return func(ctx context.Context, hc httpx.HandlerContext) {
		defer func() {
			if r := recover(); r != nil {
				hc.AbortJSON(500, map[string]string{"error": "internal error"})
			}
		}()
		next(ctx, hc)
	}
}

func loggerMiddleware(next httpx.HandlerFunc) httpx.HandlerFunc {
	return func(ctx context.Context, hc httpx.HandlerContext) {
		fmt.Println("Request received")
		next(ctx, hc)
	}
}
```

- [ ] **步骤 4：创建 main.go**

```go
package main

func main() {
	app := httpx.new()
	setupRoutes(app)
	if err := app.Run("./config.yaml"); err != nil {
		panic(err)
	}
}
```

- [ ] **步骤 5：运行 gin 测试**

运行：`cd examples/demo && go run main.go route.go`（后台运行）
curl 测试：
```bash
curl http://localhost:8080/health
curl http://localhost:8080/hello?name=Alice
```

- [ ] **步骤 6：运行 hertz 测试**

运行：`cd examples/demo && cp config_hertz.yaml config.yaml && go run main.go route.go`（后台运行）
curl 测试相同

- [ ] **步骤 7：Commit**

```bash
git add examples/demo/
git commit -m "feat: add demo example with gin and hertz support"
```

---

## 任务 10：全量编译验证

- [ ] **步骤 1：运行 `go build ./...` 确认无编译错误**

运行：`go build ./...`
预期：PASS

- [ ] **步骤 2：运行 `go vet ./...` 确认无警告**

运行：`go vet ./...`
预期：无警告

- [ ] **步骤 3：运行 `go test ./...` 确认测试通过**

运行：`go test ./...`
预期：PASS

- [ ] **步骤 4：Commit**

```bash
git add -A && git commit -m "chore: pass all build and test checks"
```

---

## 测试覆盖要求

测试必须覆盖 **gin** 和 **hertz** 两种 adapter 场景：

### 测试用例 1：App 创建和路由注册

```go
func TestGinApp_RouteRegistration(t *testing.T) {
	cfg, err := LoadConfig("testdata/config_gin.yaml")
	if err != nil {
		t.Fatal(err)
	}
	factory, _ := getAdapter(cfg.Adapter)
	srv := factory()
	router := srv.Router()
	// 注册路由并验证
}

func TestHertzApp_RouteRegistration(t *testing.T) {
	cfg, err := LoadConfig("testdata/config_hertz.yaml")
	if err != nil {
		t.Fatal(err)
	}
	factory, _ := getAdapter(cfg.Adapter)
	srv := factory()
	router := srv.Router()
	// 注册路由并验证
}
```

### 测试用例 2：全局中间件

```go
func TestGinApp_GlobalMiddleware(t *testing.T) {
	mwCalled := false
	mw := func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) {
			mwCalled = true
			next(ctx, hc)
		}
	}
	app := new()
	app.Use(mw)
	// 验证中间件被调用
}

func TestHertzApp_GlobalMiddleware(t *testing.T) {
	// 同上
}
```

### 测试用例 3：单路由中间件

```go
func TestGinApp_RouteMiddlewareChain(t *testing.T) {
	app := new()
	app.GET("/test", handler).Use(mw1, mw2)
	// 验证链式调用
}

func TestHertzApp_RouteMiddlewareChain(t *testing.T) {
	// 同上
}
```

### 测试用例 4：分组路由

```go
func TestGinApp_GroupRouting(t *testing.T) {
	app := new()
	g := app.Group("/api", mw1)
	g.GET("/users", handler)
	// 验证分组路由
}

func TestHertzApp_GroupRouting(t *testing.T) {
	// 同上
}
```

### 测试用例 5：启动和停止

```go
func TestGinApp_GracefulShutdown(t *testing.T) {
	app := new()
	app.GET("/health", handler)
	go app.Run("testdata/config_gin.yaml")
	time.Sleep(100 * time.Millisecond)
	// 验证服务运行
	app.GracefulShutdown(context.Background())
}

func TestHertzApp_GracefulShutdown(t *testing.T) {
	// 同上
}
```

---

## 依赖

```
github.com/spf13/viper v1.18.0
github.com/gin-gonic/gin v1.10.0
github.com/cloudwego/hertz v0.10.0
```