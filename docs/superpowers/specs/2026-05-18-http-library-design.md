# HTTP 库设计规格

## 目标

- 支持多 HTTP 框架（net/http, Gin, Hertz）
- 支持子路径路由（多版本/模块）
- 支持 WebSocket 升级（Server 层全局 + 路由层细粒度）
- 使用 wire 进行依赖注入
- 作为 GPX 脚手架的 BFF 层依赖库
- 团队内部开源

## 架构

### 核心接口

```go
// Server HTTP服务器接口
type Server interface {
    Start(addr string) error
    Stop(ctx context.Context) error
    Router() Router
    Use(middleware ...MiddlewareFunc) Server
    EnableWS(wsConfig *WSConfig) Server
}

// Router 路由接口
type Router interface {
    GET(path string, handlers ...HandlerFunc) Router
    POST(path string, handlers ...HandlerFunc) Router
    PUT(path string, handlers ...HandlerFunc) Router
    DELETE(path string, handlers ...HandlerFunc) Router
    PATCH(path string, handlers ...HandlerFunc) Router
    GROUP(prefix string, mw ...MiddlewareFunc) *RouterGroup
    WS(path string, handlers ...HandlerFunc) Router
}

// RouterGroup 路由组
type RouterGroup struct {
    Router Router
    Prefix string
}
```

### 适配器模式

```
┌─────────────────────────────────────────────────┐
│                   User Code                      │
│    router.GET("/users", handler)                │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────┐
│              Router interface                    │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────┐
│           ServerFactory                         │
└──────────────────┬──────────────────────────────┘
                   │
         ┌─────────┴─────────┐
         ▼                   ▼
┌─────────────────┐  ┌─────────────────┐
│  NetHTTPServer   │  │   GinServer     │
│  (adapter)       │  │   (adapter)    │
└─────────────────┘  └─────────────────┘
```

## 中间件组织

### 全局中间件（通过 wire 注入）

```go
type GlobalMiddleware struct {
    Logger     *LogMiddleware
    Recover    *RecoverMiddleware
    CORS       *CORSMiddleware
    RateLimit  *RateLimitMiddleware
}
```

### 单服务中间件（运行时按配置注入）

```go
type ServiceMiddleware struct {
    Auth *AuthMiddleware
}
```

## WebSocket 支持

### 方式1: Server层全局启用

```go
app := server.NewServer(
    server.WithWSEnable(wsConfig),
)
```

### 方式2: 路由层细粒度控制

```go
router := app.Router()
router.WS("/ws/chat", chatHandler)
router.GET("/ws/stream", streamHandler)
router.POST("/api/rest", restHandler)
```

## 子路径路由

```go
v1 := app.Router().GROUP("/api/v1")
v1.GET("/users", listUsers)

v2 := app.Router().GROUP("/api/v2")
v2.GET("/users", listUsersV2)
```

## 错误处理

开发环境 panic 透传，生产环境统一 Recover 返回 500 并记录栈。

## 目录结构

```
httpx/
├── server.go           # Server 接口
├── router.go           # Router 接口
├── adapter/
│   ├── nethttp/        # net/http 适配器
│   ├── gin/            # Gin 适配器
│   └── hertz/          # Hertz 适配器
├── middleware/
│   ├── logger.go
│   ├── recovery.go
│   ├── cors.go
│   └── ratelimit.go
├── wire/
│   ├── provider.go     # wire provider
│   └── wire.go         # wire 生成指令
└── websocket/
    └── upgrade.go      # WebSocket 升级工具
```

## wire 依赖注入

```go
var ServerSet = wire.NewSet(
    NewServerFactory,
    wire.Bind(new(Server), new(*GinServer)),
    GlobalMiddlewareProvider,
)

var MiddlewareSet = wire.NewSet(
    NewLogMiddleware,
    NewRecoverMiddleware,
    NewCORSMiddleware,
    NewRateLimitMiddleware,
)
```

## 测试策略

### 单元测试

- **接口 Mock**：为 `Server`、`Router` 等接口编写 mock 实现，用于测试业务逻辑不依赖具体 HTTP 框架
- **适配器测试**：每个适配器（nethttp/gin/hertz）独立测试，验证接口实现正确性
- **中间件测试**：独立测试每个中间件，无需启动服务器
- **工具函数测试**：WebSocket 升级、错误处理等工具函数直接测试

```go
// 适配器测试示例
func TestGinServer_Start(t *testing.T) {
    server := NewGinServer()
    // mock Router
    // verify Start/Stop behavior
}

// 中间件测试示例
func TestRecoverMiddleware(t *testing.T) {
    middleware := NewRecoverMiddleware()
    handler := middleware.Handle(nextHandler)
    // panic 场景验证
}
```

### 端到端测试

- **HTTP 层测试**：启动真实服务器，发送真实 HTTP 请求，验证完整请求链路
- **使用 httptest**：Go 标准库的 `httptest` 模拟请求，无需真实监听端口
- **测试覆盖场景**：
  - 路由注册与匹配
  - 中间件链执行顺序
  - 全局 vs 路由层 WebSocket
  - 子路径路由 GROUP
  - 错误响应格式

```go
func TestE2E_RouterGroup(t *testing.T) {
    srv := NewGinServer()
    router := srv.Router()

    v1 := router.GROUP("/api/v1")
    v1.GET("/users", func(c context.Context, hc HandlerContext) {})

    req := httptest.NewRequest("GET", "/api/v1/users", nil)
    w := httptest.NewRecorder()
    srv.ServeHTTP(w, req)

    // verify response
}

func TestE2E_WebSocket(t *testing.T) {
    srv := NewGinServer()
    srv.EnableWS(wsConfig)
    router := srv.Router()
    router.WS("/ws", wsHandler)

    // upgrade request and verify ws connection
}
```

### 测试目录结构

```
httpx/
├── server.go
├── router.go
├── adapter/
│   ├── nethttp/
│   │   └── nethttp_test.go
│   ├── gin/
│   │   └── gin_test.go
│   └── hertz/
│       └── hertz_test.go
├── middleware/
│   ├── logger_test.go
│   ├── recovery_test.go
│   ├── cors_test.go
│   └── ratelimit_test.go
├── wire/
│   ├── provider_test.go   # wire 生成的 provider 测试
│   └── wire.go
├── websocket/
│   └── upgrade_test.go
└── e2e/
    └── server_test.go    # 端到端测试

### 压力测试

使用 `wrk` 或 Go 内置 benchmark 进行性能测试。

- **Benchmark 测试**：Go 标准库 `testing.B` 进行微基准测试
- **HTTP 压测工具**：`wrk` 或 `ghz` 进行真实压力测试
- **测试指标**：
  - QPS（每秒请求数）
  - 延迟分布（P50/P90/P99）
  - 内存分配
  - 并发连接数

```go
// Benchmark 示例
func BenchmarkRouter_Get(b *testing.B) {
    srv := NewGinServer()
    router := srv.Router()
    router.GET("/users/:id", handler)

    // 使用 httptest 进行基准测试
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        req := httptest.NewRequest("GET", "/users/1", nil)
        w := httptest.NewRecorder()
        srv.ServeHTTP(w, req)
    }
}
```

```bash
# wrk 压测示例
wrk -t12 -c400 -d30s http://localhost:8080/api/v1/users

# ghz 压测示例
ghz --insecure --connections 100 --duration 30s http://localhost:8080/api/v1/users
```

### 压测场景

| 场景 | 说明 |
|------|------|
| 路由匹配 | 各种路径模式（静态/参数/通配）的匹配性能 |
| 中间件链 | 不同数量中间件的吞吐量影响 |
| WebSocket | 长连接并发数、消息吞吐量 |
| 并发连接 | 高并发下的 QPS 和延迟稳定性 |
```