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
	handler := NewCollectionsHandler(collections.NewService(collectionsRepo), authService)
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
