package wire

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gospacex/httpx"
	"github.com/gospacex/httpx/middleware"
)

type GlobalMiddleware struct {
	Logger    *middleware.LoggerMiddleware
	Recover   *middleware.RecoverMiddleware
	CORS      *middleware.CORSMiddleware
	RateLimit *middleware.RateLimitMiddleware
}

func NewGinServer() *ginServer {
	gin.SetMode(gin.ReleaseMode)
	return &ginServer{
		engine: gin.New(),
	}
}

type ginServer struct {
	engine *gin.Engine
}

func (s *ginServer) Start(addr string) error {
	return nil
}

func (s *ginServer) Stop(ctx context.Context) error {
	return nil
}

func (s *ginServer) Router() httpx.Router {
	return nil
}

func (s *ginServer) Use(middleware ...httpx.MiddlewareFunc) httpx.Server {
	return s
}

func (s *ginServer) EnableWS(wsConfig *httpx.WSConfig) httpx.Server {
	return s
}

func GlobalMiddlewareProvider() *GlobalMiddleware {
	return &GlobalMiddleware{
		Logger:    middleware.NewLoggerMiddleware(),
		Recover:   middleware.NewRecoverMiddleware("release"),
		CORS:      middleware.NewCORSMiddleware(),
		RateLimit: middleware.NewRateLimitMiddleware(),
	}
}

type ServerFactory struct {
	server     httpx.Server
	middleware *GlobalMiddleware
}

func NewServerFactory(srv httpx.Server, mw *GlobalMiddleware) *ServerFactory {
	return &ServerFactory{
		server:     srv,
		middleware: mw,
	}
}

func (f *ServerFactory) CreateServer(opts ...ServerOption) httpx.Server {
	f.server.Use(f.middleware.Logger.Handle)
	f.server.Use(f.middleware.Recover.Handle)
	return f.server
}

type ServerOption func(httpx.Server)