package gosimplerest

import (
	"net/http"
	"strings"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/gorilla/mux"
)

func AddGorillaMuxHandlers(r *mux.Router, base AddHandlersBaseParams, mid func(h http.Handler) http.HandlerFunc) *mux.Router {
	params := AddHandlersParams{
		AddHandlersBaseParams: base,
		AddRouteFunctions: AddRouteFunctions{
			Post:   GorillaAddRouteFunc(r, mid, http.MethodPost),
			Get:    GorillaAddRouteFunc(r, mid, http.MethodGet),
			Put:    GorillaAddRouteFunc(r, mid, http.MethodPut),
			Patch:  GorillaAddRouteFunc(r, mid, http.MethodPatch),
			Delete: GorillaAddRouteFunc(r, mid, http.MethodDelete),
			Head:   GorillaAddRouteFunc(r, mid, http.MethodHead),
		},
		AddParamFunc: func(name string, param string) string {
			var sb strings.Builder
			sb.WriteString(name)
			sb.WriteString("/{")
			sb.WriteString(param)
			sb.WriteString("}")
			return sb.String()
		},
	}
	AddHandlers(params)
	return r
}

// GorillaAddRouteFunc is used to add a route to the router, using the given method.
// It adds the middleware function to the handler.
// It adds params to request context.
func GorillaAddRouteFunc(r *mux.Router, mid func(h http.Handler) http.HandlerFunc, method string) AddRouteFunc {
	return func(name string, h http.HandlerFunc) {
		r.HandleFunc(name, GorillaHandlerWrapper(mid, h)).Methods(method)
	}
}

// GorillaHandlerWrapper wraps the handler function.
// It adds params to request context.
// If the middelware function is nil, it returns the handler
func GorillaHandlerWrapper(mid func(h http.Handler) http.HandlerFunc, h http.HandlerFunc) http.HandlerFunc {
	if mid != nil {
		h = mid(h)
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		r = handlers.GetRequestWithParams(r, params)
		h.ServeHTTP(rw, r)
	}
}
