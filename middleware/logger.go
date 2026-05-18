package middleware

import (
	"context"
	"log"
	"time"

	"github.com/gospacex/httpx"
)

type LoggerMiddleware struct {
	format string
}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{
		format: "[${time}] ${status} ${method} ${path} ${latency}",
	}
}

func (m *LoggerMiddleware) Handle(next httpx.HandlerFunc) httpx.HandlerFunc {
	return func(ctx context.Context, hc httpx.HandlerContext) {
		start := time.Now()
		next(ctx, hc)
		latency := time.Since(start)

		if ginCtx, ok := hc.(*ginHandlerContext); ok {
			log.Printf(m.format,
				"time", start.Format(time.RFC3339),
				"status", ginCtx.GinCtx.Writer.Status(),
				"method", ginCtx.GinCtx.Request.Method,
				"path", ginCtx.GinCtx.Request.URL.Path,
				"latency", latency.String(),
			)
		}
	}
}