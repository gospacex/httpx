package middleware

type MiddlewareConfig struct {
	RateLimitQPS int

	RecoveryMode string

	JWTSecret string
	JWTIssuer string

	CORSAllowedOrigins []string
	CORSAllowedMethods []string

	RequestIDHeader string

	LogFormat string
}

type MiddlewareOption func(*MiddlewareConfig)

func WithRateLimit(qps int) MiddlewareOption {
	return func(o *MiddlewareConfig) {
		o.RateLimitQPS = qps
	}
}

func WithRecovery(mode string) MiddlewareOption {
	return func(o *MiddlewareConfig) {
		o.RecoveryMode = mode
	}
}

func WithJWT(secret, issuer string) MiddlewareOption {
	return func(o *MiddlewareConfig) {
		o.JWTSecret = secret
		o.JWTIssuer = issuer
	}
}

func WithCORS(origins []string, methods []string) MiddlewareOption {
	return func(o *MiddlewareConfig) {
		o.CORSAllowedOrigins = origins
		o.CORSAllowedMethods = methods
	}
}

func WithRequestID(header string) MiddlewareOption {
	return func(o *MiddlewareConfig) {
		o.RequestIDHeader = header
	}
}

func WithLogFormat(format string) MiddlewareOption {
	return func(o *MiddlewareConfig) {
		o.LogFormat = format
	}
}