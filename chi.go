package gosimplerest

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
)

func AddChiHandlers(r *chi.Mux, d *sql.DB, l *logrus.Logger, v *validator.Validate, resources []resource.Resource) *chi.Mux {
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

		r.Post(name, ChiHandler(handlers.CreateHandler(base)))
		r.Get(nameID, ChiHandler(handlers.RetrieveHandler(base)))
		r.Put(name, ChiHandler(handlers.UpdateHandler(base)))
		r.Delete(nameID, ChiHandler(handlers.DeleteHandler(base)))
		r.Get(name, ChiHandler(handlers.SearchHandler(base)))

		for _, belongsTo := range resources[i].BelongsToFields {
			nameBelongsTo := fmt.Sprintf("/%s/{id}%s", strcase.KebabCase(belongsTo.Table), name)
			r.Get(nameBelongsTo, ChiHandler(handlers.GetBelongsToHandler(base, belongsTo)))
		}
	}
	return r
}

// ChiHandler wraps the handler function with the given middleware
// It adds params to request context.
// If the middelware function is nil, it returns the handler
func ChiHandler(handlerFunc http.Handler) http.HandlerFunc {
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

		r = handlers.GetRequestWithParams(r, params)
		handlerFunc.ServeHTTP(rw, r)
	}
}