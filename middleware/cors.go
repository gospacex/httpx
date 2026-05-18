package middleware

import (
	"context"

	"github.com/gospacex/httpx"
)

type CORSMiddleware struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

func NewCORSMiddleware() *CORSMiddleware {
	return &CORSMiddleware{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}
}

func (m *CORSMiddleware) Handle(next httpx.HandlerFunc) httpx.HandlerFunc {
	return func(ctx context.Context, hc httpx.HandlerContext) {
		if ginCtx, ok := hc.(*ginHandlerContext); ok {
			ginCtx.GinCtx.Header("Access-Control-Allow-Origin", "*")
			ginCtx.GinCtx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			ginCtx.GinCtx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		next(ctx, hc)
	}
}