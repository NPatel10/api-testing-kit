package history

import (
	"context"
	"testing"
	"time"
)

type historyListRepo struct {
	items     []RunRecord
	lastQuery ListQuery
}

func (r *historyListRepo) ListByUser(ctx context.Context, userID string, limit int32) ([]RunRecord, error) {
	return r.items, nil
}

func (r *historyListRepo) ListByUserWithFilters(ctx context.Context, query ListQuery) ([]RunRecord, error) {
	r.lastQuery = query
	return r.items, nil
}

func (r *historyListRepo) Create(ctx context.Context, params CreateParams) (RunRecord, error) {
	return RunRecord{}, nil
}

func TestListWithFiltersNormalizesPagination(t *testing.T) {
	t.Parallel()

	repo := &historyListRepo{
		items: []RunRecord{{ID: "run-1"}, {ID: "run-2"}},
	}
	service := NewService(repo)

	date := time.Date(2026, time.March, 26, 0, 0, 0, 0, time.UTC)
	result, err := service.ListWithFilters(context.Background(), ListQuery{
		UserID: "  user-1  ",
		Status: "success",
		Method: "all",
		Domain: "Api.Example.com",
		Date:   &date,
		Page:   2,
		Limit:  1,
	})
	if err != nil {
		t.Fatalf("list with filters failed: %v", err)
	}

	if repo.lastQuery.UserID != "user-1" {
		t.Fatalf("expected trimmed user id, got %q", repo.lastQuery.UserID)
	}
	if repo.lastQuery.Status != "success" {
		t.Fatalf("expected status to be forwarded, got %q", repo.lastQuery.Status)
	}
	if repo.lastQuery.Method != "" {
		t.Fatalf("expected all-method filter to normalize away, got %q", repo.lastQuery.Method)
	}
	if repo.lastQuery.Domain != "api.example.com" {
		t.Fatalf("expected normalized domain, got %q", repo.lastQuery.Domain)
	}
	if repo.lastQuery.Page != 2 || repo.lastQuery.Limit != 2 {
		t.Fatalf("expected sentinel fetch limit and page to be forwarded, got %+v", repo.lastQuery)
	}

	if len(result.Items) != 1 {
		t.Fatalf("expected page-sized result, got %d items", len(result.Items))
	}
	if !result.Pagination.HasMore || result.Pagination.Page != 2 || result.Pagination.Limit != 1 {
		t.Fatalf("unexpected pagination metadata: %+v", result.Pagination)
	}
}
