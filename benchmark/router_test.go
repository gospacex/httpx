package benchmark

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gospacex/httpx"
	"github.com/gospacex/httpx/adapter/gin"
)

func BenchmarkRouter_Get(b *testing.B) {
	srv := gin.NewServer()
	router := srv.Router()
	router.GET("/users/:id", func(ctx context.Context, hc httpx.HandlerContext) {})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/users/1", nil)
		w := httptest.NewRecorder()
		srv.Engine().ServeHTTP(w, req)
	}
}

func BenchmarkRouter_Group(b *testing.B) {
	srv := gin.NewServer()
	router := srv.Router()

	v1 := router.GROUP("/api/v1")
	v1.Router.GET("/users", func(ctx context.Context, hc httpx.HandlerContext) {})
	v1.Router.GET("/orders", func(ctx context.Context, hc httpx.HandlerContext) {})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/v1/users", nil)
		w := httptest.NewRecorder()
		srv.Engine().ServeHTTP(w, req)
	}
}

func BenchmarkRouter_Middleware(b *testing.B) {
	srv := gin.NewServer()
	srv.Use(func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(ctx context.Context, hc httpx.HandlerContext) {
			next(ctx, hc)
		}
	})
	router := srv.Router()
	router.GET("/test", func(ctx context.Context, hc httpx.HandlerContext) {})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		srv.Engine().ServeHTTP(w, req)
	}
}