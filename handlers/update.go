package handlers

import (
	"fmt"
	"net/http"

	"time"

	"github.com/franciscoescher/gosimplerest/resource"
)

// UpdateHandler returns a handler for the PUT method
func UpdateHandler(base *resource.Base) http.HandlerFunc {
	if base.Resource.OmitUpdateRoute {
		return NotFoundHandler
	}
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := unmarschalBody(r)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// primary key is required
		_, ok := data[base.Resource.PrimaryKey()]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			encodeJsonError(w, "primery key is required")
			return
		}

		if base.Resource.UpdatedAtField().Valid {
			data[base.Resource.UpdatedAtField().String] = time.Now()
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
		validation, err := base.Resource.ValidateFields(base.Validate, data)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(validation) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			encodeJsonError(w, fmt.Sprintf("%s", validation))
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
