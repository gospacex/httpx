package httpx

import (
	"context"
	"net/http"
	"os"
	"time"
)

type HandlerFunc func(ctx context.Context, hc HandlerContext)

type HandlerContext interface {
	Request() interface{}
	Response() interface{}
	Param(key string) string
	Query(key string) string
	Bind(into interface{}) error
	AbortWithStatus(code int)
	AbortJSON(code int, body interface{})
}

type MiddlewareFunc func(HandlerFunc) HandlerFunc

type WSConfig struct {
	ReadBufferSize  int
	WriteBufferSize int
	CheckOrigin     func(*http.Request) bool
}

type TLSConfig struct {
	CertFile string
	KeyFile  string
	Addr     string
}

type StartOption struct {
	QuitChan        <-chan os.Signal
	Timeout         time.Duration
	BeforeShutdown  func()
	AfterShutdown   func(err error)
	TLSConfig       *TLSConfig
}

func DefaultStartOption() *StartOption {
	return &StartOption{
		Timeout: 5 * time.Second,
	}
}

func WithTLS(cfg *TLSConfig) Option {
	return func(o *StartOption) {
		o.TLSConfig = cfg
	}
}