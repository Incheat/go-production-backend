package chimiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

// Test that when no X-Request-ID header is provided, the middleware
// generates a UUID, sets it in the response header, and stores it in context.
func TestUnitRequestID_GeneratesUUIDWhenMissing(t *testing.T) {
	// Handler to capture the request ID from context.
	var gotReqID string
	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		v := r.Context().Value(ContextRequestIDKey)
		if v != nil {
			if s, ok := v.(string); ok {
				gotReqID = s
			}
		}
	})

	// Wrap the handler with the middleware.
	handler := RequestID()(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Assert header is set.
	reqIDHeader := rr.Header().Get(HeaderRequestID)
	if reqIDHeader == "" {
		t.Fatalf("expected %s header to be set, got empty string", HeaderRequestID)
	}

	// Assert context has the same ID.
	if gotReqID == "" {
		t.Fatalf("expected request ID in context, got empty string")
	}
	if gotReqID != reqIDHeader {
		t.Fatalf("expected context request ID %q to match header %q", gotReqID, reqIDHeader)
	}

	// Assert it looks like a UUID.
	if _, err := uuid.Parse(reqIDHeader); err != nil {
		t.Fatalf("expected generated request ID to be a valid UUID, got %q: %v", reqIDHeader, err)
	}
}

// Test that when an X-Request-ID header is provided, the middleware
// preserves it (does not overwrite), and stores it in context.
func TestUnitRequestID_PreservesExistingID(t *testing.T) {
	const existingID = "my-custom-request-id"

	var gotReqID string
	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		v := r.Context().Value(ContextRequestIDKey)
		if v != nil {
			if s, ok := v.(string); ok {
				gotReqID = s
			}
		}
	})

	handler := RequestID()(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(HeaderRequestID, existingID)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Header should be preserved.
	reqIDHeader := rr.Header().Get(HeaderRequestID)
	if reqIDHeader != existingID {
		t.Fatalf("expected %s header %q, got %q", HeaderRequestID, existingID, reqIDHeader)
	}

	// Context should match the existing ID.
	if gotReqID != existingID {
		t.Fatalf("expected context request ID %q, got %q", existingID, gotReqID)
	}
}
