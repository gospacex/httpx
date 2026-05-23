package nethttp

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gospacex/httpx"
)

type netHttpServer struct {
	httpSrv     *http.Server
	wsConfig    *httpx.WSConfig
	middlewares []httpx.MiddlewareFunc
	running     bool
	mu          sync.RWMutex
	router      *netHttpRouter
	addr        string
}

func NewServer(opts ...ServerOption) *netHttpServer {
	srv := &netHttpServer{
		httpSrv: &http.Server{},
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

type ServerOption func(*netHttpServer)

func (s *netHttpServer) Start(addr string) error {
	s.mu.Lock()
	s.running = true
	s.addr = addr
	s.mu.Unlock()

	s.httpSrv.Addr = addr
	return s.httpSrv.ListenAndServe()
}

func (s *netHttpServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	return s.httpSrv.Shutdown(ctx)
}

func (s *netHttpServer) StartWithGraceful(opts ...*httpx.StartOption) error {
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

func (s *netHttpServer) StartTLS(cfg *httpx.TLSConfig) error {
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

func (s *netHttpServer) GracefulShutdown(ctx context.Context) error {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	return s.httpSrv.Shutdown(ctx)
}

func (s *netHttpServer) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *netHttpServer) Router() httpx.Router {
	return s.router
}

func (s *netHttpServer) Use(middleware ...httpx.MiddlewareFunc) httpx.Server {
	s.middlewares = append(s.middlewares, middleware...)
	return s
}

func (s *netHttpServer) EnableWS(wsConfig *httpx.WSConfig) httpx.Server {
	s.wsConfig = wsConfig
	return s
}

type netHttpRouter struct {
	mux *http.ServeMux
}

func NewRouter() *netHttpRouter {
	return &netHttpRouter{mux: http.NewServeMux()}
}

func (r *netHttpRouter) GET(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		hc := &netHttpHandlerContext{req: req, resp: w}
		for _, h := range handlers {
			h(req.Context(), hc)
		}
	})
	return r
}

func (r *netHttpRouter) POST(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	return r.GET(path, handlers...)
}

func (r *netHttpRouter) PUT(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	return r.GET(path, handlers...)
}

func (r *netHttpRouter) DELETE(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	return r.GET(path, handlers...)
}

func (r *netHttpRouter) PATCH(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	return r.GET(path, handlers...)
}

type netHttpRouterGroup struct {
	*netHttpRouter
	prefix      string
	middlewares []httpx.MiddlewareFunc
}

func (r *netHttpRouter) GROUP(prefix string, mw ...httpx.MiddlewareFunc) *httpx.RouterGroup {
	return &httpx.RouterGroup{
		Prefix: prefix,
		Router: &netHttpRouterGroup{netHttpRouter: r, prefix: prefix, middlewares: mw},
	}
}

func (g *netHttpRouterGroup) wrapHandlers(handlers []httpx.HandlerFunc) []httpx.HandlerFunc {
	if len(handlers) == 0 {
		return nil
	}
	var wrapped httpx.HandlerFunc = handlers[len(handlers)-1]
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		wrapped = g.middlewares[i](wrapped)
	}
	return []httpx.HandlerFunc{wrapped}
}

func (g *netHttpRouterGroup) GET(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	wrapped := g.wrapHandlers(handlers)
	g.mux.HandleFunc(g.prefix+path, func(w http.ResponseWriter, req *http.Request) {
		hc := &netHttpHandlerContext{req: req, resp: w}
		for _, h := range wrapped {
			h(req.Context(), hc)
			if hc.aborted {
				return
			}
		}
	})
	return g
}

func (g *netHttpRouterGroup) POST(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	return g.GET(path, handlers...)
}

func (g *netHttpRouterGroup) PUT(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	return g.GET(path, handlers...)
}

func (g *netHttpRouterGroup) DELETE(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	return g.GET(path, handlers...)
}

func (g *netHttpRouterGroup) PATCH(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	return g.GET(path, handlers...)
}

func (g *netHttpRouterGroup) WS(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	return g.GET(path, handlers...)
}

func (g *netHttpRouterGroup) Use(middlewares ...httpx.MiddlewareFunc) *httpx.RouterGroup {
	g.middlewares = append(g.middlewares, middlewares...)
	return &httpx.RouterGroup{
		Prefix: g.prefix,
		Router: g,
	}
}

func (r *netHttpRouter) WS(path string, handlers ...httpx.HandlerFunc) httpx.Router {
	return r.GET(path, handlers...)
}

type netHttpHandlerContext struct {
	req     *http.Request
	resp    http.ResponseWriter
	aborted bool
}

func (h *netHttpHandlerContext) Request() interface{} { return h.req }
func (h *netHttpHandlerContext) Response() interface{} { return h.resp }
func (h *netHttpHandlerContext) Param(key string) string { return "" }
func (h *netHttpHandlerContext) Query(key string) string { return h.req.URL.Query().Get(key) }
func (h *netHttpHandlerContext) Bind(into interface{}) error { return nil }

func (h *netHttpHandlerContext) AbortWithStatus(code int) {
	h.resp.WriteHeader(code)
	h.aborted = true
}

func (h *netHttpHandlerContext) AbortJSON(code int, body interface{}) {
	h.resp.Header().Set("Content-Type", "application/json")
	h.resp.WriteHeader(code)
	enc := json.NewEncoder(h.resp)
	enc.Encode(body)
	h.aborted = true
}

func init() {
	httpx.RegisterAdapter("nethttp", func() httpx.Server {
		return NewServer()
	})
}