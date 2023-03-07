package handlers

import (
	"net/http"
)

// RetrieveHandler returns a handler for the GET method
func RetrieveHandler(params *GetHandlerFuncParams) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ReadParams(r, "id")

		// validates id
		err := params.Resource.ValidateField(params.Validate, params.Resource.PrimaryKey, id)
		if err != nil {
			params.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := params.Repository.Find(params.Resource, id)
		if err != nil {
			params.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(result) == 0 {
			w.WriteHeader(http.StatusNotFound)
			err = encodeJsonError(w, r, "not found")
			if err != nil {
				params.Logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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
