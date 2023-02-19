package gosimplerest

import (
	"database/sql"
	"io"
	"net/http"
	"strings"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
)

// AddRouteFunc is type of a function that adds a route
// to a router, with a given name and handler.
// It should also add the params to the handler.
// Each router/framework should implement this function.
type AddRouteFunc func(name string, h http.HandlerFunc)

// AddParamFunc is a function that adds a param to a url/endpoint.
// It should return the route with the param added.
// Example: /users + id -> /users/:id
// Each router/framework should implement this function.
type AddParamFunc func(name string, param string) string

// AddRouteFunctions is a struct that contains the functions
// to add routes to a router, one for each request method.
// Each router/framework should implement this struct.
type AddRouteFunctions struct {
	Post   AddRouteFunc
	Get    AddRouteFunc
	Put    AddRouteFunc
	Patch  AddRouteFunc
	Delete AddRouteFunc
	Head   AddRouteFunc
}

// AddHandlers adds the routes to the router
func AddHandlers(d *sql.DB, l *logrus.Logger, v *validator.Validate, h AddRouteFunctions, apf AddParamFunc, resources []resource.Resource) {
	if v == nil {
		v = validator.New()
	}
	if l == nil {
		l = logrus.New()
		l.Out = io.Discard
	}
	for i := range resources {
		base := &resource.Base{Logger: l, DB: d, Validate: v, Resource: &resources[i]}
		var sb strings.Builder
		sb.WriteString("/")
		sb.WriteString(strcase.KebabCase(resources[i].Table))
		name := sb.String()
		nameID := apf(name, "id")

		if !resources[i].OmitCreateRoute {
			h.Post(name, handlers.CreateHandler(base))
		}
		if !resources[i].OmitRetrieveRoute {
			h.Get(nameID, handlers.RetrieveHandler(base))
		}
		if !resources[i].OmitUpdateRoute {
			h.Put(name, handlers.UpdateHandler(base))
		}
		if !resources[i].OmitPartialUpdateRoute {
			h.Patch(name, handlers.UpdateHandler(base))
		}
		if !resources[i].OmitDeleteRoute {
			h.Delete(nameID, handlers.DeleteHandler(base))
		}
		if !resources[i].OmitSearchRoute {
			h.Get(name, handlers.SearchHandler(base))
		}
		if !resources[i].OmitHeadRoutes {
			h.Head(nameID, handlers.RetrieveHandler(base))
			h.Head(name, handlers.SearchHandler(base))
		}
	}
}
