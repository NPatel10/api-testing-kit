package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"api-testing-kit/server/internal/auth"
	"api-testing-kit/server/internal/collections"
)

type CollectionsHandler struct {
	service *collections.Service
	auth    *auth.Service
}

type collectionRequest struct {
	Name        *string                 `json:"name,omitempty"`
	Slug        *string                 `json:"slug,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Visibility  *collections.Visibility `json:"visibility,omitempty"`
	Color       *string                 `json:"color,omitempty"`
	SortOrder   *int                    `json:"sortOrder,omitempty"`
	Metadata    json.RawMessage         `json:"metadata,omitempty"`
}

func NewCollectionsHandler(service *collections.Service, authService *auth.Service) *CollectionsHandler {
	return &CollectionsHandler{
		service: service,
		auth:    authService,
	}
}

func (h *CollectionsHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/collections", h.handleList)
	mux.HandleFunc("POST /api/v1/collections", h.handleCreate)
	mux.HandleFunc("PATCH /api/v1/collections/{id}", h.handleUpdate)
	mux.HandleFunc("DELETE /api/v1/collections/{id}", h.handleDelete)
}

func (h *CollectionsHandler) handleList(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	items, err := h.service.List(r.Context(), user.ID)
	if err != nil {
		writeCollectionError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"collections": items})
}

func (h *CollectionsHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	var payload collectionRequest
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		return
	}

	name := ""
	if payload.Name != nil {
		name = *payload.Name
	}

	item, err := h.service.Create(r.Context(), collections.CreateParams{
		OwnerUserID: user.ID,
		Name:        name,
		Slug:        payload.Slug,
		Description: derefString(payload.Description),
		Visibility:  derefVisibility(payload.Visibility),
		Color:       derefString(payload.Color),
		SortOrder:   derefInt(payload.SortOrder),
		Metadata:    payload.Metadata,
	})
	if err != nil {
		writeCollectionError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (h *CollectionsHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	var payload collectionRequest
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		return
	}

	var metadata *json.RawMessage
	if len(payload.Metadata) > 0 {
		copyValue := json.RawMessage(append([]byte(nil), payload.Metadata...))
		metadata = &copyValue
	}

	item, err := h.service.Update(r.Context(), collections.UpdateParams{
		ID:          r.PathValue("id"),
		OwnerUserID: user.ID,
		Name:        normalizeOptionalString(payload.Name),
		Slug:        normalizeOptionalNestedString(payload.Slug),
		Description: normalizeOptionalString(payload.Description),
		Visibility:  payload.Visibility,
		Color:       normalizeOptionalString(payload.Color),
		SortOrder:   payload.SortOrder,
		Metadata:    metadata,
	})
	if err != nil {
		writeCollectionError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *CollectionsHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), r.PathValue("id"), user.ID); err != nil {
		writeCollectionError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CollectionsHandler) requireUser(w http.ResponseWriter, r *http.Request) (auth.UserRecord, bool) {
	if h == nil || h.service == nil || h.auth == nil {
		writeError(w, http.StatusServiceUnavailable, "collections_unavailable", "collections are temporarily unavailable")
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

func writeCollectionError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, collections.ErrUnavailable):
		writeError(w, http.StatusServiceUnavailable, "collections_unavailable", "collections are temporarily unavailable")
	case errors.Is(err, collections.ErrInvalid):
		writeError(w, http.StatusBadRequest, "invalid_collection", "collection payload is invalid")
	case errors.Is(err, collections.ErrNotFound):
		writeError(w, http.StatusNotFound, "collection_not_found", "collection not found")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "unexpected collections failure")
	}
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func derefVisibility(value *collections.Visibility) collections.Visibility {
	if value == nil {
		return ""
	}
	return *value
}

func derefInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	copyValue := *value
	return &copyValue
}

func normalizeOptionalNestedString(value *string) **string {
	if value == nil {
		return nil
	}

	copyValue := *value
	inner := &copyValue
	return &inner
}
