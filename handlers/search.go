package handlers

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/sirupsen/logrus"
)

// SearchHandler returns a handler for the GET method with query params
func SearchHandler(base *resource.Base) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer cleanBodyIfHead(r, w)

		query := r.URL.Query()

		// validates that all fields in data are in the model
		for key := range query {
			// validates fields
			if !base.Resource.IsSearchable(key) {
				w.WriteHeader(http.StatusBadRequest)
				encodeJsonError(w, fmt.Sprintf("%s is not searchable", key))
				return
			}
			// validates values
			for _, v := range query[key] {
				err := base.Resource.ValidateField(base.Validate, key, v)
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
