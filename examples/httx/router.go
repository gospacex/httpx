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

func setupRoutes(app *httpx.App) {
	app.Use(recoveryMiddleware, loggerMiddleware)

	app.GET("/", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, map[string]string{
			"message":   "欢迎使用 httpx",
			"version":   "1.0.0",
			"endpoints": "/users, /articles, /hello, /time, /health, /api/v1/*, /admin/*",
		})
		return nil
	})

	app.GET("/health", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
		return nil
	})

	app.GET("/hello", func(ctx context.Context, hc httpx.HandlerContext) error {
		name := hc.Query("name")
		if name == "" {
			name = "World"
		}
		hc.AbortJSON(200, map[string]string{
			"message": fmt.Sprintf("你好, %s!", name),
		})
		return nil
	})

	app.GET("/time", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, map[string]string{
			"server_time": time.Now().Format(time.RFC3339),
			"unix":        fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil
	})

	// 路由分组示例：/api/v1 分组
	setupAPIV1Routes(app)
	// 路由分组示例：/admin 分组
	setupAdminRoutes(app)

	setupUsersRoutes(app)
	setupArticlesRoutes(app)
}

func setupUsersRoutes(app *httpx.App) {
	app.GET("/users", func(ctx context.Context, hc httpx.HandlerContext) error {
		userList := make([]User, 0, len(users))
		for _, u := range users {
			userList = append(userList, u)
		}
		hc.AbortJSON(200, map[string]interface{}{
			"count": len(users),
			"users": userList,
		})
		return nil
	})

	app.GET("/users/:id", func(ctx context.Context, hc httpx.HandlerContext) error {
		id := hc.Param("id")
		if user, ok := users[id]; ok {
			hc.AbortJSON(200, user)
		} else {
			hc.AbortJSON(404, map[string]string{"error": "用户不存在"})
		}
		return nil
	})

	app.POST("/users", func(ctx context.Context, hc httpx.HandlerContext) error {
		var newUser User
		if err := hc.Bind(&newUser); err != nil {
			hc.AbortJSON(400, map[string]string{"error": "无效的请求数据"})
			return nil
		}
		if newUser.ID == "" {
			newUser.ID = fmt.Sprintf("%d", len(users)+1)
		}
		users[newUser.ID] = newUser
		hc.AbortJSON(201, map[string]interface{}{
			"message": "用户创建成功",
			"user":    newUser,
		})
		return nil
	})

	app.PUT("/users/:id", func(ctx context.Context, hc httpx.HandlerContext) error {
		id := hc.Param("id")
		var updated User
		if err := hc.Bind(&updated); err != nil {
			hc.AbortJSON(400, map[string]string{"error": "无效的请求数据"})
			return nil
		}
		if _, ok := users[id]; !ok {
			hc.AbortJSON(404, map[string]string{"error": "用户不存在"})
			return nil
		}
		updated.ID = id
		users[id] = updated
		hc.AbortJSON(200, map[string]interface{}{
			"message": "用户更新成功",
			"user":    updated,
		})
		return nil
	})

	app.DELETE("/users/:id", func(ctx context.Context, hc httpx.HandlerContext) error {
		id := hc.Param("id")
		if _, ok := users[id]; !ok {
			hc.AbortJSON(404, map[string]string{"error": "用户不存在"})
			return nil
		}
		delete(users, id)
		hc.AbortJSON(200, map[string]string{"message": "用户删除成功"})
		return nil
	})
}

func setupArticlesRoutes(app *httpx.App) {
	app.GET("/articles", func(ctx context.Context, hc httpx.HandlerContext) error {
		articleList := make([]Article, 0, len(articles))
		for _, a := range articles {
			articleList = append(articleList, a)
		}
		hc.AbortJSON(200, map[string]interface{}{
			"count":    len(articles),
			"articles": articleList,
		})
		return nil
	})

	app.GET("/articles/:id", func(ctx context.Context, hc httpx.HandlerContext) error {
		id := hc.Param("id")
		if article, ok := articles[id]; ok {
			hc.AbortJSON(200, article)
		} else {
			hc.AbortJSON(404, map[string]string{"error": "文章不存在"})
		}
		return nil
	})

	app.POST("/articles", func(ctx context.Context, hc httpx.HandlerContext) error {
		var newArticle Article
		if err := hc.Bind(&newArticle); err != nil {
			hc.AbortJSON(400, map[string]string{"error": "无效的请求数据"})
			return nil
		}
		if newArticle.ID == "" {
			newArticle.ID = fmt.Sprintf("%d", len(articles)+1)
		}
		articles[newArticle.ID] = newArticle
		hc.AbortJSON(201, map[string]interface{}{
			"message": "文章创建成功",
			"article": newArticle,
		})
		return nil
	})

	app.DELETE("/articles/:id", func(ctx context.Context, hc httpx.HandlerContext) error {
		id := hc.Param("id")
		if _, ok := articles[id]; !ok {
			hc.AbortJSON(404, map[string]string{"error": "文章不存在"})
			return nil
		}
		delete(articles, id)
		hc.AbortJSON(200, map[string]string{"message": "文章删除成功"})
		return nil
	})
}

func recoveryMiddleware(next httpx.HandlerFunc) httpx.HandlerFunc {
	return func(ctx context.Context, hc httpx.HandlerContext) error {
		defer func() {
			if r := recover(); r != nil {
				hc.AbortJSON(500, map[string]string{"error": "internal error"})
			}
		}()
		return next(ctx, hc)
	}
}

func loggerMiddleware(next httpx.HandlerFunc) httpx.HandlerFunc {
	return func(ctx context.Context, hc httpx.HandlerContext) error {
		start := time.Now()
		fmt.Printf("📥 [%s] 请求开始\n", time.Now().Format("15:04:05"))
		if err := next(ctx, hc); err != nil {
			return err
		}
		fmt.Printf("📤 [%s] 请求完成 (耗时: %v)\n", time.Now().Format("15:04:05"), time.Since(start))
		return nil
	}
}

// API v1 分组路由示例 - 演示分组级别中间件
func setupAPIV1Routes(app *httpx.App) {
	apiMw := func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) error {
			fmt.Println("[api-v1-mw] handling request")
			return next(ctx, hc)
		}
	}

	api := app.Group("/api/v1", apiMw)

	api.Router.GET("/info", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, map[string]string{
			"version": "v1",
			"status":  "running",
		})
		return nil
	})

	api.Router.GET("/timestamp", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, map[string]string{
			"unix": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil
	})
}

// Admin 分组路由示例 - 演示分组级别中间件（多个中间件通过闭包链式调用）
func setupAdminRoutes(app *httpx.App) {
	adminMw := func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) error {
			token := hc.Query("token")
			if token != "admin123" {
				hc.AbortJSON(401, map[string]string{"error": "unauthorized"})
				return nil
			}
			return next(ctx, hc)
		}
	}

	admin := app.Group("/admin", adminMw)

	admin.Router.GET("/dashboard", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, map[string]string{
			"message": "admin dashboard",
			"stats":   "system healthy",
		})
		return nil
	})

	admin.Router.GET("/stats", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, map[string]interface{}{
			"users":    len(users),
			"articles": len(articles),
		})
		return nil
	})
}