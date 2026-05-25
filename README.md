# httpx - 统一 HTTP 框架适配层

通过统一的接口，同时支持 Gin、Hertz 和 net/http。

## 快速开始

### 安装

```bash
go get github.com/gospacex/httpx
```

### 基础用法

```go
package main

import (
    "context"
    "fmt"
    "github.com/gospacex/httpx"
    _ "github.com/gospacex/httpx/adapter/gin"
)

func main() {
    // 通过配置文件创建应用
    app, err := httpx.New("./config.yaml")
    if err != nil {
        panic(err)
    }

    // 注册路由
    app.GET("/hello", func(ctx context.Context, hc httpx.HandlerContext) {
        hc.AbortJSON(200, map[string]string{"message": "Hello, World!"})
    })

    // 启动服务
    app.Run()
}
```

### 配置文件

```yaml
adapter: gin
port: 8080
```

### 多框架支持

同一套代码，轻松切换框架：

```yaml
# Gin
adapter: gin
port: 8080
```

```yaml
# Hertz
adapter: hertz
port: 8080
```

```yaml
# net/http
adapter: nethttp
port: 8080
```

### 中间件

```go
// 全局中间件
app.Use(loggingMiddleware, recoveryMiddleware)

// 路由组中间件
admin := app.Group("/admin", authMiddleware)
admin.GET("/dashboard", dashboardHandler)
```

### 路由组

```go
api := app.Group("/api/v1", func(next httpx.HandlerFunc) httpx.HandlerFunc {
    return func(ctx context.Context, hc httpx.HandlerContext) {
        // 前置处理
        next(ctx, hc)
        // 后置处理
    }
})

api.GET("/users", usersHandler)
api.GET("/articles", articlesHandler)
```

## 架构

```
httpx
├── app.go              # App, New(), Run()
├── router.go           # Router, RouterGroup
├── types.go            # HandlerFunc, HandlerContext, MiddlewareFunc
├── config.go           # 配置加载
├── adapter_registry.go # 适配器注册
├── adapter/
│   ├── gin/           # Gin 适配器
│   ├── hertz/         # Hertz 适配器
│   └── nethttp/       # NetHTTP 适配器
└── examples/demo/      # 完整示例
```

## 适配器

| 适配器 | 路径 | 特点 |
|--------|------|------|
| Gin | `adapter/gin` | 高性能，成熟稳定 |
| Hertz | `adapter/hertz` | 字节跳动框架 |
| net/http | `adapter/nethttp` | 标准库，无依赖 |

## 接口定义

### App

```go
// 创建应用（加载配置和适配器）
func New(configPath string) (*App, error)

// 启动服务（使用配置中的端口）
func (a *App) Run() error

// 启动服务（指定地址）
func (a *App) RunOnAddr(addr string) error

// 全局中间件
func (a *App) Use(mw ...MiddlewareFunc) *App

// 路由组
func (a *App) Group(prefix string, mw ...MiddlewareFunc) *RouterGroup
```

### Server

```go
type Server interface {
    Start(addr string) error
    Stop(ctx context.Context) error
    Router() Router
    Use(middlewares ...MiddlewareFunc) Server
    EnableWS(wsConfig *WSConfig) Server
}
```

### Router

```go
type Router interface {
    GET(path string, handlers ...HandlerFunc) Router
    POST(path string, handlers ...HandlerFunc) Router
    PUT(path string, handlers ...HandlerFunc) Router
    DELETE(path string, handlers ...HandlerFunc) Router
    PATCH(path string, handlers ...HandlerFunc) Router
    WS(path string, handlers ...HandlerFunc) Router
    GROUP(prefix string, mw ...MiddlewareFunc) *RouterGroup
}
```

### HandlerContext

```go
type HandlerContext interface {
    Request() interface{}
    Response() interface{}
    Param(key string) string
    Query(key string) string
    Bind(into interface{}) error
    AbortWithStatus(code int)
    AbortJSON(code int, body interface{})
}
```

## 测试

```bash
go test -v ./...
```

## 示例

参见 `examples/demo/` 目录，包含完整的路由、中间件和路由组示例。

## 许可证

MIT
