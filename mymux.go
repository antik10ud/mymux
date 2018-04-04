package mymux

import (
	"net/http"
)

type ConsumeResult int

const (
	Unconsumed     ConsumeResult = iota
	MethodMismatch
	Consumed
)

type URLVars map[string]string

func (vars *URLVars) And(vars2 URLVars) URLVars {
	vr := URLVars{}
	for k, v := range *vars {
		vr[k] = v
	}
	for k, v := range vars2 {
		vr[k] = v
	}
	return vr
}

type Route interface {
	consume(http.ResponseWriter, *http.Request) ConsumeResult
	URL(URLVars) string
}

type ErrorFunc func(w http.ResponseWriter, status int, detail string)

type RouteHandler struct {
	routes []Route
	error  ErrorFunc
}

func NewRouteHandler() *RouteHandler {
	rh := RouteHandler{}
	return &rh

}

func (handler *RouterTemplateHandler) ErrorHandler(errorFunc ErrorFunc) {
	handler.error = errorFunc
}

func (h *RouteHandler) HandleFunc(route Route) Route {
	h.routes = append(h.routes, route)
	return route
}

func (h *RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methodMismatch := false
	for _, route := range h.routes {
		switch route.consume(w, r) {
		case MethodMismatch:
			methodMismatch = true
			break
		case Consumed:
			return
		}
	}
	status := http.StatusNotFound
	if methodMismatch {
		status = http.StatusMethodNotAllowed
	}
	if h.error != nil {
		h.error(w, status, http.StatusText(status))
	} else {
		http.Error(w, http.StatusText(status), status)
	}
	return

}
