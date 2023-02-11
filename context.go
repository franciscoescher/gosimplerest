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
	params := r.Context().Value(ContextKeyParams)
	if params == nil {
		return ""
	}
	values, ok := r.Context().Value(ContextKeyParams).(map[string]string)
	if !ok {
		return ""
	}
	val, ok := values[s]
	if !ok {
		return ""
	}
	return val
}

func GetRequestWithParams(r *http.Request, params map[string]string) *http.Request {
	if r.Context() == nil {
		r = r.WithContext(context.Background())
	}
	return r.WithContext(context.WithValue(r.Context(), ContextKeyParams, params))
}
