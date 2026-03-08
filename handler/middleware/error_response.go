package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/util/xerr"
)

func ErrorCodeFromError(err error) string {
	switch xerr.GetType(err) {
	case xerr.BadParam:
		return "bad_request"
	case xerr.NoRecord:
		return "not_found"
	case xerr.Forbidden:
		return "forbidden"
	case xerr.DB:
		return "db_error"
	case xerr.Email:
		return "email_error"
	default:
		status := xerr.GetHTTPStatus(err)
		return ErrorCodeFromStatus(status)
	}
}

func ErrorCodeFromStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not_found"
	default:
		return "internal_error"
	}
}

func BuildErrorDTO(ctx *gin.Context, status int, code, message string) *dto.BaseDTO {
	if code == "" {
		code = ErrorCodeFromStatus(status)
	}
	return &dto.BaseDTO{
		Status:    status,
		Message:   message,
		Code:      code,
		RequestID: GetRequestID(ctx),
	}
}

func AbortWithErrorJSON(ctx *gin.Context, status int, code, message string) {
	ctx.AbortWithStatusJSON(status, BuildErrorDTO(ctx, status, code, message))
}

func abortWithStatusJSON(ctx *gin.Context, status int, message string) {
	AbortWithErrorJSON(ctx, status, ErrorCodeFromStatus(status), message)
}
