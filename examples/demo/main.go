package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gospacex/httpx"
)

// User 数据结构
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Article 数据结构
type Article struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

var (
	users    = make(map[string]User)
	articles = make(map[string]Article)
)

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║   httpx HTTP 服务演示                    ║")
	fmt.Println("║   启动后访问 http://localhost:8080       ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// 创建应用（使用新的 App API）
	app := httpx.New()

	// 注册路由
	setupRoutes(app)

	// 启动服务
	if err := app.Run("./config.yaml"); err != nil {
		panic(err)
	}
}

func setupRoutes(app *httpx.App) {
	// 全局中间件
	app.Use(recoveryMiddleware, loggerMiddleware)

	// 根路径
	app.GET("/", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{
			"message":   "欢迎使用 httpx",
			"version":   "1.0.0",
			"endpoints": "/users, /articles, /hello, /time, /health",
		})
	})

	// 健康检查
	app.GET("/health", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 欢迎页
	app.GET("/hello", func(ctx context.Context, hc httpx.HandlerContext) {
		name := hc.Query("name")
		if name == "" {
			name = "World"
		}
		hc.AbortJSON(200, map[string]string{
			"message": fmt.Sprintf("你好, %s!", name),
		})
	})

	// 当前时间
	app.GET("/time", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{
			"server_time": time.Now().Format(time.RFC3339),
			"unix":        fmt.Sprintf("%d", time.Now().Unix()),
		})
	})

	// 用户路由 - 直接注册
	setupUsersRoutes(app)

	// 文章路由 - 直接注册
	setupArticlesRoutes(app)
}

func setupUsersRoutes(app *httpx.App) {
	// 列出所有用户
	app.GET("/users", func(ctx context.Context, hc httpx.HandlerContext) {
		userList := make([]User, 0, len(users))
		for _, u := range users {
			userList = append(userList, u)
		}
		hc.AbortJSON(200, map[string]interface{}{
			"count": len(users),
			"users": userList,
		})
	})

	// 获取单个用户
	app.GET("/users/:id", func(ctx context.Context, hc httpx.HandlerContext) {
		id := hc.Param("id")
		if user, ok := users[id]; ok {
			hc.AbortJSON(200, user)
		} else {
			hc.AbortJSON(404, map[string]string{"error": "用户不存在"})
		}
	})

	// 创建用户
	app.POST("/users", func(ctx context.Context, hc httpx.HandlerContext) {
		var newUser User
		if err := hc.Bind(&newUser); err != nil {
			hc.AbortJSON(400, map[string]string{"error": "无效的请求数据"})
			return
		}
		if newUser.ID == "" {
			newUser.ID = fmt.Sprintf("%d", len(users)+1)
		}
		users[newUser.ID] = newUser
		hc.AbortJSON(201, map[string]interface{}{
			"message": "用户创建成功",
			"user":    newUser,
		})
	})

	// 更新用户
	app.PUT("/users/:id", func(ctx context.Context, hc httpx.HandlerContext) {
		id := hc.Param("id")
		var updated User
		if err := hc.Bind(&updated); err != nil {
			hc.AbortJSON(400, map[string]string{"error": "无效的请求数据"})
			return
		}
		if _, ok := users[id]; !ok {
			hc.AbortJSON(404, map[string]string{"error": "用户不存在"})
			return
		}
		updated.ID = id
		users[id] = updated
		hc.AbortJSON(200, map[string]interface{}{
			"message": "用户更新成功",
			"user":    updated,
		})
	})

	// 删除用户
	app.DELETE("/users/:id", func(ctx context.Context, hc httpx.HandlerContext) {
		id := hc.Param("id")
		if _, ok := users[id]; !ok {
			hc.AbortJSON(404, map[string]string{"error": "用户不存在"})
			return
		}
		delete(users, id)
		hc.AbortJSON(200, map[string]string{"message": "用户删除成功"})
	})
}

func setupArticlesRoutes(app *httpx.App) {
	// 列出所有文章
	app.GET("/articles", func(ctx context.Context, hc httpx.HandlerContext) {
		articleList := make([]Article, 0, len(articles))
		for _, a := range articles {
			articleList = append(articleList, a)
		}
		hc.AbortJSON(200, map[string]interface{}{
			"count":    len(articles),
			"articles": articleList,
		})
	})

	// 获取单个文章
	app.GET("/articles/:id", func(ctx context.Context, hc httpx.HandlerContext) {
		id := hc.Param("id")
		if article, ok := articles[id]; ok {
			hc.AbortJSON(200, article)
		} else {
			hc.AbortJSON(404, map[string]string{"error": "文章不存在"})
		}
	})

	// 创建文章
	app.POST("/articles", func(ctx context.Context, hc httpx.HandlerContext) {
		var newArticle Article
		if err := hc.Bind(&newArticle); err != nil {
			hc.AbortJSON(400, map[string]string{"error": "无效的请求数据"})
			return
		}
		if newArticle.ID == "" {
			newArticle.ID = fmt.Sprintf("%d", len(articles)+1)
		}
		articles[newArticle.ID] = newArticle
		hc.AbortJSON(201, map[string]interface{}{
			"message": "文章创建成功",
			"article": newArticle,
		})
	})

	// 删除文章
	app.DELETE("/articles/:id", func(ctx context.Context, hc httpx.HandlerContext) {
		id := hc.Param("id")
		if _, ok := articles[id]; !ok {
			hc.AbortJSON(404, map[string]string{"error": "文章不存在"})
			return
		}
		delete(articles, id)
		hc.AbortJSON(200, map[string]string{"message": "文章删除成功"})
	})
}

// 全局恢复中间件
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

// 全局日志中间件
func loggerMiddleware(next httpx.HandlerFunc) httpx.HandlerFunc {
	return func(ctx context.Context, hc httpx.HandlerContext) {
		start := time.Now()
		fmt.Printf("📥 [%s] 请求开始\n", time.Now().Format("15:04:05"))
		next(ctx, hc)
		fmt.Printf("📤 [%s] 请求完成 (耗时: %v)\n", time.Now().Format("15:04:05"), time.Since(start))
	}
}