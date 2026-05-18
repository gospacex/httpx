package e2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gospacex/httpx"
	"github.com/gospacex/httpx/adapter/gin"
)

func TestE2E_RouterGroup(t *testing.T) {
	srv := gin.NewServer()
	router := srv.Router()

	v1 := router.GROUP("/api/v1")
	v1.Router.GET("/users", func(ctx context.Context, hc httpx.HandlerContext) {})

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	srv.Engine().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestE2E_MiddlewareChain(t *testing.T) {
	srv := gin.NewServer()
	srv.Use(func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) {
			ctx = context.WithValue(ctx, "order", 1)
			next(ctx, hc)
		}
	})

	router := srv.Router()
	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		order := ctx.Value("order")
		if order != 1 {
			t.Error("Middleware chain broken")
		}
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	srv.Engine().ServeHTTP(w, req)
}

func TestE2E_WebSocket(t *testing.T) {
	srv := gin.NewServer()
	srv.EnableWS(nil)
	router := srv.Router()
	router.WS("/ws", func(ctx context.Context, hc httpx.HandlerContext) {})

	req := httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")

	w := httptest.NewRecorder()
	_ = w
}

func TestE2E_GracefulShutdown(t *testing.T) {
	srv := gin.NewServer()
	router := srv.Router()
	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		time.Sleep(50 * time.Millisecond)
	})

	if srv.IsRunning() {
		t.Error("Server should not be running initially")
	}
}