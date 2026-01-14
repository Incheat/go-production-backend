package chimiddlewareutils

import "context"

type requestMetaKey struct{}

// RequestMeta is the metadata for the request.
type RequestMeta struct {
	RequestID string
	UserAgent string
	IPAddress string
	// Additional metadata: Referer, AcceptLanguage, etc.
}

// WithRequestMeta adds the request metadata to the context.
func WithRequestMeta(ctx context.Context, meta RequestMeta) context.Context {
	return context.WithValue(ctx, requestMetaKey{}, meta)
}

// GetRequestMeta gets the request metadata from the context.
func GetRequestMeta(ctx context.Context) (RequestMeta, bool) {
	meta, ok := ctx.Value(requestMetaKey{}).(RequestMeta)
	return meta, ok
}
