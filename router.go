package gosimplerest

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
)

var logger *logrus.Logger
var db *sql.DB

/*
HandleResources creates a REST API for the given models.
It creates the following routes (models are table names of the resources, converted to kebab case):

GET /model/{id}
POST /model
PUT /model
DELETE /model/{id}
GET /model

Also, for each belongs to relation, it creates the following routes:

GET /belongs-to/{id}/model

The handlers parameter is a function to wrap the handlers with, for example, authentication and logging
*/
func AddHandlers(d *sql.DB, l *logrus.Logger, r *mux.Router, resources []Resource, mid func(h http.Handler) http.HandlerFunc) *mux.Router {
	db = d
	logger = l

	for _, resource := range resources {
		name := fmt.Sprintf("/%s", strcase.KebabCase(resource.Table))
		nameID := fmt.Sprintf("%s/{id}", name)

		r.HandleFunc(nameID, Middelwares(mid, GetHandler(resource))).Methods(http.MethodGet)
		r.HandleFunc(nameID, Middelwares(mid, DeleteHandler(resource))).Methods(http.MethodDelete)
		r.HandleFunc(name, Middelwares(mid, CreateHandler(resource))).Methods(http.MethodPost)
		r.HandleFunc(name, Middelwares(mid, UpdateHandler(resource))).Methods(http.MethodPut)
		r.HandleFunc(name, Middelwares(mid, SearchHandler(resource))).Methods(http.MethodGet)

		for _, belongsTo := range resource.BelongsToFields {
			nameBelongsTo := fmt.Sprintf("/%s/{id}%s", strcase.KebabCase(belongsTo.Table), name)
			r.HandleFunc(nameBelongsTo, Middelwares(mid, GetBelongsToHandler(resource, belongsTo))).Methods(http.MethodGet)
		}
	}
	return r
}

// Middelwares wraps the handler function with the given middleware
// If the middelware function is nil, it returns the handler
func Middelwares(mid func(h http.Handler) http.HandlerFunc, handlerFunc http.Handler) http.HandlerFunc {
	if mid == nil {
		return func(rw http.ResponseWriter, req *http.Request) {
			handlerFunc.ServeHTTP(rw, req)
		}
	}
	return mid(handlerFunc)
}
