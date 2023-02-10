package gosimplerest

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
)

func AddGorillaMuxHandlers(r *mux.Router, d *sql.DB, l *logrus.Logger, resources []Resource, mid func(h http.Handler) http.HandlerFunc) *mux.Router {
	db = d
	logger = l

	for _, resource := range resources {
		name := fmt.Sprintf("/%s", strcase.KebabCase(resource.Table))
		nameID := fmt.Sprintf("%s/{id}", name)

		r.HandleFunc(nameID, GorillaMiddelwares(mid, GetHandler(resource))).Methods(http.MethodGet)
		r.HandleFunc(nameID, GorillaMiddelwares(mid, DeleteHandler(resource))).Methods(http.MethodDelete)
		r.HandleFunc(name, GorillaMiddelwares(mid, CreateHandler(resource))).Methods(http.MethodPost)
		r.HandleFunc(name, GorillaMiddelwares(mid, UpdateHandler(resource))).Methods(http.MethodPut)
		r.HandleFunc(name, GorillaMiddelwares(mid, SearchHandler(resource))).Methods(http.MethodGet)

		for _, belongsTo := range resource.BelongsToFields {
			nameBelongsTo := fmt.Sprintf("/%s/{id}%s", strcase.KebabCase(belongsTo.Table), name)
			r.HandleFunc(nameBelongsTo, GorillaMiddelwares(mid, GetBelongsToHandler(resource, belongsTo))).Methods(http.MethodGet)
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
		r = GetRequestWithParams(r, params)
		mid(handlerFunc).ServeHTTP(rw, r)
	}
}
