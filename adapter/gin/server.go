package gin

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/gospacex/httpx"
)

type ginServer struct {
	engine      *gin.Engine
	router      *ginRouter
	httpSrv     *http.Server
	wsConfig    *httpx.WSConfig
	middlewares []httpx.MiddlewareFunc
	running     bool
	mu          sync.RWMutex
}

func NewServer(opts ...ServerOption) *ginServer {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	srv := &ginServer{
		engine:  engine,
		router:  newGinRouter(engine),
		httpSrv: &http.Server{Addr: "", Handler: engine},
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

type ServerOption func(*ginServer)

func WithWSConfig(cfg *httpx.WSConfig) ServerOption {
	return func(s *ginServer) {
		s.wsConfig = cfg
	}
}

func (s *ginServer) Start(addr string) error {
	s.mu.Lock()
	s.running = true
	s.mu.Unlock()

	s.httpSrv.Addr = addr
	return s.httpSrv.ListenAndServe()
}

func (s *ginServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	return s.httpSrv.Shutdown(ctx)
}

func (s *ginServer) StartWithGraceful(opts ...*httpx.StartOption) error {
	opt := httpx.DefaultStartOption()
	for _, o := range opts {
		if o != nil {
			opt = o
		}
	}

	quit := opt.QuitChan
	if quit == nil {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		quit = ch
	}

	go func() {
		if err := s.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-quit
	log.Println("Shutdown Server ...")

	if opt.BeforeShutdown != nil {
		opt.BeforeShutdown()
	}

	ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout)
	defer cancel()

	err := s.GracefulShutdown(ctx)

	if opt.AfterShutdown != nil {
		opt.AfterShutdown(err)
	}

	return err
}

func (s *ginServer) GracefulShutdown(ctx context.Context) error {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	return s.httpSrv.Shutdown(ctx)
}

func (s *ginServer) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *ginServer) Router() httpx.Router {
	return s.router
}

func (s *ginServer) Use(middleware ...httpx.MiddlewareFunc) httpx.Server {
	s.middlewares = append(s.middlewares, middleware...)
	return s
}

func (s *ginServer) EnableWS(wsConfig *httpx.WSConfig) httpx.Server {
	s.wsConfig = wsConfig
	return s
}

type ginRouter struct {
	engine *gin.Engine
}

func newGinRouter(engine *gin.Engine) *ginRouter {
	return &ginRouter{engine: engine}
}

func (r *ginRouter) GET(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.engine.GET(path, r.convertHandlers(handlers)...)
	return r
}

func (r *ginRouter) POST(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.engine.POST(path, r.convertHandlers(handlers)...)
	return r
}

func (r *ginRouter) PUT(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.engine.PUT(path, r.convertHandlers(handlers)...)
	return r
}

func (r *ginRouter) DELETE(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.engine.DELETE(path, r.convertHandlers(handlers)...)
	return r
}

func (r *ginRouter) PATCH(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.engine.PATCH(path, r.convertHandlers(handlers)...)
	return r
}

func (r *ginRouter) GROUP(prefix string, mw ...httpx.MiddlewareFunc) *httpx.RouterGroup {
	return &httpx.RouterGroup{
		Prefix: prefix,
		Router: &ginRouter{engine: r.engine},
	}
}

func (r *ginRouter) WS(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.engine.GET(path, r.convertHandlers(handlers)...)
	return r
}

func (r *ginRouter) convertHandlers(handlers []httpx.HandlerFunc) []gin.HandlerFunc {
	result := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		handler := h
		result[i] = func(c *gin.Context) {
			hc := &ginHandlerContext{GinCtx: c}
			handler(c.Request.Context(), hc)
		}
	}
	return result
}

type ginHandlerContext struct {
	GinCtx *gin.Context
}

func (h *ginHandlerContext) Request() interface{}   { return h.GinCtx.Request }
func (h *ginHandlerContext) Response() interface{}  { return h.GinCtx.Writer }
func (h *ginHandlerContext) Param(key string) string { return h.GinCtx.Param(key) }
func (h *ginHandlerContext) Query(key string) string { return h.GinCtx.Query(key) }
func (h *ginHandlerContext) Bind(into interface{}) error { return h.GinCtx.Bind(into) }