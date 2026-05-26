package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gospacex/httpx"
	"github.com/gospacex/httpx/adapter/gin"
)

func main() {
	fmt.Println("=== httpx Gin Adapter 示例 ===")
	fmt.Println()

	// 1. 基本服务器启动和停止
	fmt.Println("1. 基本服务器启动和停止")
	basicServer()

	// 2. 路由示例
	fmt.Println("\n2. 路由示例")
	routingExample()

	// 3. 中间件示例
	fmt.Println("\n3. 中间件示例")
	middlewareExample()

	// 4. HandlerContext 用法
	fmt.Println("\n4. HandlerContext 用法")
	handlerContextExample()

	// 5. 优雅关闭
	fmt.Println("\n5. 优雅关闭示例")
	gracefulShutdownExample()

	fmt.Println("\n=== 所有示例完成 ===")
}

func basicServer() {
	srv := gin.NewServer()

	go func() {
		err := srv.Start(":0")
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	time.Sleep(10 * time.Millisecond)
	fmt.Printf("  Server running: %v\n", srv.IsRunning())

	srv.Stop(context.Background())
	fmt.Println("  Server stopped")
}

func routingExample() {
	srv := gin.NewServer()
	router := srv.Router()

	var getHandler, postHandler, putHandler bool

	router.GET("/hello", func(ctx context.Context, hc httpx.HandlerContext) error {
		getHandler = true
		hc.AbortJSON(200, map[string]string{"message": "hello"})
		return nil
	})

	router.POST("/data", func(ctx context.Context, hc httpx.HandlerContext) error {
		postHandler = true
		hc.AbortJSON(200, map[string]string{"status": "created"})
		return nil
	})

	router.PUT("/update", func(ctx context.Context, hc httpx.HandlerContext) error {
		putHandler = true
		hc.AbortJSON(200, map[string]string{"status": "updated"})
		return nil
	})

	router.DELETE("/remove", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, map[string]string{"status": "deleted"})
		return nil
	})

	router.PATCH("/patch", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, map[string]string{"status": "patched"})
		return nil
	})

	// 测试请求
	test := func(method, path string) {
		req := httptest.NewRequest(method, path, nil)
		w := httptest.NewRecorder()
		srv.Engine().ServeHTTP(w, req)
		fmt.Printf("  %s %s -> status: %d\n", method, path, w.Code)
	}

	test("GET", "/hello")
	test("POST", "/data")
	test("PUT", "/update")
	test("DELETE", "/remove")
	test("PATCH", "/patch")

	fmt.Printf("  GET handler called: %v\n", getHandler)
	fmt.Printf("  POST handler called: %v\n", postHandler)
	fmt.Printf("  PUT handler called: %v\n", putHandler)
}

func middlewareExample() {
	srv := gin.NewServer()

	// 添加服务器级中间件
	srv.Use(func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) error {
			fmt.Println("  [middleware] before handler")
			next(ctx, hc)
			fmt.Println("  [middleware] after handler")
			return nil
		}
	})

	router := srv.Router()
	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, map[string]string{"result": "ok"})
		return nil
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	srv.Engine().ServeHTTP(w, req)
	fmt.Printf("  Response status: %d\n", w.Code)
}

func handlerContextExample() {
	srv := gin.NewServer()
	router := srv.Router()

	router.GET("/user/:id", func(ctx context.Context, hc httpx.HandlerContext) error {
		id := hc.Param("id")
		name := hc.Query("name")
		hc.AbortJSON(200, map[string]string{
			"id":   id,
			"name": name,
		})
		return nil
	})

	req := httptest.NewRequest("GET", "/user/123?name=Alice", nil)
	w := httptest.NewRecorder()
	srv.Engine().ServeHTTP(w, req)
	fmt.Printf("  Param + Query -> status: %d, body: %s\n", w.Code, w.Body.String())

	// 测试 AbortWithStatus
	router.GET("/forbidden", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortWithStatus(403)
		return nil
	})

	req = httptest.NewRequest("GET", "/forbidden", nil)
	w = httptest.NewRecorder()
	srv.Engine().ServeHTTP(w, req)
	fmt.Printf("  AbortWithStatus(403) -> status: %d\n", w.Code)

	// 测试 AbortJSON
	router.GET("/error", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(500, map[string]string{"error": "internal error"})
		return nil
	})

	req = httptest.NewRequest("GET", "/error", nil)
	w = httptest.NewRecorder()
	srv.Engine().ServeHTTP(w, req)
	fmt.Printf("  AbortJSON(500) -> status: %d, body: %s\n", w.Code, w.Body.String())
}

func gracefulShutdownExample() {
	srv := gin.NewServer()
	router := srv.Router()

	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		time.Sleep(50 * time.Millisecond)
		hc.AbortJSON(200, map[string]string{"status": "ok"})
		return nil
	})

	go func() {
		srv.StartWithGraceful()
	}()

	time.Sleep(20 * time.Millisecond)
	fmt.Printf("  Server running: %v\n", srv.IsRunning())

	srv.GracefulShutdown(context.Background())
	fmt.Printf("  Server running after shutdown: %v\n", srv.IsRunning())
}