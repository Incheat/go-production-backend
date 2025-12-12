package chimiddleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestUnitZapLogger_WithRequestID(t *testing.T) {
	// Arrange: set up zap observer
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	// Simple handler that sets a status and writes a body
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated) // 201
		_, _ = w.Write([]byte("ok"))
	})

	// Wrap handler with ZapLogger
	handler := ZapLogger(logger)(next)

	// Build request with a request ID in context
	req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextRequestIDKey, "req-123"))

	rr := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rr, req)

	// Assert: response status
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	// Assert: we logged exactly one entry
	entries := recorded.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(entries))
	}

	entry := entries[0]

	if entry.Message != "request handled" {
		t.Fatalf("expected log message %q, got %q", "request handled", entry.Message)
	}

	fields := entry.ContextMap()

	// Check core fields
	if fields["request_id"] != "req-123" {
		t.Fatalf("expected req-123, got %#v", fields["request_id"])
	}

	if fields["status"] != int64(201) {
		t.Fatalf("expected 201, got %#v", fields["status"])
	}
	if fields["path"] != "/test/path" {
		t.Fatalf("expected /test/path, got %#v", fields["path"])
	}
	if fields["latency"] == nil {
		t.Fatalf("expected latency field to be set")
	}
}

func TestUnitZapLogger_WithoutRequestID(t *testing.T) {
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// no explicit status -> defaults to 200
		_, _ = w.Write([]byte("ok"))
	})

	handler := ZapLogger(logger)(next)

	req := httptest.NewRequest(http.MethodGet, "/no-id", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	entries := recorded.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(entries))
	}

	entry := entries[0]
	fields := entry.ContextMap()

	// When there is no request id in context, we expect "-"
	if fields["request_id"] != "-" {
		t.Fatalf("expected request_id %q when missing, got %#v", "-", fields["request_id"])
	}
	if fields["path"] != "/no-id" {
		t.Fatalf("expected path %q, got %#v", "/no-id", fields["path"])
	}
}
