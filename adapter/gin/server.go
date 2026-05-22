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

type GinServer struct {
	engine      *gin.Engine
	router      *ginRouter
	httpSrv     *http.Server
	wsConfig    *httpx.WSConfig
	middlewares []httpx.MiddlewareFunc
	running     bool
	mu          sync.RWMutex
	addr        string
}

func NewServer(opts ...ServerOption) *GinServer {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	srv := &GinServer{
		engine:  engine,
		router:  newGinRouter(engine),
		httpSrv: &http.Server{Addr: "", Handler: engine},
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

type ServerOption func(*GinServer)

func WithWSConfig(cfg *httpx.WSConfig) ServerOption {
	return func(s *GinServer) {
		s.wsConfig = cfg
	}
}

func (s *GinServer) Start(addr string) error {
	s.mu.Lock()
	s.running = true
	s.addr = addr
	s.mu.Unlock()

	s.httpSrv.Addr = addr
	return s.httpSrv.ListenAndServe()
}

func (s *GinServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	return s.httpSrv.Shutdown(ctx)
}

func (s *GinServer) StartWithGraceful(opts ...*httpx.StartOption) error {
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
		var err error
		if opt.TLSConfig != nil {
			err = s.StartTLS(opt.TLSConfig)
		} else {
			err = s.Start(s.addr)
		}
		if err != nil && err != http.ErrServerClosed {
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

func (s *GinServer) StartTLS(cfg *httpx.TLSConfig) error {
	s.mu.Lock()
	s.running = true
	s.mu.Unlock()

	addr := s.addr
	if addr == "" {
		addr = ":0"
	}
	s.httpSrv.Addr = addr
	return s.httpSrv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
}

func (s *GinServer) GracefulShutdown(ctx context.Context) error {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	return s.httpSrv.Shutdown(ctx)
}

func (s *GinServer) Engine() *gin.Engine {
	return s.engine
}

func (s *GinServer) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *GinServer) Router() httpx.Router {
	return s.router
}

func (s *GinServer) Use(middleware ...httpx.MiddlewareFunc) httpx.Server {
	for _, m := range middleware {
		s.middlewares = append(s.middlewares, m)
		s.engine.Use(toGinMiddleware(m))
	}
	return s
}

func toGinMiddleware(m httpx.MiddlewareFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		hc := &ginHandlerContext{GinCtx: c}
		currentCtx := c.Request.Context()

		var next httpx.HandlerFunc
		next = func(ctx context.Context, hc httpx.HandlerContext) {
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		}

		wrapped := m(next)
		wrapped(currentCtx, hc)
	}
}

func (s *GinServer) EnableWS(wsConfig *httpx.WSConfig) httpx.Server {
	s.wsConfig = wsConfig
	return s
}

type ginRouter struct {
	engine *gin.Engine
}

func newGinRouter(engine *gin.Engine) *ginRouter {
	return &ginRouter{engine: engine}
}

type ginRouterGroup struct {
	*ginRouter
	prefix string
}

func (r *ginRouter) GROUP(prefix string, mw ...httpx.MiddlewareFunc) *httpx.RouterGroup {
	return &httpx.RouterGroup{
		Prefix: prefix,
		Router: &ginRouterGroup{ginRouter: r, prefix: prefix},
	}
}

func (g *ginRouterGroup) GET(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	g.engine.GET(g.prefix+path, g.convertHandlers(handlers)...)
	return g
}

func (g *ginRouterGroup) POST(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	g.engine.POST(g.prefix+path, g.convertHandlers(handlers)...)
	return g
}

func (g *ginRouterGroup) PUT(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	g.engine.PUT(g.prefix+path, g.convertHandlers(handlers)...)
	return g
}

func (g *ginRouterGroup) DELETE(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	g.engine.DELETE(g.prefix+path, g.convertHandlers(handlers)...)
	return g
}

func (g *ginRouterGroup) PATCH(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	g.engine.PATCH(g.prefix+path, g.convertHandlers(handlers)...)
	return g
}

func (g *ginRouterGroup) WS(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	g.engine.GET(g.prefix+path, g.convertHandlers(handlers)...)
	return g
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
func (h *ginHandlerContext) AbortWithStatus(code int) { h.GinCtx.AbortWithStatus(code) }
func (h *ginHandlerContext) AbortJSON(code int, body interface{}) { h.GinCtx.AbortWithStatusJSON(code, body) }