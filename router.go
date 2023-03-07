package gosimplerest

import (
	"net/http"
	"strings"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/franciscoescher/gosimplerest/interfaces"
	"github.com/franciscoescher/gosimplerest/repository"
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

type AddHandlersBaseParams struct {
	Logger      interfaces.Logger
	Validator   interfaces.Validator
	Resources   []resource.Resource
	Respository repository.RepositoryInterface
}

type AddHandlersParams struct {
	AddHandlersBaseParams
	AddRouteFunctions AddRouteFunctions
	AddParamFunc      AddParamFunc
}

// AddHandlers adds the routes to the router
func AddHandlers(params AddHandlersParams) {
	if params.Validator == nil {
		params.Validator = validator.New()
	}
	if params.Logger == nil {
		params.Logger = logrus.New()
	}
	for i := range params.Resources {
		p := &handlers.GetHandlerFuncParams{
			Logger:     params.Logger,
			Validate:   params.Validator,
			Resource:   &params.Resources[i],
			Repository: params.Respository,
		}
		var sb strings.Builder
		sb.WriteString("/")
		sb.WriteString(strcase.KebabCase(params.Resources[i].Table()))
		name := sb.String()
		nameID := params.AddParamFunc(name, "id")

		h := params.AddRouteFunctions
		if !params.Resources[i].OmitCreateRoute {
			h.Post(name, handlers.CreateHandler(p))
		}
		if !params.Resources[i].OmitRetrieveRoute {
			h.Get(nameID, handlers.RetrieveHandler(p))
		}
		if !params.Resources[i].OmitUpdateRoute {
			h.Put(name, handlers.UpdateHandler(p))
		}
		if !params.Resources[i].OmitPartialUpdateRoute {
			h.Patch(name, handlers.UpdateHandler(p))
		}
		if !params.Resources[i].OmitDeleteRoute {
			h.Delete(nameID, handlers.DeleteHandler(p))
		}
		if !params.Resources[i].OmitSearchRoute {
			h.Get(name, handlers.SearchHandler(p))
		}
		if !params.Resources[i].OmitHeadRoutes {
			h.Head(nameID, handlers.RetrieveHandler(p))
			h.Head(name, handlers.SearchHandler(p))
		}
	}
}
