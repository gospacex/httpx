# httpx Examples

示例代码展示 httpx 库的各种用法。

## demo

一个完整的 HTTP 服务演示，可通过 curl 访问：

```bash
go run examples/demo/main.go
```

启动后访问 http://localhost:8080

### 可用接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/` | 首页信息 |
| GET | `/health` | 健康检查 |
| GET | `/hello?name=xxx` | 欢迎页 |
| GET | `/time` | 服务器时间 |
| GET | `/users` | 列出所有用户 |
| GET | `/users/:id` | 获取单个用户 |
| POST | `/users` | 创建用户 |
| PUT | `/users/:id` | 更新用户 |
| DELETE | `/users/:id` | 删除用户 |
| GET | `/articles` | 列出所有文章 |
| GET | `/articles/:id` | 获取单篇文章 |
| POST | `/articles` | 创建文章 |
| DELETE | `/articles/:id` | 删除文章 |

## basic

Gin adapter 的基本用法演示：

- 服务器启动和停止
- 路由注册（GET/POST/PUT/DELETE/PATCH）
- 中间件使用
- HandlerContext（Param/Query/AbortWithStatus/AbortJSON）
- 优雅关闭

### 运行

```bash
go run examples/basic/main.go
```

## hertz

Hertz adapter 演示，与 basic 类似。

## nethttp

nethttp adapter 演示，与 basic 类似。