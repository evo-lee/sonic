package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

type RequestIDMiddleware struct{}

func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{}
}

func (m *RequestIDMiddleware) RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		ctx.Set(RequestIDKey, requestID)
		ctx.Writer.Header().Set(RequestIDHeader, requestID)
		ctx.Next()
	}
}

func GetRequestID(ctx *gin.Context) string {
	requestID, ok := ctx.Get(RequestIDKey)
	if !ok {
		return ""
	}
	if s, ok := requestID.(string); ok {
		return s
	}
	return ""
}
