package middleware

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/go-sonic/sonic/handler/web"
	"github.com/go-sonic/sonic/handler/web/ginadapter"
)

type RecoveryMiddleware struct {
	logger *zap.Logger
}

func NewRecoveryMiddleware(logger *zap.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		logger: logger,
	}
}

func (r *RecoveryMiddleware) RecoveryWithLogger() gin.HandlerFunc {
	return ginadapter.Wrap(r.Handler())
}

func (r *RecoveryMiddleware) Handler() web.HandlerFunc {
	logger := r.logger.WithOptions(zap.AddCallerSkip(2))

	return func(ctx web.Context) {
		defer func() {
			if panicVal := recover(); panicVal != nil {
				var recoveredErr error
				switch e := panicVal.(type) {
				case error:
					recoveredErr = e
				default:
					recoveredErr = fmt.Errorf("%v", e)
				}
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				//nolint:errorlint
				if ne, ok := recoveredErr.(*net.OpError); ok {
					//nolint:errorlint
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					logger.Error(ctx.Path(),
						zap.Error(recoveredErr),
						zap.String("request_id", GetRequestID(ctx)),
					)
				} else {
					logger.Error("[Recovery] panic recovered",
						zap.Error(recoveredErr),
						zap.String("request_id", GetRequestID(ctx)),
						zap.String("method", ctx.Method()),
						zap.String("path", ctx.Path()),
					)
				}

				if brokenPipe {
					ctx.Abort()
				} else {
					code := http.StatusInternalServerError
					AbortWithErrorJSON(ctx, code, ErrorCodeFromStatus(code), LocalizedHTTPStatusText(ctx, code))
				}
			}
		}()
		ctx.Next()
	}
}
