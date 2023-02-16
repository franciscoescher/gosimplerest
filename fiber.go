package gosimplerest

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
)

func AddFiberHandlers(r *fiber.App, d *sql.DB, l *logrus.Logger, v *validator.Validate, resources []resource.Resource) {
	if v == nil {
		v = validator.New()
	}
	if l == nil {
		l = logrus.New()
		l.Out = io.Discard
	}
	for i := range resources {
		base := &resource.Base{Logger: l, DB: d, Validate: v, Resource: &resources[i]}
		name := fmt.Sprintf("/%s", strcase.KebabCase(resources[i].Table))
		nameID := fmt.Sprintf("%s/:id", name)

		if !resources[i].OmitCreateRoute {
			r.Post(name, FiberHandler(handlers.CreateHandler(base)))
		}
		if !resources[i].OmitRetrieveRoute {
			r.Get(nameID, FiberHandler(handlers.RetrieveHandler(base)))
		}
		if !resources[i].OmitUpdateRoute {
			r.Put(name, FiberHandler(handlers.UpdateHandler(base)))
		}
		if !resources[i].OmitPartialUpdateRoute {
			r.Patch(name, FiberHandler(handlers.UpdateHandler(base)))
		}
		if !resources[i].OmitDeleteRoute {
			r.Delete(nameID, FiberHandler(handlers.DeleteHandler(base)))
		}
		if !resources[i].OmitSearchRoute {
			r.Get(name, FiberHandler(handlers.SearchHandler(base)))
		}
		if !resources[i].OmitHeadRoutes {
			r.Head(nameID, FiberHandler(handlers.RetrieveHandler(base)))
			r.Head(name, FiberHandler(handlers.SearchHandler(base)))
		}
	}
}

// FiberHandler converts a http.HandlerFunc to a fiber.Handler
// It adds params to request context.
func FiberHandler(h http.HandlerFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := make(map[string]string, 0)

		for key, value := range c.AllParams() {
			params[key] = value
		}
		hWithParams := handlers.AddParamsToHandlerFunc(h, params)
		return adaptor.HTTPHandlerFunc(hWithParams)(c)
	}
}
