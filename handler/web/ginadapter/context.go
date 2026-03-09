package ginadapter

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/go-sonic/sonic/handler/web"
)

type Context struct {
	ctx *gin.Context
}

func NewContext(ctx *gin.Context) web.Context {
	return &Context{ctx: ctx}
}

func Unwrap(ctx web.Context) *gin.Context {
	if ctx == nil {
		return nil
	}
	if ginCtx, ok := ctx.Native().(*gin.Context); ok {
		return ginCtx
	}
	return nil
}

func Wrap(handler func(web.Context)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler(NewContext(ctx))
	}
}

func (c *Context) RequestContext() context.Context {
	return c.ctx
}

func (c *Context) Request() *http.Request {
	return c.ctx.Request
}

func (c *Context) Writer() io.Writer {
	return c.ctx.Writer
}

func (c *Context) Method() string {
	return c.ctx.Request.Method
}

func (c *Context) Path() string {
	return c.ctx.Request.URL.Path
}

func (c *Context) RawQuery() string {
	return c.ctx.Request.URL.RawQuery
}

func (c *Context) ClientIP() string {
	return c.ctx.ClientIP()
}

func (c *Context) Header(key string) string {
	return c.ctx.GetHeader(key)
}

func (c *Context) SetHeader(key, value string) {
	c.ctx.Writer.Header().Set(key, value)
}

func (c *Context) Query(key string) (string, bool) {
	return c.ctx.GetQuery(key)
}

func (c *Context) Param(key string) string {
	return c.ctx.Param(key)
}

func (c *Context) Cookie(name string) (string, error) {
	return c.ctx.Cookie(name)
}

func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	c.ctx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	return c.ctx.MultipartForm()
}

func (c *Context) Set(key string, value any) {
	c.ctx.Set(key, value)
}

func (c *Context) Get(key any) (any, bool) {
	return c.ctx.Get(key)
}

func (c *Context) Bind(value any) error {
	return c.ctx.ShouldBind(value)
}

func (c *Context) BindJSON(value any) error {
	return c.ctx.ShouldBindJSON(value)
}

func (c *Context) BindQuery(value any) error {
	return c.ctx.ShouldBindQuery(value)
}

func (c *Context) BindWith(value any, binder any) error {
	if requestBinder, ok := binder.(web.RequestBinder); ok {
		return requestBinder.Bind(c.ctx.Request, value)
	}
	ginBinder, ok := binder.(binding.Binding)
	if !ok {
		return fmt.Errorf("unsupported binder type %T", binder)
	}
	return c.ctx.ShouldBindWith(value, ginBinder)
}

func (c *Context) JSON(status int, value any) {
	c.ctx.JSON(status, value)
}

func (c *Context) AbortWithStatusJSON(status int, value any) {
	c.ctx.AbortWithStatusJSON(status, value)
}

func (c *Context) String(status int, value string) {
	c.ctx.String(status, value)
}

func (c *Context) Status(status int) {
	c.ctx.Status(status)
}

func (c *Context) Redirect(status int, location string) {
	c.ctx.Redirect(status, location)
}

func (c *Context) File(path string) {
	c.ctx.File(path)
}

func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	return c.ctx.FormFile(name)
}

func (c *Context) Abort() {
	c.ctx.Abort()
}

func (c *Context) Next() {
	c.ctx.Next()
}

func (c *Context) Native() any {
	return c.ctx
}
