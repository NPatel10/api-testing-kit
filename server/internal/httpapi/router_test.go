package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTemplatesList(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/templates", nil)
	rr := httptest.NewRecorder()

	NewRouter(RouterDeps{}).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content type, got %q", got)
	}

	var payload struct {
		Templates []struct {
			Slug string `json:"slug"`
		} `json:"templates"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	if len(payload.Templates) < 4 {
		t.Fatalf("expected at least 4 templates, got %d", len(payload.Templates))
	}

	if payload.Templates[0].Slug == "" {
		t.Fatalf("expected first template to include a slug")
	}
}

func TestTemplatesDetail(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/templates/weather-demo", nil)
	rr := httptest.NewRecorder()

	NewRouter(RouterDeps{}).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var payload struct {
		Slug     string `json:"slug"`
		Category string `json:"category"`
		Request  struct {
			Method string `json:"method"`
			URL    string `json:"url"`
		} `json:"request"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	if payload.Slug != "weather-demo" {
		t.Fatalf("expected slug %q, got %q", "weather-demo", payload.Slug)
	}

	if payload.Category != "CRUD examples" {
		t.Fatalf("expected category %q, got %q", "CRUD examples", payload.Category)
	}

	if payload.Request.Method != http.MethodGet {
		t.Fatalf("expected request method %q, got %q", http.MethodGet, payload.Request.Method)
	}

	if payload.Request.URL == "" {
		t.Fatalf("expected request URL to be populated")
	}
}

func TestTemplatesDetailNotFound(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/templates/does-not-exist", nil)
	rr := httptest.NewRecorder()

	NewRouter(RouterDeps{}).ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	var payload struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	if payload.Error.Code != "template_not_found" {
		t.Fatalf("expected error code %q, got %q", "template_not_found", payload.Error.Code)
	}
}

func TestGuestRunsRouteIsRegistered(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/guest-runs", strings.NewReader("{"))
	rr := httptest.NewRecorder()

	NewRouter(RouterDeps{}).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "invalid_request") {
		t.Fatalf("expected guest run error payload, got %s", rr.Body.String())
	}
}
