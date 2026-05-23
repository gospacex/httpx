package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gospacex/httpx"
	"github.com/gospacex/httpx/adapter/gin"
)

// 用户数据结构
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 文章数据结构
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
	fmt.Println("║   httpx HTTP 服务演示                  ║")
	fmt.Println("║   启动后访问 http://localhost:8080     ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// 创建服务器
	srv := gin.NewServer()

	// 注册路由
	setupRoutes(srv)

	// 添加全局中间件
	setupMiddleware(srv)

	// 启动服务器
	go func() {
		fmt.Println("🚀 服务器启动中...")
		if err := srv.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待信号优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n📤 正在关闭服务器...")
	srv.Stop(context.Background())
	fmt.Println("✅ 服务器已关闭")
}

func setupRoutes(srv *gin.GinServer) {
	router := srv.Router()

	// 根路径
	router.GET("/", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{
			"message":   "欢迎使用 httpx",
			"version":   "1.0.0",
			"endpoints": "/users, /articles, /hello, /time",
		})
	})

	// 健康检查
	router.GET("/health", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 欢迎页
	router.GET("/hello", func(ctx context.Context, hc httpx.HandlerContext) {
		name := hc.Query("name")
		if name == "" {
			name = "World"
		}
		hc.AbortJSON(200, map[string]string{
			"message": fmt.Sprintf("你好, %s!", name),
		})
	})

	// 当前时间
	router.GET("/time", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{
			"server_time": time.Now().Format(time.RFC3339),
			"unix":        fmt.Sprintf("%d", time.Now().Unix()),
		})
	})

	// 用户路由 - 直接在 router 上注册
	setupUsersRoutes(router)

	// 文章路由 - 直接在 router 上注册
	setupArticlesRoutes(router)
}

func setupUsersRoutes(router httpx.Router) {
	// 列出所有用户
	router.GET("/users", func(ctx context.Context, hc httpx.HandlerContext) {
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
	router.GET("/users/:id", func(ctx context.Context, hc httpx.HandlerContext) {
		id := hc.Param("id")
		if user, ok := users[id]; ok {
			hc.AbortJSON(200, user)
		} else {
			hc.AbortJSON(404, map[string]string{"error": "用户不存在"})
		}
	})

	// 创建用户
	router.POST("/users", func(ctx context.Context, hc httpx.HandlerContext) {
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
	router.PUT("/users/:id", func(ctx context.Context, hc httpx.HandlerContext) {
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
	router.DELETE("/users/:id", func(ctx context.Context, hc httpx.HandlerContext) {
		id := hc.Param("id")
		if _, ok := users[id]; !ok {
			hc.AbortJSON(404, map[string]string{"error": "用户不存在"})
			return
		}
		delete(users, id)
		hc.AbortJSON(200, map[string]string{"message": "用户删除成功"})
	})
}

func setupArticlesRoutes(router httpx.Router) {
	// 列出所有文章
	router.GET("/articles", func(ctx context.Context, hc httpx.HandlerContext) {
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
	router.GET("/articles/:id", func(ctx context.Context, hc httpx.HandlerContext) {
		id := hc.Param("id")
		if article, ok := articles[id]; ok {
			hc.AbortJSON(200, article)
		} else {
			hc.AbortJSON(404, map[string]string{"error": "文章不存在"})
		}
	})

	// 创建文章
	router.POST("/articles", func(ctx context.Context, hc httpx.HandlerContext) {
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
	router.DELETE("/articles/:id", func(ctx context.Context, hc httpx.HandlerContext) {
		id := hc.Param("id")
		if _, ok := articles[id]; !ok {
			hc.AbortJSON(404, map[string]string{"error": "文章不存在"})
			return
		}
		delete(articles, id)
		hc.AbortJSON(200, map[string]string{"message": "文章删除成功"})
	})
}

func setupMiddleware(srv *gin.GinServer) {
	// 请求日志中间件
	srv.Use(func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) {
			start := time.Now()
			fmt.Printf("📥 [%s] 请求开始\n", time.Now().Format("15:04:05"))
			next(ctx, hc)
			fmt.Printf("📤 [%s] 请求完成 (耗时: %v)\n", time.Now().Format("15:04:05"), time.Since(start))
		}
	})

	// CORS 中间件
	srv.Use(func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) {
			fmt.Println("🔧 CORS 中间件")
			next(ctx, hc)
		}
	})
}