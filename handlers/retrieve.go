package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/franciscoescher/gosimplerest/resource"
)

// RetrieveHandler returns a handler for the GET method
func RetrieveHandler(base *resource.Base) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer cleanBodyIfHead(r, w)

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

		json.NewEncoder(w).Encode(result)
	}
}

func cleanBodyIfHead(r *http.Request, w http.ResponseWriter) {
	if r.Method == http.MethodHead {
		w.Write([]byte{})
	}
}
