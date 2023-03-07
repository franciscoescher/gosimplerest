package gosimplerest

import (
	"net/http"
	"strings"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/gin-gonic/gin"
)

func AddGinHandlers(r *gin.Engine, base AddHandlersBaseParams) *gin.Engine {
	params := AddHandlersParams{
		AddHandlersBaseParams: base,
		AddRouteFunctions: AddRouteFunctions{
			Post:   GinAddRouteFunc(r.POST),
			Get:    GinAddRouteFunc(r.GET),
			Put:    GinAddRouteFunc(r.PUT),
			Patch:  GinAddRouteFunc(r.PATCH),
			Delete: GinAddRouteFunc(r.DELETE),
			Head:   GinAddRouteFunc(r.HEAD),
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

// GinAddRouteType is the type of the function that gin.Engine uses to add routes to the router.
// Example: r.POST, r.GET...
type GinAddRouteType func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes

// GinAddRouteFunc uses the f function to add a route to the router,
// wrapping the handler to add params to request context.
func GinAddRouteFunc(f GinAddRouteType) AddRouteFunc {
	return func(name string, h http.HandlerFunc) {
		f(name, GinHandlerWrapper(h))
	}
}

// GinHandlerWrapper converts a http.HandlerFunc to a gin.HandlerFunc
// It adds params to request context.
func GinHandlerWrapper(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := make(map[string]string, 0)
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}
		handlers.AddParamsToHandlerFunc(h, params)(c.Writer, c.Request)
	}
}
