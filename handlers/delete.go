package handlers

import (
	"net/http"

	"github.com/franciscoescher/gosimplerest/resource"
)

// DeleteHandler returns a handler for the DELETE method
func DeleteHandler(base *resource.Base) http.HandlerFunc {
	if base.Resource.OmitDeleteRoute {
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

		err = base.Resource.Delete(base, id)
		if err != nil {
			if err.Error() == "no rows affected" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			base.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
