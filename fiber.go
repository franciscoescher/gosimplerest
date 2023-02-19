package gosimplerest

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func AddFiberHandlers(r *fiber.App, d *sql.DB, l *logrus.Logger, v *validator.Validate, resources []resource.Resource) *fiber.App {
	h := AddRouteFunctions{
		Post:   FiberAddRouteFunc(r.Post),
		Get:    FiberAddRouteFunc(r.Get),
		Put:    FiberAddRouteFunc(r.Put),
		Patch:  FiberAddRouteFunc(r.Patch),
		Delete: FiberAddRouteFunc(r.Delete),
		Head:   FiberAddRouteFunc(r.Head),
	}
	apf := func(name string, param string) string {
		var sb strings.Builder
		sb.WriteString(name)
		sb.WriteString("/:")
		sb.WriteString(param)
		return sb.String()
	}
	AddHandlers(d, l, v, h, apf, resources)
	return r
}

// FiberAddRouteType is the type of the function that fiber.App uses to add routes to the router.
// Example: r.Post, r.Get...
type FiberAddRouteType func(path string, handlers ...fiber.Handler) fiber.Router

// FiberAddRouteFunc uses the f function to add a route to the router,
// wrapping the handler to add params to request context.
func FiberAddRouteFunc(f FiberAddRouteType) AddRouteFunc {
	return func(name string, h http.HandlerFunc) {
		f(name, FiberHandlerWrapper(h))
	}
}

// FiberHandlerWrapper wraps the handler function to a fiber handler.
// It adds params to request context.
func FiberHandlerWrapper(h http.HandlerFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := make(map[string]string, 0)

		for key, value := range c.AllParams() {
			params[key] = value
		}
		hWithParams := handlers.AddParamsToHandlerFunc(h, params)
		return adaptor.HTTPHandlerFunc(hWithParams)(c)
	}
}
