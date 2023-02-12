package gosimplerest

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
)

func AddGorillaMuxHandlers(r *mux.Router, d *sql.DB, l *logrus.Logger, resources []resource.Resource, mid func(h http.Handler) http.HandlerFunc) *mux.Router {
	for i := range resources {
		base := &resource.Base{Logger: l, DB: d, Resource: &resources[i]}
		name := fmt.Sprintf("/%s", strcase.KebabCase(resources[i].Table))
		nameID := fmt.Sprintf("%s/{id}", name)

		r.HandleFunc(name, GorillaMiddelwares(mid, handlers.CreateHandler(base))).Methods(http.MethodPost)
		r.HandleFunc(nameID, GorillaMiddelwares(mid, handlers.RetrieveHandler(base))).Methods(http.MethodGet)
		r.HandleFunc(name, GorillaMiddelwares(mid, handlers.UpdateHandler(base))).Methods(http.MethodPut)
		r.HandleFunc(nameID, GorillaMiddelwares(mid, handlers.DeleteHandler(base))).Methods(http.MethodDelete)
		r.HandleFunc(name, GorillaMiddelwares(mid, handlers.SearchHandler(base))).Methods(http.MethodGet)

		for _, belongsTo := range resources[i].BelongsToFields {
			nameBelongsTo := fmt.Sprintf("/%s/{id}%s", strcase.KebabCase(belongsTo.Table), name)
			r.HandleFunc(nameBelongsTo, GorillaMiddelwares(mid, handlers.GetBelongsToHandler(base, belongsTo))).Methods(http.MethodGet)
		}
	}
	return r
}

// GorrileGorillaMiddelwares wraps the handler function with the given middleware
// It adds params to request context.
// If the middelware function is nil, it returns the handler
func GorillaMiddelwares(mid func(h http.Handler) http.HandlerFunc, handlerFunc http.Handler) http.HandlerFunc {
	if mid != nil {
		handlerFunc = mid(handlerFunc)
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		r = handlers.GetRequestWithParams(r, params)
		handlerFunc.ServeHTTP(rw, r)
	}
}
