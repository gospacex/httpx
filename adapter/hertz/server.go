package hertz

import (
	"context"
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
}

func NewServer(opts ...ServerOption) *hertzServer {
	srv := server.Default()
	return &hertzServer{
		srv:    srv,
		router: &hertzRouter{hertz: srv},
	}
}

type ServerOption func(*hertzServer)

func (s *hertzServer) Start(addr string) error {
	s.mu.Lock()
	s.running = true
	s.mu.Unlock()

	return s.srv.Run()
}

func (s *hertzServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	return nil
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

	go s.srv.Run()

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

	s.srv.Shutdown(ctx)
	return nil
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

func (r *hertzRouter) GROUP(prefix string, mw ...httpx.MiddlewareFunc) *httpx.RouterGroup {
	return &httpx.RouterGroup{
		Router: &hertzRouter{hertz: r.hertz},
		Prefix: prefix,
	}
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