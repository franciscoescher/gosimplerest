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

		// adds missing fields if method is PUT
		if r.Method == http.MethodPut {
			for field := range base.Resource.Fields {
				if _, ok := data[field]; !ok {
					data[field] = nil
				}
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
				err = encodeJsonError(w, r, fmt.Sprintf("%s not in the model", key))
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
