package db

import (
	"strings"
	"testing"
	"time"

	"api-testing-kit/server/internal/history"
)

func TestBuildRequestHistoryListQueryAppliesFilters(t *testing.T) {
	t.Parallel()

	date := time.Date(2026, time.March, 26, 0, 0, 0, 0, time.UTC)
	query, args := buildRequestHistoryListQuery(history.ListQuery{
		UserID: "user-1",
		Status: "success",
		Method: "GET",
		Domain: "api.example.com",
		Date:   &date,
		Page:   2,
		Limit:  21,
	})

	if !strings.Contains(query, "status::text = ANY($2::text[])") {
		t.Fatalf("expected status filter in query, got %s", query)
	}
	if !strings.Contains(query, "method::text = $3") {
		t.Fatalf("expected method filter in query, got %s", query)
	}
	if !strings.Contains(query, "COALESCE(target_host, '') = $4") {
		t.Fatalf("expected domain filter in query, got %s", query)
	}
	if !strings.Contains(query, "created_at >= $5 AND created_at < ($5 + INTERVAL '1 day')") {
		t.Fatalf("expected date filter in query, got %s", query)
	}
	if !strings.Contains(query, "LIMIT $6 OFFSET $7") {
		t.Fatalf("expected pagination in query, got %s", query)
	}

	statuses, ok := args[1].([]string)
	if !ok || len(statuses) != 1 || statuses[0] != "succeeded" {
		t.Fatalf("expected normalized success status, got %#v", args[1])
	}
	if got := args[2].(string); got != "GET" {
		t.Fatalf("expected normalized method, got %q", got)
	}
	if got := args[3].(string); got != "api.example.com" {
		t.Fatalf("expected normalized domain, got %q", got)
	}
	if got := args[4].(time.Time); !got.Equal(date.UTC()) {
		t.Fatalf("expected normalized date, got %v", got)
	}
	if got := args[5].(int32); got != 21 {
		t.Fatalf("expected fetch limit with sentinel, got %d", got)
	}
	if got := args[6].(int32); got != 20 {
		t.Fatalf("expected offset for page 2, got %d", got)
	}
}
