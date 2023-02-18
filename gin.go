package gosimplerest

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func AddGinHandlers(r *gin.Engine, d *sql.DB, l *logrus.Logger, v *validator.Validate, resources []resource.Resource) *gin.Engine {
	h := AddRouteFunctions{
		Post:   GinAddRouteFunc(r.POST),
		Get:    GinAddRouteFunc(r.GET),
		Put:    GinAddRouteFunc(r.PUT),
		Patch:  GinAddRouteFunc(r.PATCH),
		Delete: GinAddRouteFunc(r.DELETE),
		Head:   GinAddRouteFunc(r.HEAD),
	}
	apf := func(name string, param string) string {
		return fmt.Sprintf("%s/:%s", name, param)
	}
	AddHandlers(d, l, v, h, apf, resources)
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
