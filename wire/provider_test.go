package wire

import (
	"testing"

	"github.com/gospacex/httpx/middleware"
)

func TestGlobalMiddlewareProvider(t *testing.T) {
	mw := GlobalMiddlewareProvider()
	if mw.Logger == nil {
		t.Error("Logger middleware should not be nil")
	}
	if mw.Recover == nil {
		t.Error("Recover middleware should not be nil")
	}
	if mw.CORS == nil {
		t.Error("CORS middleware should not be nil")
	}
	if mw.RateLimit == nil {
		t.Error("RateLimit middleware should not be nil")
	}
}

func TestServerFactory(t *testing.T) {
	mw := &GlobalMiddleware{
		Logger:    middleware.NewLoggerMiddleware(),
		Recover:   middleware.NewRecoverMiddleware("release"),
		CORS:      middleware.NewCORSMiddleware(),
		RateLimit: middleware.NewRateLimitMiddleware(),
	}
	srv := NewGinServer()
	factory := NewServerFactory(srv, mw)
	if factory == nil {
		t.Error("ServerFactory should not be nil")
	}
	if factory.server == nil {
		t.Error("Server should not be nil")
	}
	if factory.middleware == nil {
		t.Error("Middleware should not be nil")
	}
}

func TestNewGinServer(t *testing.T) {
	srv := NewGinServer()
	if srv == nil {
		t.Error("GinServer should not be nil")
	}
	if srv.engine == nil {
		t.Error("Engine should not be nil")
	}
}