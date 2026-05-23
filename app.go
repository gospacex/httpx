package httpx

import (
	"fmt"
)

type App struct {
	router      *routerImpl
	middlewares []MiddlewareFunc
}

func newApp() *App {
	return &App{
		router: newRouter(),
	}
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

func (a *App) Run(configPath string) error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	factory, err := getAdapter(cfg.Adapter)
	if err != nil {
		return err
	}

	srv := factory()

	for _, mw := range a.middlewares {
		srv.Use(mw)
	}

	router := srv.Router()
	a.router.setupToRouter(router)

	addr := fmt.Sprintf(":%d", cfg.Port)
	return srv.Start(addr)
}