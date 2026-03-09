package web

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
)

type RequestBinder interface {
	Bind(*http.Request, any) error
}

// Context is the minimal request/response surface shared by framework adapters.
type Context interface {
	context.Context
	RequestContext() context.Context
	Request() *http.Request
	Writer() io.Writer
	Method() string
	Path() string
	RawQuery() string
	ClientIP() string
	Header(string) string
	ResponseHeader(string) string
	SetHeader(string, string)
	StatusCode() int
	Query(string) (string, bool)
	Param(string) string
	Cookie(string) (string, error)
	SetCookie(string, string, int, string, string, bool, bool)
	MultipartForm() (*multipart.Form, error)
	Set(string, any)
	Get(any) (any, bool)
	Bind(any) error
	BindJSON(any) error
	BindQuery(any) error
	BindWith(any, any) error
	JSON(int, any)
	AbortWithStatusJSON(int, any)
	String(int, string)
	Status(int)
	Redirect(int, string)
	File(string)
	FormFile(string) (*multipart.FileHeader, error)
	Abort()
	Next()
	Native() any
}
