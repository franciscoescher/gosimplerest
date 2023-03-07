package handlers

import (
	"net/http"
)

// DeleteHandler returns a handler for the DELETE method
func DeleteHandler(params *GetHandlerFuncParams) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ReadParams(r, "id")

		// validates id
		err := params.Resource.ValidateField(params.Validate, params.Resource.PrimaryKey, id)
		if err != nil {
			params.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = params.Repository.Delete(params.Resource, id)
		if err != nil {
			if err.Error() == "no rows affected" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			params.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
