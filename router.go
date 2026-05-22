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