package main

import (
	"fmt"

	"github.com/gospacex/httpx"
	_ "github.com/gospacex/httpx/adapter/gin"   // 副作用：注册 gin adapter
	_ "github.com/gospacex/httpx/adapter/hertz" // 副作用：注册 hertz adapter
)

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║   httpx HTTP 服务演示                    ║")
	fmt.Println("║   启动后访问 http://localhost:8080       ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	//app, err := httpx.New("/Users/hyx/work/gowork/src/gospacex/httpx/examples/demo/config_hertz.yaml")
	app, err := httpx.New("")
	if err != nil {
		panic(err)
	}

	setupRoutes(app)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
