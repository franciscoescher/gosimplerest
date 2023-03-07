package gosimplerest

import (
	"net/http"
	"strings"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/go-chi/chi"
)

func AddChiHandlers(r *chi.Mux, base AddHandlersBaseParams) *chi.Mux {
	params := AddHandlersParams{
		AddHandlersBaseParams: base,
		AddRouteFunctions: AddRouteFunctions{
			Post:   ChiAddRouteFunc(r.Post),
			Get:    ChiAddRouteFunc(r.Get),
			Put:    ChiAddRouteFunc(r.Put),
			Patch:  ChiAddRouteFunc(r.Patch),
			Delete: ChiAddRouteFunc(r.Delete),
			Head:   ChiAddRouteFunc(r.Head),
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

// ChiAddRouteFunc uses the f function to add a route to the router,
// wrapping the handler to add params to request context.
func ChiAddRouteFunc(f AddRouteFunc) AddRouteFunc {
	return func(name string, h http.HandlerFunc) {
		f(name, ChiHandlerWrapper(h))
	}
}

// ChiHandlerWrapper wraps the handler function.
// It adds params to request context.
func ChiHandlerWrapper(h http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// get params from chi context
		keys := make([]string, 0)
		values := make([]string, 0)
		rctx := chi.RouteContext(r.Context())
		if rctx != nil {
			keys = rctx.URLParams.Keys
			values = rctx.URLParams.Values
		}
		params := make(map[string]string, len(keys))
		for i := range keys {
			params[keys[i]] = values[i]
		}

		handlers.AddParamsToHandlerFunc(h, params)(rw, r)
	}
}
