package hertz

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/gospacex/httpx"
)

type hertzServer struct {
	srv      *server.Hertz
	wsConfig *httpx.WSConfig
	running  bool
	mu       sync.RWMutex
	router   *hertzRouter
	addr     string
}

func NewServer(opts ...ServerOption) *hertzServer {
	srv := server.Default()
	h := &hertzServer{
		srv:    srv,
		router: &hertzRouter{hertz: srv},
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

type ServerOption func(*hertzServer)

func WithHostPorts(addr string) ServerOption {
	return func(s *hertzServer) {
		s.addr = addr
		s.srv = server.New(server.WithHostPorts(addr))
		s.router = &hertzRouter{hertz: s.srv}
	}
}

func WithTLS(cfg *httpx.TLSConfig) ServerOption {
	return func(s *hertzServer) {
		s.srv = server.New(
			server.WithHostPorts(s.addr),
			server.WithTLS(loadTLSConfig(cfg)),
		)
		s.router = &hertzRouter{hertz: s.srv}
	}
}

func loadTLSConfig(cfg *httpx.TLSConfig) *tls.Config {
	tlsCfg, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return &tls.Config{}
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCfg},
	}
}

func (s *hertzServer) Start(addr string) error {
	s.mu.Lock()
	s.running = true
	if s.addr == "" {
		s.addr = addr
	}
	s.mu.Unlock()

	return s.srv.Run()
}

func (s *hertzServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	return s.srv.Shutdown(ctx)
}

func (s *hertzServer) StartWithGraceful(opts ...*httpx.StartOption) error {
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
		if err := s.Start(s.addr); err != nil {
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

func (s *hertzServer) GracefulShutdown(ctx context.Context) error {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	return s.srv.Shutdown(ctx)
}

func (s *hertzServer) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *hertzServer) Router() httpx.Router {
	return s.router
}

func (s *hertzServer) Use(middleware ...httpx.MiddlewareFunc) httpx.Server {
	return s
}

func (s *hertzServer) EnableWS(wsConfig *httpx.WSConfig) httpx.Server {
	s.wsConfig = wsConfig
	return s
}

type hertzRouter struct {
	hertz *server.Hertz
}

func (r *hertzRouter) GET(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.hertz.GET(path, func(ctx context.Context, c *app.RequestContext) {
		hc := &hertzHandlerContext{RequestContext: c}
		for _, h := range handlers {
			h(ctx, hc)
		}
	})
	return r
}

func (r *hertzRouter) POST(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.hertz.POST(path, func(ctx context.Context, c *app.RequestContext) {
		hc := &hertzHandlerContext{RequestContext: c}
		for _, h := range handlers {
			h(ctx, hc)
		}
	})
	return r
}

func (r *hertzRouter) PUT(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.hertz.PUT(path, func(ctx context.Context, c *app.RequestContext) {
		hc := &hertzHandlerContext{RequestContext: c}
		for _, h := range handlers {
			h(ctx, hc)
		}
	})
	return r
}

func (r *hertzRouter) DELETE(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.hertz.DELETE(path, func(ctx context.Context, c *app.RequestContext) {
		hc := &hertzHandlerContext{RequestContext: c}
		for _, h := range handlers {
			h(ctx, hc)
		}
	})
	return r
}

func (r *hertzRouter) PATCH(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.hertz.PATCH(path, func(ctx context.Context, c *app.RequestContext) {
		hc := &hertzHandlerContext{RequestContext: c}
		for _, h := range handlers {
			h(ctx, hc)
		}
	})
	return r
}

type hertzRouterGroup struct {
	*hertzRouter
	prefix      string
	middlewares []httpx.MiddlewareFunc
}

func (r *hertzRouter) GROUP(prefix string, mw ...httpx.MiddlewareFunc) *httpx.RouterGroup {
	return &httpx.RouterGroup{
		Prefix: prefix,
		Router: &hertzRouterGroup{hertzRouter: r, prefix: prefix, middlewares: mw},
	}
}

func (g *hertzRouterGroup) wrapHandlers(handlers []httpx.HandlerFunc) []httpx.HandlerFunc {
	if len(handlers) == 0 {
		return nil
	}
	var wrapped httpx.HandlerFunc = handlers[len(handlers)-1]
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		wrapped = g.middlewares[i](wrapped)
	}
	return []httpx.HandlerFunc{wrapped}
}

func (g *hertzRouterGroup) GET(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	wrapped := g.wrapHandlers(handlers)
	g.hertz.GET(g.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		hc := &hertzHandlerContext{RequestContext: c}
		for _, h := range wrapped {
			h(ctx, hc)
		}
	})
	return g
}

func (g *hertzRouterGroup) POST(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	wrapped := g.wrapHandlers(handlers)
	g.hertz.POST(g.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		hc := &hertzHandlerContext{RequestContext: c}
		for _, h := range wrapped {
			h(ctx, hc)
		}
	})
	return g
}

func (g *hertzRouterGroup) PUT(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	wrapped := g.wrapHandlers(handlers)
	g.hertz.PUT(g.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		hc := &hertzHandlerContext{RequestContext: c}
		for _, h := range wrapped {
			h(ctx, hc)
		}
	})
	return g
}

func (g *hertzRouterGroup) DELETE(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	wrapped := g.wrapHandlers(handlers)
	g.hertz.DELETE(g.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		hc := &hertzHandlerContext{RequestContext: c}
		for _, h := range wrapped {
			h(ctx, hc)
		}
	})
	return g
}

func (g *hertzRouterGroup) PATCH(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	wrapped := g.wrapHandlers(handlers)
	g.hertz.PATCH(g.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		hc := &hertzHandlerContext{RequestContext: c}
		for _, h := range wrapped {
			h(ctx, hc)
		}
	})
	return g
}

func (g *hertzRouterGroup) WS(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	wrapped := g.wrapHandlers(handlers)
	g.hertz.GET(g.prefix+path, func(ctx context.Context, c *app.RequestContext) {
		hc := &hertzHandlerContext{RequestContext: c}
		for _, h := range wrapped {
			h(ctx, hc)
		}
	})
	return g
}

func (r *hertzRouter) WS(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.hertz.GET(path, func(ctx context.Context, c *app.RequestContext) {
		hc := &hertzHandlerContext{RequestContext: c}
		for _, h := range handlers {
			h(ctx, hc)
		}
	})
	return r
}

type hertzHandlerContext struct {
	RequestContext *app.RequestContext
}

func (h *hertzHandlerContext) Request() interface{} { return h.RequestContext.Request }
func (h *hertzHandlerContext) Response() interface{} { return h.RequestContext.Response }
func (h *hertzHandlerContext) Param(key string) string { return h.RequestContext.Param(key) }
func (h *hertzHandlerContext) Query(key string) string { return h.RequestContext.Query(key) }
func (h *hertzHandlerContext) Bind(into interface{}) error { return h.RequestContext.Bind(into) }
func (h *hertzHandlerContext) AbortWithStatus(code int) { h.RequestContext.AbortWithStatus(code) }
func (h *hertzHandlerContext) AbortJSON(code int, body interface{}) { h.RequestContext.JSON(code, body); h.RequestContext.Abort() }