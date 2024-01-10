package router

import "net/http"

type Middleware func(r *http.Request) error
