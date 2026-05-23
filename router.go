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

func (r *routerImpl) GET(path string, handlers ...HandlerFunc) *Route {
	return r.addRoute("GET", path, handlers)
}

func (r *routerImpl) POST(path string, handlers ...HandlerFunc) *Route {
	return r.addRoute("POST", path, handlers)
}

func (r *routerImpl) PUT(path string, handlers ...HandlerFunc) *Route {
	return r.addRoute("PUT", path, handlers)
}

func (r *routerImpl) DELETE(path string, handlers ...HandlerFunc) *Route {
	return r.addRoute("DELETE", path, handlers)
}

func (r *routerImpl) PATCH(path string, handlers ...HandlerFunc) *Route {
	return r.addRoute("PATCH", path, handlers)
}

func (r *routerImpl) WS(path string, handlers ...HandlerFunc) *Route {
	return r.addRoute("WS", path, handlers)
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

func (r *routerImpl) addRoute(method, path string, handlers []HandlerFunc) *Route {
	r.routes = append(r.routes, routeRecord{
		method:      method,
		path:        path,
		handlers:    handlers,
		middlewares: nil,
	})
	return &Route{
		path:       path,
		method:     method,
		handlers:   handlers,
		middlewares: nil,
		routerImpl: r,
	}
}

func (r *routerImpl) setupToRouter(router Router) {
	for _, rec := range r.routes {
		var allMW []MiddlewareFunc
		allMW = append(allMW, rec.middlewares...)
		wrapped := make([]HandlerFunc, len(allMW)+len(rec.handlers))

		for i, mw := range allMW {
			wrapped[i] = mw
		}
		copy(wrapped[len(allMW):], rec.handlers)

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
			var allMW []MiddlewareFunc
			allMW = append(allMW, g.middlewares...)
			allMW = append(allMW, rec.middlewares...)
			wrapped := make([]HandlerFunc, len(allMW)+len(rec.handlers))

			for i, mw := range allMW {
				wrapped[i] = mw
			}
			copy(wrapped[len(allMW):], rec.handlers)

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

type groupRouter struct {
	impl        *routerImpl
	prefix      string
	middlewares []MiddlewareFunc
	routes      []routeRecord
}

func (g *groupRouter) GET(path string, handlers ...HandlerFunc) *Route {
	g.routes = append(g.routes, routeRecord{
		method:      "GET",
		path:        path,
		handlers:    handlers,
		middlewares: g.middlewares,
	})
	return &Route{
		path:       path,
		method:     "GET",
		handlers:   handlers,
		middlewares: g.middlewares,
		routerImpl: g.impl,
	}
}

func (g *groupRouter) POST(path string, handlers ...HandlerFunc) *Route {
	g.routes = append(g.routes, routeRecord{
		method:      "POST",
		path:        path,
		handlers:    handlers,
		middlewares: g.middlewares,
	})
	return &Route{
		path:       path,
		method:     "POST",
		handlers:   handlers,
		middlewares: g.middlewares,
		routerImpl: g.impl,
	}
}

func (g *groupRouter) PUT(path string, handlers ...HandlerFunc) *Route {
	g.routes = append(g.routes, routeRecord{
		method:      "PUT",
		path:        path,
		handlers:    handlers,
		middlewares: g.middlewares,
	})
	return &Route{
		path:       path,
		method:     "PUT",
		handlers:   handlers,
		middlewares: g.middlewares,
		routerImpl: g.impl,
	}
}

func (g *groupRouter) DELETE(path string, handlers ...HandlerFunc) *Route {
	g.routes = append(g.routes, routeRecord{
		method:      "DELETE",
		path:        path,
		handlers:    handlers,
		middlewares: g.middlewares,
	})
	return &Route{
		path:       path,
		method:     "DELETE",
		handlers:   handlers,
		middlewares: g.middlewares,
		routerImpl: g.impl,
	}
}

func (g *groupRouter) PATCH(path string, handlers ...HandlerFunc) *Route {
	g.routes = append(g.routes, routeRecord{
		method:      "PATCH",
		path:        path,
		handlers:    handlers,
		middlewares: g.middlewares,
	})
	return &Route{
		path:       path,
		method:     "PATCH",
		handlers:   handlers,
		middlewares: g.middlewares,
		routerImpl: g.impl,
	}
}

func (g *groupRouter) WS(path string, handlers ...HandlerFunc) *Route {
	g.routes = append(g.routes, routeRecord{
		method:      "WS",
		path:        path,
		handlers:    handlers,
		middlewares: g.middlewares,
	})
	return &Route{
		path:       path,
		method:     "WS",
		handlers:   handlers,
		middlewares: g.middlewares,
		routerImpl: g.impl,
	}
}

func (g *groupRouter) GROUP(prefix string, mw ...MiddlewareFunc) *RouterGroup {
	return &RouterGroup{
		Prefix:      prefix,
		Router:      g,
		middlewares: mw,
	}
}