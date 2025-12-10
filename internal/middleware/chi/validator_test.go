package chimiddleware_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	middleware "github.com/incheat/go-playground/internal/middleware/chi"
)

func TestUnitNewValidatorOptions_ProdMode(t *testing.T) {
	var logged string

	cfg := middleware.ValidatorConfig{
		ProdMode:  true,
		ProdError: "invalid request (prod)",
		Logger: func(format string, args ...any) {
			logged = fmt.Sprintf(format, args...)
		},
	}

	opts := middleware.NewValidatorOptions(cfg)
	if opts == nil {
		t.Fatalf("expected options to be non-nil")
	}

	rr := httptest.NewRecorder()

	// Simulate a validation error from oapi-codegen
	opts.ErrorHandler(rr, "detailed dev error message", http.StatusBadRequest)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	if got := rr.Header().Get("Content-Type"); !strings.Contains(got, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q", got)
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if body["error"] != "invalid request (prod)" {
		t.Fatalf("expected error %q, got %q", "invalid request (prod)", body["error"])
	}

	if !strings.Contains(logged, "validation error (400): detailed dev error message") {
		t.Fatalf("expected log to contain detailed message, got %q", logged)
	}
}

func TestUnitNewValidatorOptions_DevMode(t *testing.T) {
	var logged string

	cfg := middleware.ValidatorConfig{
		ProdMode: false,
		Logger: func(format string, args ...any) {
			logged = fmt.Sprintf(format, args...)
		},
		// ProdError is ignored in dev mode
		ProdError: "ignored in dev",
	}

	opts := middleware.NewValidatorOptions(cfg)
	if opts == nil {
		t.Fatalf("expected options to be non-nil")
	}

	rr := httptest.NewRecorder()

	msg := "some detailed validation error"
	opts.ErrorHandler(rr, msg, http.StatusUnprocessableEntity)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, rr.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if body["error"] != msg {
		t.Fatalf("expected error %q, got %q", msg, body["error"])
	}

	if !strings.Contains(logged, "validation error (422): some detailed validation error") {
		t.Fatalf("expected log to contain detailed message, got %q", logged)
	}
}

func TestUnitNewValidatorOptions_Defaults(t *testing.T) {
	// No logger and no prod error: should not panic and should use defaults.
	cfg := middleware.ValidatorConfig{
		ProdMode: true,
		// Logger: nil
		// ProdError: ""
	}

	opts := middleware.NewValidatorOptions(cfg)
	if opts == nil {
		t.Fatalf("expected options to be non-nil")
	}

	rr := httptest.NewRecorder()
	opts.ErrorHandler(rr, "whatever", http.StatusBadRequest)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	// Default ProdError should be "invalid request"
	if body["error"] != "invalid request" {
		t.Fatalf("expected default error %q, got %q", "invalid request", body["error"])
	}
}
