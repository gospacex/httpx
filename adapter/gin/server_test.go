package gin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gospacex/httpx"
)

func TestGinServer_StartStop(t *testing.T) {
	srv := NewServer()
	errCh := make(chan error, 1)

	go func() {
		errCh <- srv.Start(":0")
	}()

	time.Sleep(50 * time.Millisecond)

	err := srv.Stop(context.Background())
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			t.Fatalf("Start returned unexpected error: %v", err)
		}
	default:
	}

	_ = errCh
}

func TestGinServer_GracefulShutdown(t *testing.T) {
	srv := NewServer()
	router := srv.Router()
	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		time.Sleep(100 * time.Millisecond)
	})

	go func() {
		srv.StartWithGraceful()
	}()

	time.Sleep(100 * time.Millisecond)
	if !srv.IsRunning() {
		t.Error("Server should be running")
	}

	srv.GracefulShutdown(context.Background())

	if srv.IsRunning() {
		t.Error("Server should be stopped")
	}
}

func TestGinServer_Router(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	called := false
	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		called = true
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if !called {
		t.Error("Handler was not called")
	}
}

func TestGinServer_Router_POST(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	called := false
	router.POST("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		called = true
	})

	req := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if !called {
		t.Error("Handler was not called")
	}
}