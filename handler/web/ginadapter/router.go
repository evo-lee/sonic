package ginadapter

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/handler/web"
)

type Router struct {
	router ginRouter
}

type ginRouter interface {
	Use(...gin.HandlerFunc) gin.IRoutes
	Group(string, ...gin.HandlerFunc) *gin.RouterGroup
	GET(string, ...gin.HandlerFunc) gin.IRoutes
	POST(string, ...gin.HandlerFunc) gin.IRoutes
	PUT(string, ...gin.HandlerFunc) gin.IRoutes
	DELETE(string, ...gin.HandlerFunc) gin.IRoutes
	StaticFS(string, http.FileSystem) gin.IRoutes
}

func NewRouter(router ginRouter) web.Router {
	return &Router{router: router}
}

func wrapHandlers(handlers ...web.HandlerFunc) []gin.HandlerFunc {
	if len(handlers) == 0 {
		return nil
	}
	wrapped := make([]gin.HandlerFunc, 0, len(handlers))
	for _, handler := range handlers {
		wrapped = append(wrapped, Wrap(handler))
	}
	return wrapped
}

func (r *Router) Use(handlers ...web.HandlerFunc) {
	r.router.Use(wrapHandlers(handlers...)...)
}

func (r *Router) Group(relativePath string, handlers ...web.HandlerFunc) web.Router {
	return &Router{router: r.router.Group(relativePath, wrapHandlers(handlers...)...)}
}

func (r *Router) GET(relativePath string, handlers ...web.HandlerFunc) {
	r.router.GET(relativePath, wrapHandlers(handlers...)...)
}

func (r *Router) POST(relativePath string, handlers ...web.HandlerFunc) {
	r.router.POST(relativePath, wrapHandlers(handlers...)...)
}

func (r *Router) PUT(relativePath string, handlers ...web.HandlerFunc) {
	r.router.PUT(relativePath, wrapHandlers(handlers...)...)
}

func (r *Router) DELETE(relativePath string, handlers ...web.HandlerFunc) {
	r.router.DELETE(relativePath, wrapHandlers(handlers...)...)
}

func (r *Router) StaticFS(relativePath, root string) {
	r.router.StaticFS(relativePath, gin.Dir(root, false))
}

func (r *Router) Native() any {
	return r.router
}
