package httpx

type Router interface {
	GET(path string, handlers ...HandlerFunc) Router
	POST(path string, handlers ...HandlerFunc) Router
	PUT(path string, handlers ...HandlerFunc) Router
	DELETE(path string, handlers ...HandlerFunc) Router
	PATCH(path string, handlers ...HandlerFunc) Router
	GROUP(prefix string, mw ...MiddlewareFunc) *RouterGroup
	WS(path string, handlers ...HandlerFunc) Router
}

type RouterGroup struct {
	Router      Router
	Prefix      string
	middlewares []MiddlewareFunc
}

type routerImpl struct {
	routes []routeRecord
	groups []*RouterGroup
}

type routeRecord struct {
	method      string
	path        string
	handlers    []HandlerFunc
	middlewares []MiddlewareFunc
}

func newRouter() *routerImpl {
	return &routerImpl{}
}

func (r *routerImpl) addRoute(method, path string, handlers []HandlerFunc, mw []MiddlewareFunc) *Route {
	r.routes = append(r.routes, routeRecord{
		method:      method,
		path:        path,
		handlers:    handlers,
		middlewares: mw,
	})
	return &Route{
		path:       path,
		method:     method,
		handlers:   handlers,
		middlewares: mw,
	}
}

func (r *routerImpl) GET(path string, handlers ...HandlerFunc) Router {
	r.addRoute("GET", path, handlers, nil)
	return r
}

func (r *routerImpl) POST(path string, handlers ...HandlerFunc) Router {
	r.addRoute("POST", path, handlers, nil)
	return r
}

func (r *routerImpl) PUT(path string, handlers ...HandlerFunc) Router {
	r.addRoute("PUT", path, handlers, nil)
	return r
}

func (r *routerImpl) DELETE(path string, handlers ...HandlerFunc) Router {
	r.addRoute("DELETE", path, handlers, nil)
	return r
}

func (r *routerImpl) PATCH(path string, handlers ...HandlerFunc) Router {
	r.addRoute("PATCH", path, handlers, nil)
	return r
}

func (r *routerImpl) WS(path string, handlers ...HandlerFunc) Router {
	r.addRoute("WS", path, handlers, nil)
	return r
}

func (r *routerImpl) GROUP(prefix string, mw ...MiddlewareFunc) *RouterGroup {
	g := &RouterGroup{
		Prefix:      prefix,
		Router:      &groupRouter{impl: r, prefix: prefix, middlewares: mw},
		middlewares: mw,
	}
	r.groups = append(r.groups, g)
	return g
}

func (r *routerImpl) setupToRouter(router Router) {
	for _, rec := range r.routes {
		wrapped := r.wrapWithMiddlewares(rec.handlers, rec.middlewares)
		switch rec.method {
		case "GET":
			router.GET(rec.path, wrapped...)
		case "POST":
			router.POST(rec.path, wrapped...)
		case "PUT":
			router.PUT(rec.path, wrapped...)
		case "DELETE":
			router.DELETE(rec.path, wrapped...)
		case "PATCH":
			router.PATCH(rec.path, wrapped...)
		case "WS":
			router.WS(rec.path, wrapped...)
		}
	}

	for _, g := range r.groups {
		grouter := g.Router.(*groupRouter)
		for _, rec := range grouter.routes {
			allMW := make([]MiddlewareFunc, 0, len(g.middlewares)+len(rec.middlewares))
			allMW = append(allMW, g.middlewares...)
			allMW = append(allMW, rec.middlewares...)
			wrapped := r.wrapWithMiddlewares(rec.handlers, allMW)

			fullPath := grouter.prefix + rec.path
			switch rec.method {
			case "GET":
				router.GET(fullPath, wrapped...)
			case "POST":
				router.POST(fullPath, wrapped...)
			case "PUT":
				router.PUT(fullPath, wrapped...)
			case "DELETE":
				router.DELETE(fullPath, wrapped...)
			case "PATCH":
				router.PATCH(fullPath, wrapped...)
			case "WS":
				router.WS(fullPath, wrapped...)
			}
		}
	}
}

func (r *routerImpl) wrapWithMiddlewares(handlers []HandlerFunc, middlewares []MiddlewareFunc) []HandlerFunc {
	if len(middlewares) == 0 {
		return handlers
	}
	if len(handlers) == 0 {
		return nil
	}

	// Chain middlewares: mw1(mw2(mw3(...(finalHandler))))
	result := handlers[len(handlers)-1]
	for i := len(middlewares) - 1; i >= 0; i-- {
		mw := middlewares[i]
		next := result
		result = mw(next)
	}

	// Return wrapped handlers (only the first handler is used since we chained all middlewares to it)
	wrapped := make([]HandlerFunc, len(handlers))
	wrapped[0] = result
	for i := 1; i < len(handlers); i++ {
		wrapped[i] = handlers[i]
	}
	return wrapped
}

type groupRouter struct {
	impl        *routerImpl
	prefix      string
	middlewares []MiddlewareFunc
	routes      []routeRecord
}

func (g *groupRouter) addRoute(method, path string, handlers []HandlerFunc, mw []MiddlewareFunc) *Route {
	g.routes = append(g.routes, routeRecord{
		method:      method,
		path:        path,
		handlers:    handlers,
		middlewares: mw,
	})
	return &Route{
		path:       path,
		method:     method,
		handlers:   handlers,
		middlewares: mw,
	}
}

func (g *groupRouter) GET(path string, handlers ...HandlerFunc) Router {
	g.addRoute("GET", path, handlers, g.middlewares)
	return g
}

func (g *groupRouter) POST(path string, handlers ...HandlerFunc) Router {
	g.addRoute("POST", path, handlers, g.middlewares)
	return g
}

func (g *groupRouter) PUT(path string, handlers ...HandlerFunc) Router {
	g.addRoute("PUT", path, handlers, g.middlewares)
	return g
}

func (g *groupRouter) DELETE(path string, handlers ...HandlerFunc) Router {
	g.addRoute("DELETE", path, handlers, g.middlewares)
	return g
}

func (g *groupRouter) PATCH(path string, handlers ...HandlerFunc) Router {
	g.addRoute("PATCH", path, handlers, g.middlewares)
	return g
}

func (g *groupRouter) WS(path string, handlers ...HandlerFunc) Router {
	g.addRoute("WS", path, handlers, g.middlewares)
	return g
}

func (g *groupRouter) GROUP(prefix string, mw ...MiddlewareFunc) *RouterGroup {
	return &RouterGroup{
		Prefix:      prefix,
		Router:      g,
		middlewares: mw,
	}
}