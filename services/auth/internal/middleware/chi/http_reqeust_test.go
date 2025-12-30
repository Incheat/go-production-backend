package chimiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	chimiddlewareutils "github.com/incheat/go-production-backend/services/auth/internal/middleware/chi/utils"
)

// Test that HTTPRequest middleware stores the *http.Request in the context
// and that downstream handlers can retrieve it using GetHTTPRequest.
func TestUnitHTTPRequestMiddlewareStoresRequestInContext(t *testing.T) {
	// A handler that reads the request from context using GetHTTPRequest.
	var gotReq *http.Request
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReq = chimiddlewareutils.GetHTTPRequest(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the middleware
	mw := HTTPRequest()
	wrapped := mw(handler)

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/test-path?x=1", nil)
	req.Header.Set("User-Agent", "test-agent")
	rr := httptest.NewRecorder()

	// Serve the request
	wrapped.ServeHTTP(rr, req)

	// Ensure handler was called and request was stored
	if gotReq == nil {
		t.Fatalf("expected request to be stored in context, got nil")
	}

	// Validate some properties of the stored request
	if gotReq.Method != http.MethodGet {
		t.Errorf("expected method %q, got %q", http.MethodGet, gotReq.Method)
	}
	if gotReq.URL.Path != "/test-path" {
		t.Errorf("expected path %q, got %q", "/test-path", gotReq.URL.Path)
	}
	if gotReq.Header.Get("User-Agent") != "test-agent" {
		t.Errorf("expected User-Agent %q, got %q", "test-agent", gotReq.Header.Get("User-Agent"))
	}
}

// Test that GetHTTPRequest returns nil when there is no request in the context.
func TestUnitGetHTTPRequestEmptyContext(t *testing.T) {
	ctx := context.Background()
	req := chimiddlewareutils.GetHTTPRequest(ctx)
	if req != nil {
		t.Fatalf("expected nil when no request stored in context, got %#v", req)
	}
}
