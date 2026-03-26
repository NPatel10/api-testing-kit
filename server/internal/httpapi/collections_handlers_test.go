package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"api-testing-kit/server/internal/auth"
	"api-testing-kit/server/internal/collections"
	"api-testing-kit/server/internal/requests"
)

func TestCollectionsCRUD(t *testing.T) {
	t.Parallel()

	authRepo := newFakeAuthRepo()
	authService := auth.NewService(authRepo)
	authResult, err := authService.Signup(context.Background(), auth.SignupInput{
		Email:    "collections@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("signup failed: %v", err)
	}

	collectionsRepo := newFakeCollectionsRepo()
	requestsRepo := newFakeCollectionsSavedRequestRepo()
	handler := NewCollectionsHandler(collections.NewService(collectionsRepo), requests.NewService(requestsRepo), authService)
	mux := http.NewServeMux()
	handler.Register(mux)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/collections", strings.NewReader(`{"name":"Primary","description":"Saved requests"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: authResult.Token})
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)

	if createRR.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, createRR.Code)
	}

	var created collections.Collection
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode created collection: %v", err)
	}

	ownerID := authResult.User.ID
	requestsRepo.items["request-2"] = requests.SavedRequest{
		ID:           "request-2",
		CollectionID: &created.ID,
		OwnerUserID:  &ownerID,
		Name:         "Second",
		Method:       http.MethodPost,
		URL:          "https://api.example.com/second",
		SortOrder:    2,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	requestsRepo.items["request-1"] = requests.SavedRequest{
		ID:           "request-1",
		CollectionID: &created.ID,
		OwnerUserID:  &ownerID,
		Name:         "First",
		Method:       http.MethodGet,
		URL:          "https://api.example.com/first",
		SortOrder:    1,
		CreatedAt:    time.Now().UTC().Add(-time.Minute),
		UpdatedAt:    time.Now().UTC().Add(-time.Minute),
	}

	detailReq := httptest.NewRequest(http.MethodGet, "/api/v1/collections/"+created.ID, nil)
	detailReq.SetPathValue("id", created.ID)
	detailReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: authResult.Token})
	detailRR := httptest.NewRecorder()
	mux.ServeHTTP(detailRR, detailReq)

	if detailRR.Code != http.StatusOK {
		t.Fatalf("expected detail status %d, got %d", http.StatusOK, detailRR.Code)
	}

	var detailPayload struct {
		Collection    collections.Collection  `json:"collection"`
		SavedRequests []requests.SavedRequest `json:"savedRequests"`
	}
	if err := json.Unmarshal(detailRR.Body.Bytes(), &detailPayload); err != nil {
		t.Fatalf("failed to decode collection detail: %v", err)
	}
	if detailPayload.Collection.ID != created.ID {
		t.Fatalf("expected detail for collection %q, got %q", created.ID, detailPayload.Collection.ID)
	}
	if len(detailPayload.SavedRequests) != 2 {
		t.Fatalf("expected 2 saved requests, got %d", len(detailPayload.SavedRequests))
	}
	if detailPayload.SavedRequests[0].SortOrder != 1 || detailPayload.SavedRequests[1].SortOrder != 2 {
		t.Fatalf("expected saved requests to be sorted by sort order, got %+v", detailPayload.SavedRequests)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/collections", nil)
	listReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: authResult.Token})
	listRR := httptest.NewRecorder()
	mux.ServeHTTP(listRR, listReq)

	if listRR.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d", http.StatusOK, listRR.Code)
	}

	updateReq := httptest.NewRequest(http.MethodPatch, "/api/v1/collections/"+created.ID, strings.NewReader(`{"name":"Renamed"}`))
	updateReq.SetPathValue("id", created.ID)
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: authResult.Token})
	updateRR := httptest.NewRecorder()
	mux.ServeHTTP(updateRR, updateReq)

	if updateRR.Code != http.StatusOK {
		t.Fatalf("expected update status %d, got %d", http.StatusOK, updateRR.Code)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/collections/"+created.ID, nil)
	deleteReq.SetPathValue("id", created.ID)
	deleteReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: authResult.Token})
	deleteRR := httptest.NewRecorder()
	mux.ServeHTTP(deleteRR, deleteReq)

	if deleteRR.Code != http.StatusNoContent {
		t.Fatalf("expected delete status %d, got %d", http.StatusNoContent, deleteRR.Code)
	}
}

type fakeCollectionsRepo struct {
	items map[string]collections.Collection
	next  int
}

func newFakeCollectionsRepo() *fakeCollectionsRepo {
	return &fakeCollectionsRepo{
		items: make(map[string]collections.Collection),
	}
}

func (r *fakeCollectionsRepo) ListByOwner(ctx context.Context, ownerUserID string) ([]collections.Collection, error) {
	items := make([]collections.Collection, 0)
	for _, item := range r.items {
		if item.OwnerUserID != nil && *item.OwnerUserID == ownerUserID && item.DeletedAt == nil {
			items = append(items, item)
		}
	}
	return items, nil
}

func (r *fakeCollectionsRepo) GetByID(ctx context.Context, id string, ownerUserID string) (collections.Collection, error) {
	item, ok := r.items[id]
	if !ok || item.OwnerUserID == nil || *item.OwnerUserID != ownerUserID || item.DeletedAt != nil {
		return collections.Collection{}, collections.ErrNotFound
	}
	return item, nil
}

func (r *fakeCollectionsRepo) Create(ctx context.Context, params collections.CreateParams) (collections.Collection, error) {
	r.next++
	now := time.Now().UTC()
	ownerID := params.OwnerUserID
	item := collections.Collection{
		ID:          fmt.Sprintf("collection-%d", r.next),
		OwnerUserID: &ownerID,
		Name:        params.Name,
		Slug:        params.Slug,
		Description: params.Description,
		Visibility:  params.Visibility,
		Color:       params.Color,
		SortOrder:   params.SortOrder,
		Metadata:    params.Metadata,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.items[item.ID] = item
	return item, nil
}

func (r *fakeCollectionsRepo) Update(ctx context.Context, params collections.UpdateParams) (collections.Collection, error) {
	item, ok := r.items[params.ID]
	if !ok || item.OwnerUserID == nil || *item.OwnerUserID != params.OwnerUserID {
		return collections.Collection{}, collections.ErrNotFound
	}
	if params.Name != nil {
		item.Name = *params.Name
	}
	item.UpdatedAt = time.Now().UTC()
	r.items[item.ID] = item
	return item, nil
}

func (r *fakeCollectionsRepo) Delete(ctx context.Context, id string, ownerUserID string) error {
	item, ok := r.items[id]
	if !ok || item.OwnerUserID == nil || *item.OwnerUserID != ownerUserID {
		return collections.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.items[id] = item
	return nil
}

type fakeCollectionsSavedRequestRepo struct {
	items map[string]requests.SavedRequest
}

func newFakeCollectionsSavedRequestRepo() *fakeCollectionsSavedRequestRepo {
	return &fakeCollectionsSavedRequestRepo{items: make(map[string]requests.SavedRequest)}
}

func (r *fakeCollectionsSavedRequestRepo) GetByID(ctx context.Context, id string, ownerUserID string) (requests.SavedRequest, error) {
	item, ok := r.items[id]
	if !ok || item.OwnerUserID == nil || *item.OwnerUserID != ownerUserID {
		return requests.SavedRequest{}, requests.ErrNotFound
	}
	return item, nil
}

func (r *fakeCollectionsSavedRequestRepo) ListByCollection(ctx context.Context, collectionID string, ownerUserID string) ([]requests.SavedRequest, error) {
	items := make([]requests.SavedRequest, 0)
	for _, item := range r.items {
		if item.CollectionID != nil && *item.CollectionID == collectionID && item.OwnerUserID != nil && *item.OwnerUserID == ownerUserID {
			items = append(items, item)
		}
	}
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].SortOrder < items[i].SortOrder {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	return items, nil
}

func (r *fakeCollectionsSavedRequestRepo) Create(ctx context.Context, params requests.CreateParams) (requests.SavedRequest, error) {
	return requests.SavedRequest{}, requests.ErrUnavailable
}

func (r *fakeCollectionsSavedRequestRepo) Update(ctx context.Context, params requests.UpdateParams) (requests.SavedRequest, error) {
	return requests.SavedRequest{}, requests.ErrUnavailable
}

func (r *fakeCollectionsSavedRequestRepo) Delete(ctx context.Context, id string, ownerUserID string) error {
	return requests.ErrUnavailable
}
