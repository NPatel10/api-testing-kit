package db

import (
	"strings"
	"testing"
	"time"

	"api-testing-kit/server/internal/usage"
)

func TestBuildUsageRecentQueryDefaultsLimit(t *testing.T) {
	t.Parallel()

	query, args := buildUsageRecentQuery(usage.RecentFilter{})

	if got, want := len(args), 1; got != want {
		t.Fatalf("expected %d arg, got %d", want, got)
	}

	if got, want := args[0], int32(100); got != want {
		t.Fatalf("expected default limit %v, got %v", want, got)
	}

	if !strings.Contains(query, "FROM usage_events") {
		t.Fatalf("expected usage_events table in query")
	}

	if !strings.Contains(query, "ORDER BY occurred_at DESC LIMIT $1") {
		t.Fatalf("expected limit placeholder at end of query, got %q", query)
	}
}

func TestBuildUsageSummaryQueryWithFilters(t *testing.T) {
	t.Parallel()

	userID := "user-123"
	since := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)

	query, args := buildUsageSummaryQuery(usage.SummaryFilter{
		UserID:         &userID,
		Bucket:         "day",
		EventKey:       "api.request",
		OccurredAfter:  &since,
		OccurredBefore: &until,
		Limit:          25,
	})

	if got, want := len(args), 6; got != want {
		t.Fatalf("expected %d args, got %d", want, got)
	}

	if got, want := args[len(args)-1], int32(25); got != want {
		t.Fatalf("expected limit %v, got %v", want, got)
	}

	for _, want := range []string{"user_id = $1", "bucket = $2", "event_key = $3", "occurred_at >=", "occurred_at <=", "GROUP BY bucket, event_key"} {
		if !strings.Contains(query, want) {
			t.Fatalf("expected query to contain %q, got %q", want, query)
		}
	}
}
