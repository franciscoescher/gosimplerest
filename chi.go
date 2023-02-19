package gosimplerest

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func AddChiHandlers(r *chi.Mux, d *sql.DB, l *logrus.Logger, v *validator.Validate, resources []resource.Resource) *chi.Mux {
	h := AddRouteFunctions{
		Post:   ChiAddRouteFunc(r.Post),
		Get:    ChiAddRouteFunc(r.Get),
		Put:    ChiAddRouteFunc(r.Put),
		Patch:  ChiAddRouteFunc(r.Patch),
		Delete: ChiAddRouteFunc(r.Delete),
		Head:   ChiAddRouteFunc(r.Head),
	}
	apf := func(name string, param string) string {
		var sb strings.Builder
		sb.WriteString(name)
		sb.WriteString("/{")
		sb.WriteString(param)
		sb.WriteString("}")
		return sb.String()
	}
	AddHandlers(d, l, v, h, apf, resources)
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
