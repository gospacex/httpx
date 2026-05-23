package httpx

type Route struct {
	path        string
	method      string
	handlers    []HandlerFunc
	middlewares []MiddlewareFunc
	routerImpl  *routerImpl
}

func (r *Route) Use(mw ...MiddlewareFunc) *Route {
	r.middlewares = append(r.middlewares, mw...)
	return r
}