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

func TestGinServer_Router_GET(t *testing.T) {
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

func TestGinServer_Router_PUT(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	called := false
	router.PUT("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		called = true
	})

	req := httptest.NewRequest("PUT", "/test", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if !called {
		t.Error("Handler was not called")
	}
}

func TestGinServer_Router_DELETE(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	called := false
	router.DELETE("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		called = true
	})

	req := httptest.NewRequest("DELETE", "/test", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if !called {
		t.Error("Handler was not called")
	}
}

func TestGinServer_Router_PATCH(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	called := false
	router.PATCH("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		called = true
	})

	req := httptest.NewRequest("PATCH", "/test", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if !called {
		t.Error("Handler was not called")
	}
}

func TestGinServer_RouterGroup(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	group := router.GROUP("/api")
	rg := group.Router.(*ginRouterGroup)
	rg.GET("/hello", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{"message": "hello"})
	})

	req := httptest.NewRequest("GET", "/api/hello", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGinServer_RouterGroup_GET(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	group := router.GROUP("/api")
	rg := group.Router.(*ginRouterGroup)
	rg.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, nil)
	})

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGinServer_RouterGroup_POST(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	group := router.GROUP("/api")
	rg := group.Router.(*ginRouterGroup)
	rg.POST("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, nil)
	})

	req := httptest.NewRequest("POST", "/api/test", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGinServer_HandlerContext_Query(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	var queryValue string
	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		queryValue = hc.Query("name")
		hc.AbortJSON(200, nil)
	})

	req := httptest.NewRequest("GET", "/test?name=alice", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if queryValue != "alice" {
		t.Errorf("Expected query 'alice', got '%s'", queryValue)
	}
}

func TestGinServer_HandlerContext_Param(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	var paramValue string
	router.GET("/users/:id", func(ctx context.Context, hc httpx.HandlerContext) {
		paramValue = hc.Param("id")
		hc.AbortJSON(200, nil)
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if paramValue != "123" {
		t.Errorf("Expected param '123', got '%s'", paramValue)
	}
}

func TestGinServer_HandlerContext_AbortWithStatus(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortWithStatus(404)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestGinServer_HandlerContext_AbortJSON(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, map[string]string{"key": "value"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() == "" {
		t.Error("Expected body, got empty")
	}
}

func TestGinServer_ServerMiddleware(t *testing.T) {
	srv := NewServer()

	middlewareCalled := false
	srv.Use(func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) {
			middlewareCalled = true
			next(ctx, hc)
		}
	})

	router := srv.Router()
	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, nil)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if !middlewareCalled {
		t.Error("Server middleware was not called")
	}
}

func TestGinServer_WS(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	handlerCalled := false
	router.WS("/ws", func(ctx context.Context, hc httpx.HandlerContext) {
		handlerCalled = true
		hc.AbortJSON(200, nil)
	})

	req := httptest.NewRequest("GET", "/ws", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	_ = req
	if !handlerCalled {
		t.Error("WS handler was not called")
	}
}

func TestGinServer_RouterGroup_WS(t *testing.T) {
	srv := NewServer()
	router := srv.Router()

	group := router.GROUP("/api")
	rg := group.Router.(*ginRouterGroup)
	rg.WS("/ws", func(ctx context.Context, hc httpx.HandlerContext) {
		hc.AbortJSON(200, nil)
	})

	req := httptest.NewRequest("GET", "/api/ws", nil)
	w := httptest.NewRecorder()
	srv.engine.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}