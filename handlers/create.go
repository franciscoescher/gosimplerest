package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"encoding/json"
	"time"

	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/sirupsen/logrus"
)

// CreateHandler returns a handler for the POST method
func CreateHandler(base *resource.Base) http.HandlerFunc {
	if base.Resource.OmitCreateRoute {
		return NotFoundHandler
	}
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
				encodeJsonError(w, fmt.Sprintf("%s not in the model", key))
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

		result := map[string]any{base.Resource.PrimaryKey: pk}
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

// encodeJsonError encodes an error message in json to the response writer
func encodeJsonError(w http.ResponseWriter, msg string) {
	json.NewEncoder(w).Encode(map[string]any{
		"error": msg,
	})
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
