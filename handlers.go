package gosimplerest

import (
	"bytes"
	"fmt"
	"net/http"

	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

// GetHandler returns a handler for the GET method
func GetHandler(resource Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ReadParams(r, "id")

		// validates id
		err := resource.ValidateField(resource.PrimaryKey, id)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := resource.Find(id)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(result) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(result)
		w.Header().Set("Content-Type", "application/json")
	}
}

// DeleteHandler returns a handler for the DELETE method
func DeleteHandler(resource Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ReadParams(r, "id")

		// validates id
		err := resource.ValidateField(resource.PrimaryKey, id)
		if err != nil {
			logger.Error(err)
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
		if resource.UpdatedAtField.Valid {
			data[resource.UpdatedAtField.String] = time.Now()
		}
		if resource.SoftDeleteField.Valid {
			data[resource.SoftDeleteField.String] = nil
		}

		// perform data validation
		for key := range data {
			// validates field exists in the model
			if !resource.HasField(key) {
				json.NewEncoder(w).Encode(fmt.Sprintf("%s not in the model", key))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// validates value
			err := resource.ValidateField(key, data[key])
			if err != nil {
				logrus.Error(err)
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

		// primary key is required
		_, ok := data[resource.PrimaryKey]
		if !ok {
			json.NewEncoder(w).Encode("primery key is required")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if resource.UpdatedAtField.Valid {
			data[resource.UpdatedAtField.String] = time.Now()
		}

		// perform data validation
		for key := range data {
			// validates field exists in the model
			if !resource.HasField(key) {
				json.NewEncoder(w).Encode(fmt.Sprintf("%s not in the model", key))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// validates value
			err := resource.ValidateField(key, data[key])
			if err != nil {
				logrus.Error(err)
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

// GetBelongsToHandler returns a handler for the GET method of the belongs to relationship
func GetBelongsToHandler(resource Resource, belongsTo BelongsTo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ReadParams(r, "id")

		// validates id
		err := resource.ValidateField(resource.PrimaryKey, id)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := resource.FindFromBelongsTo(id, belongsTo)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(result) == 0 {
			w.WriteHeader(http.StatusNotFound)
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
			// validates fields
			if !resource.IsSearchable(key) {
				json.NewEncoder(w).Encode(fmt.Sprintf("%s is not searchable", key))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// validates values
			for _, v := range query[key] {
				err := resource.ValidateField(key, v)
				if err != nil {
					logrus.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
		}

		result, err := resource.Search(query)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(result) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		json.NewEncoder(w).Encode(result)
		w.Header().Set("Content-Type", "application/json")
	}
}

// unmarschalBody converts the body of the request to a map of strings and interfaces
func unmarschalBody(r *http.Request) (map[string]any, error) {
	b := new(bytes.Buffer)
	b.ReadFrom(r.Body)
	var objmap map[string]any
	err := json.Unmarshal(b.Bytes(), &objmap)
	return objmap, err
}
