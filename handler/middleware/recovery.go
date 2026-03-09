package middleware

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

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
	logger := r.logger.WithOptions(zap.AddCallerSkip(2))

	return func(ctx *gin.Context) {
		webCtx := ginadapter.NewContext(ctx)
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
					logger.Error(ctx.Request.URL.Path,
						zap.Error(recoveredErr),
						zap.String("request_id", GetRequestID(webCtx)),
					)
				} else {
					logger.Error("[Recovery] panic recovered",
						zap.Error(recoveredErr),
						zap.String("request_id", GetRequestID(webCtx)),
						zap.String("method", ctx.Request.Method),
						zap.String("path", ctx.Request.URL.Path),
					)
				}

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					ctx.Error(recoveredErr) // nolint: errcheck
					ctx.Abort()
				} else {
					code := http.StatusInternalServerError
					AbortWithErrorJSON(webCtx, code, ErrorCodeFromStatus(code), LocalizedHTTPStatusText(webCtx, code))
				}
			}
		}()
		ctx.Next()
	}
}
