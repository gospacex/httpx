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
    "github.com/gospacex/httpx/adapter/gin"
)

func main() {
    // 创建服务器（使用 Gin 适配器）
    srv := gin.NewServer()

    // 注册路由
    router := srv.Router()
    router.GET("/hello", func(ctx context.Context, hc httpx.HandlerContext) {
        hc.AbortJSON(200, map[string]string{"message": "Hello, World!"})
    })

    // 启动服务器
    srv.Start(":8080")
}
```

### 多框架支持

同一套代码，轻松切换框架：

```go
// Gin
srv := gin.NewServer()

// Hertz  
srv := hertz.NewServer()

// net/http
srv := nethttp.NewServer()
```

### 路由组

```go
router := srv.Router()
group := router.GROUP("/api")

group.GET("/users", handler1)
group.POST("/users", handler2)

// 在组上添加中间件
group.Use(middleware1, middleware2)
```

### 中间件

```go
srv.Use(httpx.Logger(), httpx.Recovery())
```

## 适配器

| 适配器 | 路径 | 特点 |
|--------|------|------|
| Gin | `adapter/gin` | 高性能，成熟稳定 |
| Hertz | `adapter/hertz` | 字节跳动框架 |
| net/http | `adapter/nethttp` | 标准库，无依赖 |

## 接口定义

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

## 许可证

MIT
