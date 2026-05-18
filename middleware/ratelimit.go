package middleware

import (
	"context"
	"sync"
	"time"

	"github.com/gospacex/httpx"
)

type RateLimitMiddleware struct {
	rate     int
	burst    int
	interval time.Duration
	tokens   map[string]int
	mu       sync.Mutex
}

func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{
		rate:     100,
		burst:    200,
		interval: time.Second,
		tokens:   make(map[string]int),
	}
}

func (m *RateLimitMiddleware) Handle(next httpx.HandlerFunc) httpx.HandlerFunc {
	return func(ctx context.Context, hc httpx.HandlerContext) {
		if ginCtx, ok := hc.(*ginHandlerContext); ok {
			key := ginCtx.GinCtx.ClientIP()
			m.mu.Lock()
			count, ok := m.tokens[key]
			if !ok {
				count = m.rate
			}
			if count > 0 {
				m.tokens[key] = count - 1
				m.mu.Unlock()
			} else {
				m.mu.Unlock()
				ginCtx.GinCtx.AbortWithStatusJSON(429, map[string]interface{}{
					"code":    429,
					"message": "too many requests",
				})
				return
			}
		}
		next(ctx, hc)
	}
}