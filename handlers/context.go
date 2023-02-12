package handlers

import (
	"context"
	"net/http"
)

type GinContextKey int

const (
	ContextKeyParams GinContextKey = iota
)

// GetRequestWithParams returns a new request with the given params in the request context
// The params are stored with the key defined by ContextKeyParams
func GetRequestWithParams(r *http.Request, params map[string]string) *http.Request {
	if r.Context() == nil {
		r = r.WithContext(context.Background())
	}
	return r.WithContext(context.WithValue(r.Context(), ContextKeyParams, params))
}

// ReadParams returns the value of the given params previously written in
// the request context with the key defined by ContextKeyParams
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
