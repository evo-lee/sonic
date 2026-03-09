package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/go-sonic/sonic/handler/web"
	"github.com/go-sonic/sonic/handler/web/ginadapter"
)

type LoggerMiddleware struct {
	logger *zap.Logger
}

func NewLoggerMiddleware(logger *zap.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		logger: logger,
	}
}

// LoggerConfig defines the config for Logger middleware
type LoggerConfig struct {
	// SkipPaths is an url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

// LoggerWithConfig instance a Logger middleware with config.
func (g *LoggerMiddleware) LoggerWithConfig(conf LoggerConfig) gin.HandlerFunc {
	return ginadapter.Wrap(g.HandlerWithConfig(conf))
}

func (g *LoggerMiddleware) HandlerWithConfig(conf LoggerConfig) web.HandlerFunc {
	logger := g.logger.WithOptions(zap.WithCaller(false))
	notLogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}

	return func(ctx web.Context) {
		// Start timer
		start := time.Now()
		path := ctx.Path()
		raw := ctx.RawQuery()

		// Process request
		ctx.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			if raw != "" {
				path = path + "?" + raw
			}
			path = strings.ReplaceAll(path, "\n", "")
			path = strings.ReplaceAll(path, "\r", "")
			clientIP := strings.ReplaceAll(ctx.ClientIP(), "\n", "")
			clientIP = strings.ReplaceAll(clientIP, "\r", "")

			logger.Info("[HTTP]",
				zap.Time("beginTime", start),
				zap.Int("status", ctx.StatusCode()),
				zap.Duration("latency", time.Since(start)),
				zap.String("clientIP", clientIP),
				zap.String("method", ctx.Method()),
				zap.String("path", path),
				zap.String("request_id", GetRequestID(ctx)))
		}
	}
}
