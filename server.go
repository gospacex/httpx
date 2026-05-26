package httpx

import (
	"context"
	"os"
	"os/signal"
	"time"
)

type Server interface {
	Start(addr string) error
	Stop(ctx context.Context) error
	Router() Router
	Use(middleware ...MiddlewareFunc) Server
	EnableWS(wsConfig *WSConfig) Server
}

type TestableServer interface {
	Server
	// Engine returns the underlying HTTP server engine for testing purposes
	Engine() interface{}
}

type GracefulServer interface {
	Server
	StartWithGraceful(opts ...*StartOption) error
	GracefulShutdown(ctx context.Context) error
	IsRunning() bool
}

type Option func(*StartOption)

func WithQuitSignal(signals ...os.Signal) Option {
	return func(o *StartOption) {
		if len(signals) > 0 {
			sigChan := make(chan os.Signal, 1)
			for _, s := range signals {
				signal.Notify(sigChan, s)
			}
			o.QuitChan = sigChan
		}
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *StartOption) {
		o.Timeout = timeout
	}
}