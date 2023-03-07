package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// SearchHandler returns a handler for the GET method with query params
func SearchHandler(params *GetHandlerFuncParams) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// validates that all fields in data are in the model
		for key := range query {
			// validates fields
			if !params.Resource.IsSearchable(key) {
				w.WriteHeader(http.StatusBadRequest)
				err := encodeJsonError(w, r, key+" is not searchable")
				if err != nil {
					params.Logger.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
			// validates values
			for _, v := range query[key] {
				err := params.Resource.ValidateField(params.Validate, key, v)
				if err != nil {
					logrus.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
		}

		result, err := params.Repository.Search(params.Resource, query)
		if err != nil {
			params.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(result) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		err = encodeJson(w, r, result)
		if err != nil {
			params.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
