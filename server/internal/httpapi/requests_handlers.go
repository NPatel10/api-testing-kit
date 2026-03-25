package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"api-testing-kit/server/internal/auth"
	"api-testing-kit/server/internal/history"
	"api-testing-kit/server/internal/requests"
)

type RequestsHandler struct {
	requests *requests.Service
	history  *history.Service
	auth     *auth.Service
}

type savedRequestPayload struct {
	CollectionID    *string         `json:"collectionId,omitempty"`
	Name            *string         `json:"name,omitempty"`
	Description     *string         `json:"description,omitempty"`
	Method          *string         `json:"method,omitempty"`
	URL             *string         `json:"url,omitempty"`
	QueryParams     json.RawMessage `json:"queryParams,omitempty"`
	Headers         json.RawMessage `json:"headers,omitempty"`
	AuthScheme      *string         `json:"authScheme,omitempty"`
	AuthConfig      json.RawMessage `json:"authConfig,omitempty"`
	BodyMode        *string         `json:"bodyMode,omitempty"`
	BodyConfig      json.RawMessage `json:"bodyConfig,omitempty"`
	ExampleResponse json.RawMessage `json:"exampleResponse,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
}

func NewRequestsHandler(requestsService *requests.Service, historyService *history.Service, authService *auth.Service) *RequestsHandler {
	return &RequestsHandler{
		requests: requestsService,
		history:  historyService,
		auth:     authService,
	}
}

func (h *RequestsHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/requests/{id}", h.handleGetRequest)
	mux.HandleFunc("POST /api/v1/requests", h.handleCreateRequest)
	mux.HandleFunc("PATCH /api/v1/requests/{id}", h.handleUpdateRequest)
	mux.HandleFunc("DELETE /api/v1/requests/{id}", h.handleDeleteRequest)
	mux.HandleFunc("GET /api/v1/history", h.handleListHistory)
}

func (h *RequestsHandler) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}
	item, err := h.requests.Get(r.Context(), r.PathValue("id"), user.ID)
	if err != nil {
		writeSavedRequestError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *RequestsHandler) handleCreateRequest(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}
	var payload savedRequestPayload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		return
	}
	item, err := h.requests.Create(r.Context(), requests.CreateParams{
		CollectionID:    payload.CollectionID,
		OwnerUserID:     user.ID,
		Name:            derefString(payload.Name),
		Description:     derefString(payload.Description),
		Method:          derefString(payload.Method),
		URL:             derefString(payload.URL),
		QueryParams:     payload.QueryParams,
		Headers:         payload.Headers,
		AuthScheme:      derefString(payload.AuthScheme),
		AuthConfig:      payload.AuthConfig,
		BodyMode:        derefString(payload.BodyMode),
		BodyConfig:      payload.BodyConfig,
		ExampleResponse: payload.ExampleResponse,
		Metadata:        payload.Metadata,
	})
	if err != nil {
		writeSavedRequestError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *RequestsHandler) handleUpdateRequest(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}
	var payload savedRequestPayload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		return
	}
	item, err := h.requests.Update(r.Context(), requests.UpdateParams{
		ID:              r.PathValue("id"),
		OwnerUserID:     user.ID,
		CollectionID:    normalizeOptionalNestedString(payload.CollectionID),
		Name:            normalizeOptionalString(payload.Name),
		Description:     normalizeOptionalString(payload.Description),
		Method:          normalizeOptionalString(payload.Method),
		URL:             normalizeOptionalString(payload.URL),
		QueryParams:     normalizeOptionalJSON(payload.QueryParams),
		Headers:         normalizeOptionalJSON(payload.Headers),
		AuthScheme:      normalizeOptionalString(payload.AuthScheme),
		AuthConfig:      normalizeOptionalJSON(payload.AuthConfig),
		BodyMode:        normalizeOptionalString(payload.BodyMode),
		BodyConfig:      normalizeOptionalJSON(payload.BodyConfig),
		ExampleResponse: normalizeOptionalJSON(payload.ExampleResponse),
		Metadata:        normalizeOptionalJSON(payload.Metadata),
	})
	if err != nil {
		writeSavedRequestError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *RequestsHandler) handleDeleteRequest(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}
	if err := h.requests.Delete(r.Context(), r.PathValue("id"), user.ID); err != nil {
		writeSavedRequestError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *RequestsHandler) handleListHistory(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}
	items, err := h.history.List(r.Context(), user.ID, 20)
	if err != nil {
		writeHistoryError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"history": items})
}

func (h *RequestsHandler) requireUser(w http.ResponseWriter, r *http.Request) (auth.UserRecord, bool) {
	if h == nil || h.auth == nil || h.requests == nil || h.history == nil {
		writeError(w, http.StatusServiceUnavailable, "requests_unavailable", "saved requests are temporarily unavailable")
		return auth.UserRecord{}, false
	}
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "missing session")
		return auth.UserRecord{}, false
	}
	user, _, err := h.auth.CurrentUser(r.Context(), cookie.Value)
	if err != nil {
		writeAuthError(w, err)
		return auth.UserRecord{}, false
	}
	return user, true
}

func writeSavedRequestError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, requests.ErrUnavailable):
		writeError(w, http.StatusServiceUnavailable, "requests_unavailable", "saved requests are temporarily unavailable")
	case errors.Is(err, requests.ErrInvalid):
		writeError(w, http.StatusBadRequest, "invalid_saved_request", "saved request payload is invalid")
	case errors.Is(err, requests.ErrNotFound):
		writeError(w, http.StatusNotFound, "saved_request_not_found", "saved request not found")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "unexpected saved request failure")
	}
}

func writeHistoryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, history.ErrUnavailable):
		writeError(w, http.StatusServiceUnavailable, "history_unavailable", "request history is temporarily unavailable")
	case errors.Is(err, history.ErrInvalid):
		writeError(w, http.StatusBadRequest, "invalid_history", "request history input is invalid")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "unexpected history failure")
	}
}

func normalizeOptionalJSON(value json.RawMessage) *json.RawMessage {
	if len(value) == 0 {
		return nil
	}
	copyValue := json.RawMessage(append([]byte(nil), value...))
	return &copyValue
}
