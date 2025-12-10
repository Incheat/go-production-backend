package chimiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestUnitZapRecovery_PanicHandled(t *testing.T) {
	// Capture logs using zaptest observer
	core, recorded := observer.New(zap.ErrorLevel)
	logger := zap.New(core)

	// handler that panics
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("boom!")
	})

	handler := ZapRecovery(logger)(next)

	req := httptest.NewRequest(http.MethodGet, "/panic-test", nil)
	rr := httptest.NewRecorder()

	// The middleware should recover — this must NOT panic
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("middleware should recover panic, but panic escaped: %#v", rec)
			}
		}()
		handler.ServeHTTP(rr, req)
	}()

	// Assert response code
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}

	// Assert 1 log entry
	entries := recorded.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(entries))
	}

	entry := entries[0]

	// Message check
	if entry.Message != "panic recovered" {
		t.Fatalf("expected log message %q, got %q", "panic recovered", entry.Message)
	}

	fields := entry.ContextMap()

	// error field check
	if fields["error"] != "boom!" {
		t.Fatalf(`expected error field "boom!", got %#v`, fields["error"])
	}

	// path field check
	if fields["path"] != "/panic-test" {
		t.Fatalf(`expected path "/panic-test", got %#v`, fields["path"])
	}
}
