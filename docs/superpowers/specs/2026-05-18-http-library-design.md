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