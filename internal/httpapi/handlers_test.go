package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"config-audit/internal/app"
)

func TestAnalyzeHandler(t *testing.T) {
	handler := NewHandler(app.NewDefaultAnalyzer())

	request := httptest.NewRequest(http.MethodPost, "/v1/analyze?filename=config.json", strings.NewReader(`{"tls":{"enabled":false}}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	if !strings.Contains(response.Body.String(), "tls-disabled") {
		t.Fatalf("expected tls-disabled response, got %q", response.Body.String())
	}
}

func TestAnalyzeHandlerRejectsInvalidConfig(t *testing.T) {
	handler := NewHandler(app.NewDefaultAnalyzer())

	request := httptest.NewRequest(http.MethodPost, "/v1/analyze?filename=config.json", strings.NewReader(`{"tls":`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}

	if !strings.Contains(response.Body.String(), "decode json") {
		t.Fatalf("expected decode error response, got %q", response.Body.String())
	}
}

func TestHealthHandler(t *testing.T) {
	handler := NewHandler(app.NewDefaultAnalyzer())

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if !strings.Contains(response.Body.String(), `"status":"ok"`) {
		t.Fatalf("expected health response, got %q", response.Body.String())
	}
}
