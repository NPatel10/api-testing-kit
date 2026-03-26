package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api-testing-kit/server/internal/billing"
)

func TestBillingCheckoutHooks(t *testing.T) {
	t.Parallel()

	handler := NewBillingHandler(nil)
	mux := http.NewServeMux()
	handler.Register(mux)

	successReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/checkout/success", nil)
	successRR := httptest.NewRecorder()
	mux.ServeHTTP(successRR, successReq)

	if successRR.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected unavailable status when service is nil, got %d", successRR.Code)
	}

	runtimeHandler := NewBillingHandler(billing.NewStubService())
	runtimeMux := http.NewServeMux()
	runtimeHandler.Register(runtimeMux)

	successRR = httptest.NewRecorder()
	runtimeMux.ServeHTTP(successRR, successReq)
	if successRR.Code != http.StatusOK {
		t.Fatalf("expected success hook status %d, got %d", http.StatusOK, successRR.Code)
	}

	var successPayload struct {
		Status       string `json:"status"`
		Message      string `json:"message"`
		NextAction   string `json:"nextAction"`
		ProviderNote string `json:"providerNote"`
	}
	if err := json.Unmarshal(successRR.Body.Bytes(), &successPayload); err != nil {
		t.Fatalf("failed to decode checkout success payload: %v", err)
	}
	if successPayload.Status != "checkout_success_stub" || successPayload.ProviderNote == "" {
		t.Fatalf("unexpected checkout success payload: %+v", successPayload)
	}

	cancelReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/checkout/cancel", nil)
	cancelRR := httptest.NewRecorder()
	runtimeMux.ServeHTTP(cancelRR, cancelReq)
	if cancelRR.Code != http.StatusOK {
		t.Fatalf("expected cancel hook status %d, got %d", http.StatusOK, cancelRR.Code)
	}
}

func TestBillingWebhookAcceptsStubPayload(t *testing.T) {
	t.Parallel()

	handler := NewBillingHandler(billing.NewStubService())
	mux := http.NewServeMux()
	handler.Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/webhooks", strings.NewReader(`{"type":"subscription.updated"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected accepted status, got %d", rr.Code)
	}

	var payload struct {
		Status       string `json:"status"`
		Message      string `json:"message"`
		ProviderNote string `json:"providerNote"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode webhook payload: %v", err)
	}
	if payload.Status != "webhook_received_stub" || payload.ProviderNote == "" {
		t.Fatalf("unexpected webhook payload: %+v", payload)
	}
}
