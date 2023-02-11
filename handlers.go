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
func GetHandler(base Base) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ReadParams(r, "id")

		// validates id
		err := base.Resource.ValidateField(base.Resource.PrimaryKey, id)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := base.Resource.Find(base, id)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(result) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("not found")
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

// DeleteHandler returns a handler for the DELETE method
func DeleteHandler(base Base) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ReadParams(r, "id")

		// validates id
		err := base.Resource.ValidateField(base.Resource.PrimaryKey, id)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = base.Resource.Delete(base, id)
		if err != nil {
			if err.Error() == "no rows affected" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// CreateHandler returns a handler for the POST method
func CreateHandler(base Base) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := unmarschalBody(r)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		pk := base.Resource.GeneratePrimaryKey()
		data[base.Resource.PrimaryKey] = pk
		if base.Resource.CreatedAtField.Valid {
			data[base.Resource.CreatedAtField.String] = time.Now()
		}
		if base.Resource.UpdatedAtField.Valid {
			data[base.Resource.UpdatedAtField.String] = time.Now()
		}
		if base.Resource.SoftDeleteField.Valid {
			data[base.Resource.SoftDeleteField.String] = nil
		}

		// perform data validation
		for key := range data {
			// validates field exists in the model
			if !base.Resource.HasField(key) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(fmt.Sprintf("%s not in the model", key))
				return
			}
			// validates value
			err := base.Resource.ValidateField(key, data[key])
			if err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		err = base.Resource.Insert(base, data)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(pk)
	}
}

// UpdateHandler returns a handler for the PUT method
func UpdateHandler(base Base) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := unmarschalBody(r)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// primary key is required
		_, ok := data[base.Resource.PrimaryKey]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("primery key is required")
			return
		}

		if base.Resource.UpdatedAtField.Valid {
			data[base.Resource.UpdatedAtField.String] = time.Now()
		}

		// perform data validation
		for key := range data {
			// validates field exists in the model
			if !base.Resource.HasField(key) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(fmt.Sprintf("%s not in the model", key))
				return
			}
			// validates value
			err := base.Resource.ValidateField(key, data[key])
			if err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		affected, err := base.Resource.Update(base, data)
		if err != nil {
			base.Logger.Error(err)
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
func GetBelongsToHandler(base Base, belongsTo BelongsTo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ReadParams(r, "id")

		// validates id
		err := base.Resource.ValidateField(base.Resource.PrimaryKey, id)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := base.Resource.FindFromBelongsTo(base, id, belongsTo)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(result) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

// SearchHandler returns a handler for the GET method with query params
func SearchHandler(base Base) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// validates that all fields in data are in the model
		for key := range query {
			// validates fields
			if !base.Resource.IsSearchable(key) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(fmt.Sprintf("%s is not searchable", key))
				return
			}
			// validates values
			for _, v := range query[key] {
				err := base.Resource.ValidateField(key, v)
				if err != nil {
					logrus.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
		}

		result, err := base.Resource.Search(base, query)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(result) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		json.NewEncoder(w).Encode(result)
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
