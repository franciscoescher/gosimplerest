package gosimplerest

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

// GetHandler returns a handler for the GET method
func GetHandler(resource Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		// validates id
		err := validateField(resource, resource.PrimaryKey, id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := resource.Find(id)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
		w.Header().Set("Content-Type", "application/json")
	}
}

// DeleteHandler returns a handler for the DELETE method
func DeleteHandler(resource Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		// validates id
		err := validateField(resource, resource.PrimaryKey, id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = resource.Delete(id)
		if err != nil {
			if err.Error() == "no rows affected" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// CreateHandler returns a handler for the POST method
func CreateHandler(resource Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := unmarschalBody(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		pk := resource.GeneratePrimaryKey()
		data[resource.PrimaryKey] = pk
		if resource.CreatedAtField.Valid {
			data[resource.CreatedAtField.String] = time.Now()
		}
		if resource.SoftDeleteField.Valid {
			data[resource.SoftDeleteField.String] = nil
		}
		if resource.UpdatedAtField.Valid {
			data[resource.UpdatedAtField.String] = time.Now()
		}

		// validates that all fields in data are in the model
		for key := range data {
			if !resource.HasField(key) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// validates fields
			err := validateField(resource, key, data[key])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		err = resource.Insert(data)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(pk)
	}
}

// UpdateHandler returns a handler for the PUT method
func UpdateHandler(resource Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := unmarschalBody(r)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, ok := data[resource.PrimaryKey]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if resource.UpdatedAtField.Valid {
			data[resource.UpdatedAtField.String] = time.Now()
		}

		// validates that all fields in data are in the model
		for key := range data {
			if !resource.HasField(key) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		affected, err := resource.Update(data)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if affected == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

// unmarschalBody converts the body of the request to a map of strings and interfaces
func unmarschalBody(r *http.Request) (map[string]interface{}, error) {
	b := new(bytes.Buffer)
	b.ReadFrom(r.Body)
	var objmap map[string]interface{}
	err := json.Unmarshal(b.Bytes(), &objmap)
	return objmap, err
}

func validateField(resource Resource, field string, value interface{}) error {
	vf := resource.Fields[field].Validator
	if vf != nil {
		err := vf(field, value)
		return err
	}
	return nil
}

// GetBelongsToHandler returns a handler for the GET method of the belongs to relationship
func GetBelongsToHandler(resource Resource, belongsTo BelongsTo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		// validates id
		err := validateField(resource, resource.PrimaryKey, id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := resource.FindFromBelongsTo(id, belongsTo)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
		w.Header().Set("Content-Type", "application/json")
	}
}

// SearchHandler returns a handler for the GET method with query params
func SearchHandler(resource Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// validates that all fields in data are in the model
		for key := range query {
			if !resource.HasField(key) {
				logrus.Error(fmt.Errorf("field %s unkown", key))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			for _, v := range query[key] {
				// validates fields
				err := validateField(resource, key, v)
				if err != nil {
					logrus.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
		}

		result, err := resource.Search(query)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
		w.Header().Set("Content-Type", "application/json")
	}
}
