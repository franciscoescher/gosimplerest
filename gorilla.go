package gosimplerest

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
)

// AddGorillaMuxHandlers contains an extra parameter for the middleware since it can't be added to the router directly
func AddGorillaMuxHandlers(r *mux.Router, d *sql.DB, l *logrus.Logger, v *validator.Validate, resources []resource.Resource, mid func(h http.Handler) http.HandlerFunc) *mux.Router {
	if v == nil {
		v = validator.New()
	}
	if l == nil {
		l = logrus.New()
		l.Out = io.Discard
	}
	for i := range resources {
		base := &resource.Base{Logger: l, DB: d, Validate: v, Resource: &resources[i]}
		name := fmt.Sprintf("/%s", strcase.KebabCase(resources[i].Table))
		nameID := fmt.Sprintf("%s/{id}", name)

		r.HandleFunc(name, GorillaHandler(mid, handlers.CreateHandler(base))).Methods(http.MethodPost)
		r.HandleFunc(nameID, GorillaHandler(mid, handlers.RetrieveHandler(base))).Methods(http.MethodGet)
		r.HandleFunc(name, GorillaHandler(mid, handlers.UpdateHandler(base))).Methods(http.MethodPut)
		r.HandleFunc(nameID, GorillaHandler(mid, handlers.DeleteHandler(base))).Methods(http.MethodDelete)
		r.HandleFunc(name, GorillaHandler(mid, handlers.SearchHandler(base))).Methods(http.MethodGet)

		for _, belongsTo := range resources[i].BelongsToFields {
			nameBelongsTo := fmt.Sprintf("/%s/{id}%s", strcase.KebabCase(belongsTo.Table), name)
			r.HandleFunc(nameBelongsTo, GorillaHandler(mid, handlers.GetBelongsToHandler(base, belongsTo))).Methods(http.MethodGet)
		}
	}
	return r
}

// GorillaHandler wraps the handler function with the given middleware
// It adds params to request context.
// If the middelware function is nil, it returns the handler
func GorillaHandler(mid func(h http.Handler) http.HandlerFunc, h http.Handler) http.HandlerFunc {
	if mid != nil {
		h = mid(h)
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		r = handlers.GetRequestWithParams(r, params)
		h.ServeHTTP(rw, r)
	}
}
