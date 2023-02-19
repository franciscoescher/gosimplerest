package handlers

import (
	"fmt"
	"net/http"

	"time"

	"github.com/franciscoescher/gosimplerest/resource"
)

// UpdateHandler returns a handler for the PATCH method
func UpdateHandler(base *resource.Base) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := unmarshalBody(r)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// primary key is required
		_, ok := data[base.Resource.PrimaryKey]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			err = encodeJsonError(w, r, "primery key is required")
			if err != nil {
				base.Logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		// Checks for immutable fields being updated and adds missing fields if method is PUT
		for key, field := range base.Resource.Fields {
			_, ok := data[key]
			// checks for tentative of updating immutable fields
			if ok && field.Immutable {
				w.WriteHeader(http.StatusBadRequest)
				err = encodeJsonError(w, r, key+" is immutable")
				if err != nil {
					base.Logger.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			} else if r.Method == http.MethodPut && !ok && !field.Immutable {
				// if method is PUT and field is not immutable and not present in the request,
				// adds it to the data for update
				data[key] = nil
			}
		}

		if base.Resource.UpdatedAtField.Valid {
			data[base.Resource.UpdatedAtField.String] = time.Now()
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
		errs := base.Resource.ValidateInputFields(base.Validate, data)
		if len(errs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			err = encodeJsonError(w, r, fmt.Sprintf("%s", errs))
			if err != nil {
				base.Logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
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
