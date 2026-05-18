package middleware

import (
	"context"
	"testing"

	"github.com/gospacex/httpx"
)

func TestRecoverMiddleware_DebugMode(t *testing.T) {
	recover := NewRecoverMiddleware("debug")
	handlerCalled := false
	handler := func(ctx context.Context, hc httpx.HandlerContext) {
		handlerCalled = true
	}
	wrapped := recover.Handle(handler)
	wrapped(context.Background(), nil)
	if !handlerCalled {
		t.Error("Handler should be called")
	}
}

func TestRecoverMiddleware_ReleaseMode(t *testing.T) {
	recoverMW := NewRecoverMiddleware("release")
	wrapped := recoverMW.Handle(func(ctx context.Context, hc httpx.HandlerContext) {
		panic("test panic")
	})
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Recover caught unexpected panic: %v", r)
		}
	}()
	wrapped(context.Background(), nil)
}

func TestLoggerMiddleware(t *testing.T) {
	logger := NewLoggerMiddleware()
	handler := func(ctx context.Context, hc httpx.HandlerContext) {}
	wrapped := logger.Handle(handler)
	wrapped(context.Background(), nil)
}

func TestCORSMiddleware(t *testing.T) {
	cors := NewCORSMiddleware()
	handler := func(ctx context.Context, hc httpx.HandlerContext) {}
	wrapped := cors.Handle(handler)
	wrapped(context.Background(), nil)
}

func TestRateLimitMiddleware(t *testing.T) {
	rateLimit := NewRateLimitMiddleware()
	handler := func(ctx context.Context, hc httpx.HandlerContext) {}
	wrapped := rateLimit.Handle(handler)
	wrapped(context.Background(), nil)
}