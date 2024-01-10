package router

import "net/http"

type ErrorHandler func(w http.ResponseWriter)

type RouteHandler func(*http.Request) (ApiResponse, error)

type RouteHandlers struct {
	GET    *RouteHandler
	POST   *RouteHandler
	PUT    *RouteHandler
	PATCH  *RouteHandler
	DELETE *RouteHandler
}

func (m *RouteHandlers) RegisterMethodHandler(method string, handler *RouteHandler) {
	switch method {
	case "GET":
		m.GET = handler
	case "POST":
		m.POST = handler
	case "PUT":
		m.PUT = handler
	case "PATCH":
		m.PATCH = handler
	case "DELETE":
		m.DELETE = handler
	}
}

func (m *RouteHandlers) UnregisterMethodHandler(method string) {
	switch method {
	case "GET":
		m.GET = nil
	case "POST":
		m.POST = nil
	case "PUT":
		m.PUT = nil
	case "PATCH":
		m.PATCH = nil
	case "DELETE":
		m.DELETE = nil
	}
}

func (m *RouteHandlers) GetHandler(method string) *RouteHandler {
	switch method {
	case "GET":
		return m.GET
	case "POST":
		return m.POST
	case "PUT":
		return m.PUT
	case "PATCH":
		return m.PATCH
	case "DELETE":
		return m.DELETE
	}

	return nil
}

func (m *RouteHandlers) IsAllowed(method string) bool {
	switch method {
	case "GET":
		return m.GET != nil
	case "POST":
		return m.POST != nil
	case "PUT":
		return m.PUT != nil
	case "PATCH":
		return m.PATCH != nil
	case "DELETE":
		return m.DELETE != nil
	}

	return false
}
