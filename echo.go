package gosimplerest

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
)

func AddEchoHandlers(r *echo.Echo, d *sql.DB, l *logrus.Logger, v *validator.Validate, resources []resource.Resource) {
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
			r.POST(name, EchoHandler(handlers.CreateHandler(base)))
		}
		if !resources[i].OmitRetrieveRoute {
			r.GET(nameID, EchoHandler(handlers.RetrieveHandler(base)))
		}
		if !resources[i].OmitUpdateRoute {
			r.PUT(name, EchoHandler(handlers.UpdateHandler(base)))
		}
		if !resources[i].OmitPartialUpdateRoute {
			r.PATCH(name, EchoHandler(handlers.UpdateHandler(base)))
		}
		if !resources[i].OmitDeleteRoute {
			r.DELETE(nameID, EchoHandler(handlers.DeleteHandler(base)))
		}
		if !resources[i].OmitSearchRoute {
			r.GET(name, EchoHandler(handlers.SearchHandler(base)))
		}
		if !resources[i].OmitHeadRoutes {
			r.HEAD(nameID, EchoHandler(handlers.RetrieveHandler(base)))
			r.HEAD(name, EchoHandler(handlers.SearchHandler(base)))
		}
		if !resources[i].OmitBelongsToRoutes {
			for _, belongsTo := range resources[i].BelongsToFields {
				nameBelongsTo := fmt.Sprintf("/%s/:id%s", strcase.KebabCase(belongsTo.Table), name)
				r.GET(nameBelongsTo, EchoHandler(handlers.GetBelongsToHandler(base, belongsTo)))
			}
		}
	}
}

// EchoHandler converts a http.HandlerFunc to a echo.HandlerFunc
// It adds params to request context.
func EchoHandler(h http.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := make(map[string]string, 0)

		for _, param := range c.ParamNames() {
			params[param] = c.Param(param)
		}
		handlers.AddParamsToHandlerFunc(h, params)(c.Response().Writer, c.Request())
		return nil
	}
}
