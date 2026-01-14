package chimiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	chimiddlewareutils "github.com/incheat/go-production-backend/services/auth/internal/middleware/chi/utils"
)

func TestRequestMeta_PopulatesContextFromHeaderAndRequest(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta, ok := chimiddlewareutils.GetRequestMeta(r.Context())
		if !ok {
			t.Fatal("expected request meta in context")
		}

		if meta.RequestID != "req-123" {
			t.Fatalf("expected RequestID %q, got %q", "req-123", meta.RequestID)
		}
		if meta.UserAgent != "test-agent/1.0" {
			t.Fatalf("expected UserAgent %q, got %q", "test-agent/1.0", meta.UserAgent)
		}
		if meta.IPAddress != "203.0.113.10" {
			t.Fatalf("expected IPAddress %q, got %q", "203.0.113.10", meta.IPAddress)
		}

		w.WriteHeader(http.StatusOK)
	})

	handler := RequestMeta()(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(HeaderRequestID, "req-123")
	req.Header.Set("User-Agent", "test-agent/1.0")
	req.Header.Set("X-Forwarded-For", "203.0.113.10")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestRequestMeta_IPPrefersXForwardedForOverXRealIPAndRemoteAddr(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta, ok := chimiddlewareutils.GetRequestMeta(r.Context())
		if !ok {
			t.Fatal("expected request meta in context")
		}
		if meta.IPAddress != "198.51.100.1" {
			t.Fatalf("expected IPAddress %q, got %q", "198.51.100.1", meta.IPAddress)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestMeta()(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Forwarded-For", "198.51.100.1")
	req.Header.Set("X-Real-IP", "198.51.100.2") // should be ignored when XFF exists
	req.RemoteAddr = "192.0.2.9:54321"          // should be ignored when XFF exists

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestRequestMeta_IPFallsBackToXRealIPWhenNoXForwardedFor(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta, ok := chimiddlewareutils.GetRequestMeta(r.Context())
		if !ok {
			t.Fatal("expected request meta in context")
		}
		if meta.IPAddress != "198.51.100.2" {
			t.Fatalf("expected IPAddress %q, got %q", "198.51.100.2", meta.IPAddress)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestMeta()(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Real-IP", "198.51.100.2")
	req.RemoteAddr = "192.0.2.9:54321" // should be ignored when X-Real-IP exists

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestRequestMeta_IPFallsBackToRemoteAddrAndStripsPort(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta, ok := chimiddlewareutils.GetRequestMeta(r.Context())
		if !ok {
			t.Fatal("expected request meta in context")
		}
		if meta.IPAddress != "192.0.2.9" {
			t.Fatalf("expected IPAddress %q, got %q", "192.0.2.9", meta.IPAddress)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestMeta()(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.0.2.9:54321"

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestRequestMeta_IPUsesRemoteAddrAsIsWhenSplitHostPortFails(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta, ok := chimiddlewareutils.GetRequestMeta(r.Context())
		if !ok {
			t.Fatal("expected request meta in context")
		}
		if meta.IPAddress != "not-a-hostport" {
			t.Fatalf("expected IPAddress %q, got %q", "not-a-hostport", meta.IPAddress)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestMeta()(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "not-a-hostport"

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
