package httpapi

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"api-testing-kit/server/internal/guest"
	"api-testing-kit/server/internal/safety"
)

type GuestRunsHandler struct {
	service *guest.Service
}

func NewGuestRunsHandler(service *guest.Service) *GuestRunsHandler {
	return &GuestRunsHandler{service: service}
}

func (h *GuestRunsHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/guest-runs", h.handleRun)
}

func (h *GuestRunsHandler) handleRun(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "guest_unavailable", "guest execution is temporarily unavailable")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)

	var payload guest.RunRequest
	if err := decodeJSON(r, &payload); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeError(w, http.StatusRequestEntityTooLarge, "guest_request_too_large", "guest request is too large")
			return
		}
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		return
	}

	result, err := h.service.Execute(r.Context(), clientIPFromRequest(r), payload)
	if err != nil {
		var guestErr *guest.GuestError
		if errors.As(err, &guestErr) {
			writeError(w, guestErr.Status, string(guestErr.Code), guestErr.Message)
			return
		}
		var validationErr *safety.ValidationError
		if errors.As(err, &validationErr) {
			writeError(w, http.StatusForbidden, "blocked_target", validationErr.Message)
			return
		}
		writeError(w, http.StatusBadGateway, "guest_request_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func clientIPFromRequest(r *http.Request) string {
	if r == nil {
		return ""
	}

	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			return normalizeIP(parts[0])
		}
	}

	return normalizeIP(r.RemoteAddr)
}

func normalizeIP(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	if host, _, err := net.SplitHostPort(value); err == nil {
		value = host
	}

	return strings.TrimSpace(value)
}
