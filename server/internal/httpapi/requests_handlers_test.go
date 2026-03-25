package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"api-testing-kit/server/internal/auth"
	"api-testing-kit/server/internal/history"
	"api-testing-kit/server/internal/requests"
)

func TestSavedRequestsAndHistoryHandlers(t *testing.T) {
	t.Parallel()

	authRepo := newFakeAuthRepo()
	authService := auth.NewService(authRepo)
	authResult, err := authService.Signup(context.Background(), auth.SignupInput{
		Email:    "requests@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("signup failed: %v", err)
	}

	requestRepo := newFakeSavedRequestRepo()
	historyRepo := newFakeHistoryRepo()
	historyService := history.NewService(historyRepo)
	if _, err := historyService.Create(context.Background(), history.CreateParams{
		UserID: authResult.User.ID,
		Method: "GET",
		URL:    "https://api.example.com/users",
	}); err != nil {
		t.Fatalf("create history failed: %v", err)
	}

	handler := NewRequestsHandler(requests.NewService(requestRepo), historyService, authService)
	mux := http.NewServeMux()
	handler.Register(mux)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/requests", strings.NewReader(`{"name":"Fetch users","method":"GET","url":"https://api.example.com/users"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: authResult.Token})
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)

	if createRR.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, createRR.Code)
	}

	var created requests.SavedRequest
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created request: %v", err)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/requests/"+created.ID, nil)
	getReq.SetPathValue("id", created.ID)
	getReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: authResult.Token})
	getRR := httptest.NewRecorder()
	mux.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusOK {
		t.Fatalf("expected get status %d, got %d", http.StatusOK, getRR.Code)
	}

	historyReq := httptest.NewRequest(http.MethodGet, "/api/v1/history", nil)
	historyReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: authResult.Token})
	historyRR := httptest.NewRecorder()
	mux.ServeHTTP(historyRR, historyReq)

	if historyRR.Code != http.StatusOK {
		t.Fatalf("expected history status %d, got %d", http.StatusOK, historyRR.Code)
	}
}

type fakeSavedRequestRepo struct {
	items map[string]requests.SavedRequest
	next  int
}

func newFakeSavedRequestRepo() *fakeSavedRequestRepo {
	return &fakeSavedRequestRepo{items: make(map[string]requests.SavedRequest)}
}

func (r *fakeSavedRequestRepo) GetByID(ctx context.Context, id string, ownerUserID string) (requests.SavedRequest, error) {
	item, ok := r.items[id]
	if !ok || item.OwnerUserID == nil || *item.OwnerUserID != ownerUserID {
		return requests.SavedRequest{}, requests.ErrNotFound
	}
	return item, nil
}

func (r *fakeSavedRequestRepo) Create(ctx context.Context, params requests.CreateParams) (requests.SavedRequest, error) {
	r.next++
	now := time.Now().UTC()
	ownerID := params.OwnerUserID
	item := requests.SavedRequest{
		ID:              "request-" + time.Now().Format("150405.000"),
		CollectionID:    params.CollectionID,
		OwnerUserID:     &ownerID,
		Name:            params.Name,
		Description:     params.Description,
		Method:          params.Method,
		URL:             params.URL,
		QueryParams:     params.QueryParams,
		Headers:         params.Headers,
		AuthScheme:      params.AuthScheme,
		AuthConfig:      params.AuthConfig,
		BodyMode:        params.BodyMode,
		BodyConfig:      params.BodyConfig,
		ExampleResponse: params.ExampleResponse,
		Metadata:        params.Metadata,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	r.items[item.ID] = item
	return item, nil
}

func (r *fakeSavedRequestRepo) Update(ctx context.Context, params requests.UpdateParams) (requests.SavedRequest, error) {
	item, ok := r.items[params.ID]
	if !ok || item.OwnerUserID == nil || *item.OwnerUserID != params.OwnerUserID {
		return requests.SavedRequest{}, requests.ErrNotFound
	}
	if params.Name != nil {
		item.Name = *params.Name
	}
	item.UpdatedAt = time.Now().UTC()
	r.items[item.ID] = item
	return item, nil
}

func (r *fakeSavedRequestRepo) Delete(ctx context.Context, id string, ownerUserID string) error {
	item, ok := r.items[id]
	if !ok || item.OwnerUserID == nil || *item.OwnerUserID != ownerUserID {
		return requests.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.items[id] = item
	return nil
}

type fakeHistoryRepo struct {
	items []history.RunRecord
}

func newFakeHistoryRepo() *fakeHistoryRepo {
	return &fakeHistoryRepo{}
}

func (r *fakeHistoryRepo) ListByUser(ctx context.Context, userID string, limit int32) ([]history.RunRecord, error) {
	return r.items, nil
}

func (r *fakeHistoryRepo) Create(ctx context.Context, params history.CreateParams) (history.RunRecord, error) {
	now := time.Now().UTC()
	userID := params.UserID
	item := history.RunRecord{
		ID:        "run-1",
		UserID:    &userID,
		Source:    params.Source,
		Status:    params.Status,
		Method:    params.Method,
		URL:       params.URL,
		CreatedAt: now,
		Metadata:  params.Metadata,
	}
	r.items = append(r.items, item)
	return item, nil
}
