package router

type Route struct {
	RouteHandlers RouteHandlers
	Middlewares   []Middleware
}
