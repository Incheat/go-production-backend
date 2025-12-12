// Package chimiddlewareutils defines the context for the chi middleware utils.
package chimiddlewareutils

import (
	"context"
	"net/http"
)

type httpRequestKey struct{}

// WithHTTPRequest adds a http request to the context.
func WithHTTPRequest(ctx context.Context, request *http.Request) context.Context {
	return context.WithValue(ctx, httpRequestKey{}, request)
}

// GetHTTPRequest gets the http request from the context.
func GetHTTPRequest(ctx context.Context) *http.Request {
	request, ok := ctx.Value(httpRequestKey{}).(*http.Request)
	if !ok {
		// fallback to global http request
		return nil
	}
	return request
}
