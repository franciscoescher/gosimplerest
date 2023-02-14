package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/franciscoescher/gosimplerest/resource"
)

// GetBelongsToHandler returns a handler for the GET method of the belongs to relationship
func GetBelongsToHandler(base *resource.Base, belongsTo resource.BelongsTo) http.HandlerFunc {
	if base.Resource.OmitBelongsToRoutes {
		return NotFoundHandler
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := ReadParams(r, "id")

		// validates id
		err := base.Resource.ValidateField(base.Validate, base.Resource.PrimaryKey(), id)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := base.Resource.FindFromBelongsTo(base, id, belongsTo)
		if err != nil {
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(result) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}
