package gosimplerest

import (
	"context"
	"net/http"
)

type GinContextKey int

const (
	ContextKeyParams GinContextKey = iota
)

func ReadParams(r *http.Request, s string) string {
	return r.Context().Value(ContextKeyParams).(map[string]string)[s]
}

func GetRequestWithParams(r *http.Request, params map[string]string) *http.Request {
	if r.Context() == nil {
		r = r.WithContext(context.Background())
	}
	return r.WithContext(context.WithValue(r.Context(), ContextKeyParams, params))
}
