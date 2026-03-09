package web

type HandlerFunc func(Context)

type Router interface {
	Use(...HandlerFunc)
	Group(string, ...HandlerFunc) Router
	GET(string, ...HandlerFunc)
	POST(string, ...HandlerFunc)
	PUT(string, ...HandlerFunc)
	DELETE(string, ...HandlerFunc)
	StaticFS(string, string)
	Native() any
}
