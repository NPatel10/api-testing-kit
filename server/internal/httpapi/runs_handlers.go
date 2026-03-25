package httpapi

import (
	"errors"
	"net/http"

	"api-testing-kit/server/internal/auth"
	"api-testing-kit/server/internal/runner"
	"api-testing-kit/server/internal/safety"
)

type RunsHandler struct {
	runner *runner.Service
	auth   *auth.Service
}

func NewRunsHandler(runnerService *runner.Service, authService *auth.Service) *RunsHandler {
	return &RunsHandler{runner: runnerService, auth: authService}
}

func (h *RunsHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/runs", h.handleRun)
}

func (h *RunsHandler) handleRun(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.runner == nil || h.auth == nil {
		writeError(w, http.StatusServiceUnavailable, "runner_unavailable", "request execution is temporarily unavailable")
		return
	}
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	user, _, err := h.auth.CurrentUser(r.Context(), cookie.Value)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	var payload runner.RunInput
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		return
	}

	result, err := h.runner.Execute(r.Context(), user.ID, payload)
	if err != nil {
		var validationErr *safety.ValidationError
		if errors.As(err, &validationErr) {
			writeError(w, http.StatusForbidden, "blocked_target", validationErr.Message)
			return
		}
		if errors.Is(err, runner.ErrUnavailable) {
			writeError(w, http.StatusServiceUnavailable, "runner_unavailable", "request execution is temporarily unavailable")
			return
		}
		if errors.Is(err, runner.ErrInvalid) {
			writeError(w, http.StatusBadRequest, "invalid_run", "run payload is invalid")
			return
		}
		writeError(w, http.StatusBadGateway, "upstream_request_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}
