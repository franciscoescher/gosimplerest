package handlers

import (
	"fmt"
	"net/http"

	"time"
)

// UpdateHandler returns a handler for the PATCH method
func UpdateHandler(params *GetHandlerFuncParams) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := unmarshalBody(r)
		if err != nil {
			params.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// primary key is required
		_, ok := data[params.Resource.PrimaryKey]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			err = encodeJsonError(w, r, "primery key is required")
			if err != nil {
				params.Logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		// Checks for immutable fields being updated and adds missing fields if method is PUT
		for key, field := range params.Resource.Fields {
			_, ok := data[key]
			// checks for tentative of updating immutable fields
			if ok && field.Immutable {
				w.WriteHeader(http.StatusBadRequest)
				err = encodeJsonError(w, r, key+" is immutable")
				if err != nil {
					params.Logger.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			} else if r.Method == http.MethodPut && !ok && !field.Immutable {
				// if method is PUT and field is not immutable and not present in the request,
				// adds it to the data for update
				data[key] = nil
			}
		}

		if params.Resource.UpdatedAtField.Valid {
			data[params.Resource.UpdatedAtField.String] = time.Now()
		}

		// perform data validation
		for key := range data {
			// validates field exists in the model
			if !params.Resource.HasField(key) {
				w.WriteHeader(http.StatusBadRequest)
				err = encodeJsonError(w, r, key+" not in the model")
				if err != nil {
					params.Logger.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
		}
		// validates values
		fmt.Println("AQUIIIIIIIIIIIIIIIIIIIIIIIIIIII")
		fmt.Println(data)
		errs := params.Resource.ValidateInputFields(params.Validate, data)
		if len(errs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			err = encodeJsonError(w, r, fmt.Sprintf("%s", errs))
			if err != nil {
				params.Logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		fmt.Println("AQUIIIIIIIIIIIIIIIIIIIIIIIIIIII2222")

		affected, err := params.Repository.Update(params.Resource, data)
		if err != nil {
			params.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !affected {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}
