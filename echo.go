package gosimplerest

import (
	"net/http"
	"strings"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/labstack/echo/v4"
)

func AddEchoHandlers(r *echo.Echo, base AddHandlersBaseParams) *echo.Echo {
	params := AddHandlersParams{
		AddHandlersBaseParams: base,
		AddRouteFunctions: AddRouteFunctions{
			Post:   EchoAddRouteFunc(r.POST),
			Get:    EchoAddRouteFunc(r.GET),
			Put:    EchoAddRouteFunc(r.PUT),
			Patch:  EchoAddRouteFunc(r.PATCH),
			Delete: EchoAddRouteFunc(r.DELETE),
			Head:   EchoAddRouteFunc(r.HEAD),
		},
		AddParamFunc: func(name string, param string) string {
			var sb strings.Builder
			sb.WriteString(name)
			sb.WriteString("/:")
			sb.WriteString(param)
			return sb.String()
		},
	}
	AddHandlers(params)
	return r
}

// EchoAddRouteType is the type of the function that echo.Echo uses to add routes to the router.
// Example: r.POST, r.GET...
type EchoAddRouteType func(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route

// EchoAddRouteFunc uses the f function to add a route to the router,
// wrapping the handler to add params to request context.
func EchoAddRouteFunc(f EchoAddRouteType) AddRouteFunc {
	return func(name string, h http.HandlerFunc) {
		f(name, EchoHandlerWrapper(h))
	}
}

// EchoHandlerWrapper converts a http.HandlerFunc to a echo.HandlerFunc
// It adds params to request context.
func EchoHandlerWrapper(h http.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := make(map[string]string, 0)

		for _, param := range c.ParamNames() {
			params[param] = c.Param(param)
		}
		handlers.AddParamsToHandlerFunc(h, params)(c.Response().Writer, c.Request())
		return nil
	}
}
