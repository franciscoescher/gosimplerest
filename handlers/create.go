package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"encoding/json"
	"time"

	"github.com/franciscoescher/gosimplerest/resource"
)

// CreateHandler returns a handler for the POST method
func CreateHandler(base *resource.Base) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := unmarshalBody(r)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !base.Resource.AutoIncrementalPK {
			pk := base.Resource.GeneratePrimaryKey()
			data[base.Resource.PrimaryKey] = pk
		}
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
		}
		// validates values
		errs := base.Resource.ValidateAllFields(base.Validate, data)
		if len(errs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			encodeJsonError(w, fmt.Sprintf("%s", errs))
			return
		}

		id, err := base.Resource.Insert(base, data)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if base.Resource.AutoIncrementalPK {
			data[base.Resource.PrimaryKey] = id
		}

		result := map[string]any{base.Resource.PrimaryKey: data[base.Resource.PrimaryKey]}
		json.NewEncoder(w).Encode(result)
	}
}

// unmarshalBody converts the body of the request to a map of strings and interfaces
func unmarshalBody(r *http.Request) (map[string]any, error) {
	b := new(bytes.Buffer)
	_, err := b.ReadFrom(r.Body)
	if err != nil {
		return nil, err
	}
	var objmap map[string]any
	err = json.Unmarshal(b.Bytes(), &objmap)
	return objmap, err
}

// encodeJsonError encodes an error message in json to the response writer
func encodeJsonError(w http.ResponseWriter, msg string) {
	json.NewEncoder(w).Encode(map[string]any{
		"error": msg,
	})
}
