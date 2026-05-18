package middleware

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/gospacex/httpx"
)

type ginHandlerContext struct {
	GinCtx *gin.Context
}

func (h *ginHandlerContext) Request() interface{}   { return h.GinCtx.Request }
func (h *ginHandlerContext) Response() interface{}  { return h.GinCtx.Writer }
func (h *ginHandlerContext) Param(key string) string { return h.GinCtx.Param(key) }
func (h *ginHandlerContext) Query(key string) string { return h.GinCtx.Query(key) }
func (h *ginHandlerContext) Bind(into interface{}) error { return h.GinCtx.Bind(into) }

type RecoverMiddleware struct {
	mode string
}

func NewRecoverMiddleware(mode string) *RecoverMiddleware {
	return &RecoverMiddleware{mode: mode}
}

func (m *RecoverMiddleware) Handle(next httpx.HandlerFunc) httpx.HandlerFunc {
	return func(ctx context.Context, hc httpx.HandlerContext) {
		defer func() {
			if r := recover(); r != nil {
				if m.mode == "debug" {
					panic(r)
				}
				log.Printf("panic recovered: %v\n%s", r, debug.Stack())
				if ginCtx, ok := hc.(*ginHandlerContext); ok {
					ginCtx.GinCtx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
						"code":    500,
						"message": "internal server error",
					})
				}
			}
		}()
		next(ctx, hc)
	}
}