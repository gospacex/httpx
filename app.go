package httpx

import (
	"fmt"
)

type App struct {
	router      *routerImpl
	middlewares []MiddlewareFunc
	config      *Config
	adapter     AdapterFactory
}

// New creates a new App instance with the configuration from the specified file path.
func New(configPath string) (*App, error) {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	factory, err := getAdapter(cfg.Adapter)
	if err != nil {
		return nil, err
	}

	return &App{
		router:   newRouter(),
		config:   cfg,
		adapter:  factory,
	}, nil
}

func (a *App) GET(path string, handlers ...HandlerFunc) *Route {
	return a.router.addRoute("GET", path, handlers, nil)
}

func (a *App) POST(path string, handlers ...HandlerFunc) *Route {
	return a.router.addRoute("POST", path, handlers, nil)
}

func (a *App) PUT(path string, handlers ...HandlerFunc) *Route {
	return a.router.addRoute("PUT", path, handlers, nil)
}

func (a *App) DELETE(path string, handlers ...HandlerFunc) *Route {
	return a.router.addRoute("DELETE", path, handlers, nil)
}

func (a *App) PATCH(path string, handlers ...HandlerFunc) *Route {
	return a.router.addRoute("PATCH", path, handlers, nil)
}

func (a *App) WS(path string, handlers ...HandlerFunc) *Route {
	return a.router.addRoute("WS", path, handlers, nil)
}

func (a *App) Group(prefix string, mw ...MiddlewareFunc) *RouterGroup {
	return a.router.GROUP(prefix, mw...)
}

func (a *App) Use(mw ...MiddlewareFunc) *App {
	a.middlewares = append(a.middlewares, mw...)
	return a
}

func (a *App) Run() error {
	srv := a.adapter()

	for _, mw := range a.middlewares {
		srv.Use(mw)
	}

	router := srv.Router()
	a.router.setupToRouter(router)

	addr := fmt.Sprintf(":%d", a.config.Port)
	return srv.Start(addr)
}

func (a *App) RunOnAddr(addr string) error {
	srv := a.adapter()

	for _, mw := range a.middlewares {
		srv.Use(mw)
	}

	router := srv.Router()
	a.router.setupToRouter(router)

	return srv.Start(addr)
}

func (a *App) Adapter() Server {
	srv := a.adapter()

	for _, mw := range a.middlewares {
		srv.Use(mw)
	}

	router := srv.Router()
	a.router.setupToRouter(router)

	return srv
}