package httpapi

import (
	"encoding/json"
	"io"
	"net/http"

	"api-testing-kit/server/internal/billing"
)

type BillingHandler struct {
	service *billing.StubService
}

func NewBillingHandler(service *billing.StubService) *BillingHandler {
	return &BillingHandler{service: service}
}

func (h *BillingHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/billing/webhooks", h.handleWebhook)
	mux.HandleFunc("GET /api/v1/billing/checkout/success", h.handleCheckoutSuccess)
	mux.HandleFunc("GET /api/v1/billing/checkout/cancel", h.handleCheckoutCancel)
}

func (h *BillingHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "billing_unavailable", "billing integration is temporarily unavailable")
		return
	}

	// Keep the stub permissive: providers can POST their own event shapes later.
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil && err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		return
	}

	writeJSON(w, http.StatusAccepted, h.service.WebhookReceived())
}

func (h *BillingHandler) handleCheckoutSuccess(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "billing_unavailable", "billing integration is temporarily unavailable")
		return
	}

	writeJSON(w, http.StatusOK, h.service.CheckoutSuccess())
}

func (h *BillingHandler) handleCheckoutCancel(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "billing_unavailable", "billing integration is temporarily unavailable")
		return
	}

	writeJSON(w, http.StatusOK, h.service.CheckoutCancel())
}
