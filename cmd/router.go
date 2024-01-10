package router

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Router struct {
	server            *http.Server
	globalMiddlewares []Middleware
	routes            map[string]*Route
	errorHandlers     map[string]*ErrorHandler
}

func New(s *http.Server) Router {
	return Router{
		s,
		[]Middleware{},
		make(map[string]*Route),
		make(map[string]*ErrorHandler),
	}
}

func (router *Router) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		return
	}

	if errHandler, found := router.errorHandlers[err.Error()]; found {
		(*errHandler)(w)
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (router *Router) handleMiddlewares(pattern string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Global middlewares
		for m := range router.globalMiddlewares {
			err := router.globalMiddlewares[m](r)

			if err != nil {
				router.handleError(w, err)
				return
			}
		}

		// Individual middlewares
		route, found := router.routes[pattern]
		if found {
			for m := range route.Middlewares {
				err := route.Middlewares[m](r)

				if err != nil {
					router.handleError(w, err)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (router *Router) makeRoute(pattern string) http.HandlerFunc {
	return router.handleMiddlewares(pattern, func(w http.ResponseWriter, r *http.Request) {
		route := router.routes[pattern]

		if !route.RouteHandlers.IsAllowed(r.Method) {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler := route.RouteHandlers.GetHandler(r.Method)

		apiResponse, err := (*handler)(r)

		if err != nil {
			router.handleError(w, err)
			return
		}

		response, err := json.Marshal(apiResponse)

		if err != nil {
			router.handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})
}

func (router *Router) register(method string, pattern string, handler *RouteHandler) {
	_, found := router.routes[pattern]

	// Bind route pattern if not already binded
	if !found {
		router.routes[pattern] = &Route{RouteHandlers{}, []Middleware{}}
		http.HandleFunc(pattern, router.makeRoute(pattern))
	}

	route := router.routes[pattern]

	if route.RouteHandlers.IsAllowed(method) {
		panic(fmt.Errorf("Route %s: %s already assigned", pattern, method))
	}

	route.RouteHandlers.RegisterMethodHandler(method, handler)
}

func (router *Router) Get(pattern string, handler RouteHandler) {
	router.register("GET", pattern, &handler)
}

func (router *Router) Post(pattern string, handler RouteHandler) {
	router.register("POST", pattern, &handler)
}

func (router *Router) Put(pattern string, handler RouteHandler) {
	router.register("PUT", pattern, &handler)
}

func (router *Router) Patch(pattern string, handler RouteHandler) {
	router.register("PATCH", pattern, &handler)
}

func (router *Router) Delete(pattern string, handler RouteHandler) {
	router.register("DELETE", pattern, &handler)
}

func (router *Router) UseGlobalMiddleware(middleware Middleware) {
	router.globalMiddlewares = append(router.globalMiddlewares, middleware)
}

func (router *Router) UseMiddleware(pattern string, middleware Middleware) {
	route, found := router.routes[pattern]

	if !found {
		panic(fmt.Sprintf("Route %s not registered", pattern))
	}

	route.Middlewares = append(route.Middlewares, middleware)
}

func (router *Router) UseErrorHandler(err string, errorHandler ErrorHandler) {
	router.errorHandlers[err] = &errorHandler
}

func (router *Router) Start() error {
	fmt.Println("Server running on port " + router.server.Addr)

	return router.server.ListenAndServe()
}
