package gosimplerest

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
)

func AddGinHandlers(r *gin.Engine, d *sql.DB, l *logrus.Logger, v *validator.Validate, resources []resource.Resource) {
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
			r.POST(name, GinHandler(handlers.CreateHandler(base)))
		}
		if !resources[i].OmitRetrieveRoute {
			r.GET(nameID, GinHandler(handlers.RetrieveHandler(base)))
		}
		if !resources[i].OmitUpdateRoute {
			r.PUT(name, GinHandler(handlers.UpdateHandler(base)))
		}
		if !resources[i].OmitDeleteRoute {
			r.DELETE(nameID, GinHandler(handlers.DeleteHandler(base)))
		}
		if !resources[i].OmitSearchRoute {
			r.GET(name, GinHandler(handlers.SearchHandler(base)))
		}
		if !resources[i].OmitBelongsToRoutes {
			for _, belongsTo := range resources[i].BelongsToFields {
				nameBelongsTo := fmt.Sprintf("/%s/:id%s", strcase.KebabCase(belongsTo.Table), name)
				r.GET(nameBelongsTo, GinHandler(handlers.GetBelongsToHandler(base, belongsTo)))
			}
		}
	}
}

// GinHandler converts a http.HandlerFunc to a gin.HandlerFunc
// It adds params to request context.
func GinHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := make(map[string]string, 0)
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}
		handlers.AddParamsToHandlerFunc(h, params)(c.Writer, c.Request)
	}
}
