package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gospacex/httpx"
	"github.com/gospacex/httpx/adapter/hertz"
)

func main() {
	fmt.Println("=== httpx Hertz Adapter 示例 ===")
	fmt.Println()

	// 1. 基本服务器
	fmt.Println("1. 基本服务器")
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
	srv := hertz.NewServer()
	fmt.Printf("  Server created, running: %v\n", srv.IsRunning())
}

func routingExample() {
	srv := hertz.NewServer()
	router := srv.Router()

	router.GET("/hello", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{"message": "hello"})
	})

	router.POST("/data", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{"status": "created"})
	})

	router.PUT("/update", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{"status": "updated"})
	})

	router.DELETE("/remove", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{"status": "deleted"})
	})

	router.PATCH("/patch", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{"status": "patched"})
	})

	fmt.Println("  Hertz router registered with GET/POST/PUT/DELETE/PATCH")
}

func middlewareExample() {
	srv := hertz.NewServer()

	srv.Use(func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) {
			fmt.Println("  [middleware] before handler")
			next(ctx, hc)
			fmt.Println("  [middleware] after handler")
		}
	})

	router := srv.Router()
	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{"result": "ok"})
	})

	fmt.Println("  Middleware registered")
}

func handlerContextExample() {
	srv := hertz.NewServer()
	router := srv.Router()

	router.GET("/user/:id", func(ctx context.Context, hc httpx.HandlerContext) {
		id := hc.Param("id")
		name := hc.Query("name")
		hc.AbortJSON(200, map[string]string{
			"id":   id,
			"name": name,
		})
	})

	router.GET("/forbidden", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortWithStatus(403)
	})

	router.GET("/error", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(500, map[string]string{"error": "internal error"})
	})

	fmt.Println("  Routes registered: /user/:id, /forbidden, /error")
}

func gracefulShutdownExample() {
	srv := hertz.NewServer(hertz.WithHostPorts(":0"))
	router := srv.Router()

	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		time.Sleep(50 * time.Millisecond)
		hc.AbortJSON(200, map[string]string{"status": "ok"})
	})

	go func() {
		srv.StartWithGraceful()
	}()

	time.Sleep(20 * time.Millisecond)
	fmt.Printf("  Server running: %v\n", srv.IsRunning())

	srv.GracefulShutdown(context.Background())
	fmt.Printf("  Server running after shutdown: %v\n", srv.IsRunning())
}