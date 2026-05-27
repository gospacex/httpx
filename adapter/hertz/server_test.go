package hertz

import (
	"context"
	"testing"
	"time"

	"github.com/gospacex/httpx"
)

func TestHertzServer_StartStop(t *testing.T) {
	srv := NewServer()

	if srv.IsRunning() {
		t.Error("Server should not be running initially")
	}
}

func TestHertzServer_GracefulShutdown(t *testing.T) {
	srv := NewServer()

	srv.GracefulShutdown(context.Background())

	if srv.IsRunning() {
		t.Error("Server should be stopped after GracefulShutdown")
	}
}

func TestHertzServer_Router_GET(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, nil)
		return nil
	})
}

func TestHertzServer_Router_POST(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	router.POST("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, nil)
		return nil
	})
}

func TestHertzServer_RouterGroup(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	group := router.GROUP("/api")
	rg := group.Router.(*hertzRouterGroup)
	rg.GET("/hello", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, map[string]string{"message": "hello"})
		return nil
	})
}

func TestHertzServer_RouterGroup_Use(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	group := router.GROUP("/api")
	rg := group.Router.(*hertzRouterGroup)
	rg.Use(func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) error {
			hc.AbortJSON(200, map[string]string{"middleware": "called"})
			return nil
		}
	})
	rg.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, nil)
		return nil
	})
}

func TestHertzServer_RouterGroup_HTTPMethods(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	group := router.GROUP("/api")
	rg := group.Router.(*hertzRouterGroup)

	rg.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, nil)
		return nil
	})
	rg.POST("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, nil)
		return nil
	})
	rg.PUT("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, nil)
		return nil
	})
	rg.DELETE("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, nil)
		return nil
	})
	rg.PATCH("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, nil)
		return nil
	})
}

func TestHertzServer_HandlerContext_Query(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		_ = hc.Query("name")
		hc.AbortJSON(200, nil)
		return nil
	})
}

func TestHertzServer_HandlerContext_Param(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	router.GET("/users/:id", func(ctx context.Context, hc httpx.HandlerContext) error {
		_ = hc.Param("id")
		hc.AbortJSON(200, nil)
		return nil
	})
}

func TestHertzServer_HandlerContext_AbortWithStatus(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortWithStatus(404)
		return nil
	})
}

func TestHertzServer_HandlerContext_AbortJSON(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, map[string]string{"key": "value"})
		return nil
	})
}

func TestHertzServer_Middleware_Chaining(t *testing.T) {
	srv := NewServer()

	srv.Use(func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) error {
			next(ctx, hc)
			return nil
		}
	})

	router := srv.Router()
	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, nil)
		return nil
	})
}

func TestHertzServer_WS(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	router.WS("/ws", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, nil)
		return nil
	})
}

func TestHertzServer_WithHostPorts(t *testing.T) {
	srv := NewServer(WithHostPorts(":8080"))

	if srv.addr != ":8080" {
		t.Errorf("Expected addr ':8080', got '%s'", srv.addr)
	}
}

func TestHertzServer_EnableWS(t *testing.T) {
	srv := NewServer()
	srv.EnableWS(&httpx.WSConfig{})

	if srv.wsConfig == nil {
		t.Error("WSConfig should be set")
	}
}

func TestHertzServer_StartWithGraceful(t *testing.T) {
	srv := NewServer(WithHostPorts(":0"))
	router := srv.Router()
	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) error {
		hc.AbortJSON(200, nil)
		return nil
	})

	go func() {
		srv.StartWithGraceful()
	}()

	time.Sleep(50 * time.Millisecond)
	if !srv.IsRunning() {
		t.Error("Server should be running after StartWithGraceful")
	}

	srv.GracefulShutdown(context.Background())

	if srv.IsRunning() {
		t.Error("Server should be stopped after GracefulShutdown")
	}
}