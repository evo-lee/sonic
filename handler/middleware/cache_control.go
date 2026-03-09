package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/handler/web"
	"github.com/go-sonic/sonic/handler/web/ginadapter"
)

type CacheControlMiddleware struct {
	MaxAge time.Duration
	Public bool
}

type CacheControlOption func(*CacheControlMiddleware)

func NewCacheControlMiddleware(opts ...CacheControlOption) *CacheControlMiddleware {
	c := &CacheControlMiddleware{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *CacheControlMiddleware) CacheControl() gin.HandlerFunc {
	return ginadapter.Wrap(c.Handler())
}

func (c *CacheControlMiddleware) Handler() web.HandlerFunc {
	value := ""
	if c.Public {
		value = "public,"
	}
	if c.MaxAge > 0 {
		value = "max-age=" + strconv.FormatInt(int64(c.MaxAge.Seconds()), 10)
	}
	return func(ctx web.Context) {
		ctx.SetHeader("Cache-Control", value)
	}
}

func WithMaxAge(maxAge time.Duration) CacheControlOption {
	return func(c *CacheControlMiddleware) {
		c.MaxAge = maxAge
	}
}

func WithPublic(public bool) CacheControlOption {
	return func(c *CacheControlMiddleware) {
		c.Public = public
	}
}
