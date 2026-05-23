package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gospacex/httpx"
	"github.com/gospacex/httpx/adapter/nethttp"
)

func main() {
	fmt.Println("=== httpx nethttp Adapter 示例 ===")
	fmt.Println()

	// 注意: nethttp adapter 需要在同一个包内初始化 router
	// 以下示例展示 API 设计和使用方式

	// 1. 基本服务器
	fmt.Println("1. 基本服务器")
	basicServer()

	// 2. 路由示例
	fmt.Println("\n2. 路由示例")
	routingExample()

	// 3. 优雅关闭
	fmt.Println("\n3. 优雅关闭示例")
	gracefulShutdownExample()

	fmt.Println("\n=== 所有示例完成 ===")
}

func basicServer() {
	srv := nethttp.NewServer()
	fmt.Printf("  Server created, running: %v\n", srv.IsRunning())
}

func routingExample() {
	srv := nethttp.NewServer()

	// nethttp 需要在包内设置 router
	// 这里演示如何创建 router 并注册路由
	router := nethttp.NewRouter()

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

	// 使用 srv.Use() 添加中间件
	srv.Use(func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) {
			fmt.Println("  [middleware] before handler")
			next(ctx, hc)
		}
	})

	fmt.Println("  nethttp router 创建完成（需在包内关联到 server）")
}

func gracefulShutdownExample() {
	srv := nethttp.NewServer()

	router := nethttp.NewRouter()
	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		time.Sleep(50 * time.Millisecond)
		hc.AbortJSON(200, map[string]string{"status": "ok"})
	})

	go func() {
		srv.Start(":0")
	}()

	time.Sleep(20 * time.Millisecond)
	fmt.Printf("  Server running: %v\n", srv.IsRunning())

	srv.Stop(context.Background())
	fmt.Printf("  Server running after stop: %v\n", srv.IsRunning())
}