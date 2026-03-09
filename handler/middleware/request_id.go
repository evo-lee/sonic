package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/go-sonic/sonic/handler/web"
	"github.com/go-sonic/sonic/handler/web/ginadapter"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

type RequestIDMiddleware struct{}

func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{}
}

func (m *RequestIDMiddleware) apply(ctx web.Context) {
	requestID := ctx.Header(RequestIDHeader)
	if requestID == "" {
		requestID = uuid.NewString()
	}
	ctx.Set(RequestIDKey, requestID)
	ctx.SetHeader(RequestIDHeader, requestID)
	ctx.Next()
}

func (m *RequestIDMiddleware) Handler() web.HandlerFunc {
	return m.apply
}

func (m *RequestIDMiddleware) RequestID() gin.HandlerFunc {
	return ginadapter.Wrap(m.Handler())
}

func GetRequestID(ctx web.Context) string {
	requestID, ok := ctx.Get(RequestIDKey)
	if !ok {
		return ""
	}
	if s, ok := requestID.(string); ok {
		return s
	}
	return ""
}
