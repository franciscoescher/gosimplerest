package gosimplerest

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/franciscoescher/gosimplerest/handlers"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
)

func AddGinHandlers(r *gin.Engine, d *sql.DB, l *logrus.Logger, resources []resource.Resource, mid ...gin.HandlerFunc) {

	for i := range resources {
		base := &resource.Base{Logger: l, DB: d, Resource: &resources[i]}
		name := fmt.Sprintf("/%s", strcase.KebabCase(resources[i].Table))
		nameID := fmt.Sprintf("%s/:id", name)
		r.POST(name, GinGetHandlersChain(handlers.CreateHandler(base), mid...)...)
		r.GET(nameID, GinGetHandlersChain(handlers.RetrieveHandler(base), mid...)...)
		r.PUT(name, GinGetHandlersChain(handlers.UpdateHandler(base), mid...)...)
		r.DELETE(nameID, GinGetHandlersChain(handlers.DeleteHandler(base), mid...)...)
		r.GET(name, GinGetHandlersChain(handlers.SearchHandler(base), mid...)...)

		for _, belongsTo := range resources[i].BelongsToFields {
			nameBelongsTo := fmt.Sprintf("/%s/:id%s", strcase.KebabCase(belongsTo.Table), name)
			r.GET(nameBelongsTo, GinGetHandlersChain(handlers.GetBelongsToHandler(base, belongsTo), mid...)...)
		}
	}
}

func GinGetHandlersChain(h http.HandlerFunc, mid ...gin.HandlerFunc) []gin.HandlerFunc {
	handlers := make([]gin.HandlerFunc, 0)
	handlers = append(handlers, mid...)
	handlers = append(handlers, ConvertHttpHandlerToGinHandler(h))
	return handlers
}

// ConvertHttpHandlerToGinHandler converts a http.HandlerFunc to a gin.HandlerFunc
// It adds params to request context.
func ConvertHttpHandlerToGinHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := make(map[string]string, 0)
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}
		r := handlers.GetRequestWithParams(c.Request, params)
		h(c.Writer, r)
	}
}
