package gosimplerest

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
)

func AddGinHandlers(r *gin.Engine, d *sql.DB, l *logrus.Logger, resources []Resource, mid ...gin.HandlerFunc) {
	db = d
	logger = l

	for _, resource := range resources {
		name := fmt.Sprintf("/%s", strcase.KebabCase(resource.Table))
		nameID := fmt.Sprintf("%s/:id", name)
		r.GET(nameID, GinGetHandlersChain(GetHandler(resource), mid...)...)
		r.DELETE(nameID, GinGetHandlersChain(DeleteHandler(resource), mid...)...)
		r.POST(name, GinGetHandlersChain(CreateHandler(resource), mid...)...)
		r.PUT(name, GinGetHandlersChain(UpdateHandler(resource), mid...)...)
		r.GET(name, GinGetHandlersChain(SearchHandler(resource), mid...)...)

		for _, belongsTo := range resource.BelongsToFields {
			nameBelongsTo := fmt.Sprintf("/%s/:id%s", strcase.KebabCase(belongsTo.Table), name)
			r.GET(nameBelongsTo, GinGetHandlersChain(GetBelongsToHandler(resource, belongsTo), mid...)...)
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
		r := GetRequestWithParams(c.Request, params)
		h(c.Writer, r)
	}
}
