package handlers

import (
	"net/http"

	"github.com/franciscoescher/gosimplerest/resource"
)

// RetrieveHandler returns a handler for the GET method
func RetrieveHandler(base *resource.Base) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ReadParams(r, "id")

		// validates id
		err := base.Resource.ValidateField(base.Validate, base.Resource.PrimaryKey, id)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := base.Resource.Find(base, id)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(result) == 0 {
			w.WriteHeader(http.StatusNotFound)
			encodeJsonError(w, "not found")
			return
		}

		err = encodeJson(w, r, result)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
