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
				err = encodeJsonError(w, r, key+" not in the model")
				if err != nil {
					base.Logger.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
		}
		// validates values
		errs := base.Resource.ValidateAllFields(base.Validate, data)
		if len(errs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			err = encodeJsonError(w, r, fmt.Sprintf("%s", errs))
			if err != nil {
				base.Logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
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

		err = encodeJson(w, r, map[string]any{base.Resource.PrimaryKey: data[base.Resource.PrimaryKey]})
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// unmarshalBody converts the body of the request to a map where
// the keys are the field names and the values are the field values
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
func encodeJsonError(w http.ResponseWriter, r *http.Request, msg string) error {
	return encodeJson(w, r, map[string]any{
		"error": msg,
	})
}

// encodeJson encodes a json to the response writer.
// if the method is HEAD, it does not write the body, only the headers.
func encodeJson(w http.ResponseWriter, r *http.Request, data interface{}) error {
	jsonResponnse, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(jsonResponnse)))
	if r.Method != http.MethodHead {
		_, err = w.Write(jsonResponnse)
	}

	return err
}
