package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"encoding/json"
	"time"
)

// CreateHandler returns a handler for the POST method
func CreateHandler(params *GetHandlerFuncParams) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := unmarshalBody(r)
		if err != nil {
			params.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !params.Resource.AutoIncrementalPK {
			pk := params.Resource.GeneratePrimaryKey()
			data[params.Resource.PrimaryKey] = pk
		}
		if params.Resource.CreatedAtField.Valid {
			data[params.Resource.CreatedAtField.String] = time.Now()
		}
		if params.Resource.UpdatedAtField.Valid {
			data[params.Resource.UpdatedAtField.String] = time.Now()
		}
		if params.Resource.SoftDeleteField.Valid {
			data[params.Resource.SoftDeleteField.String] = nil
		}

		// perform data validation
		for key := range data {
			// validates field exists in the model
			if !params.Resource.HasField(key) {
				w.WriteHeader(http.StatusBadRequest)
				err = encodeJsonError(w, r, key+" not in the model")
				if err != nil {
					params.Logger.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
		}
		// validates values
		errs := params.Resource.ValidateAllFields(params.Validate, data)
		if len(errs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			err = encodeJsonError(w, r, fmt.Sprintf("%s", errs))
			if err != nil {
				params.Logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		id, err := params.Repository.Insert(params.Resource, data)
		if err != nil {
			params.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if params.Resource.AutoIncrementalPK {
			data[params.Resource.PrimaryKey] = id
		}

		err = encodeJson(w, r, map[string]any{params.Resource.PrimaryKey: data[params.Resource.PrimaryKey]})
		if err != nil {
			params.Logger.Error(err)
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
