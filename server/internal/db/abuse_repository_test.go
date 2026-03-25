package db

import (
	"strings"
	"testing"
	"time"

	"api-testing-kit/server/internal/abuse"
)

func TestBuildAbuseRecentQueryIncludesSeverityAndAction(t *testing.T) {
	t.Parallel()

	since := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	query, args := buildAbuseRecentQuery(abuse.RecentFilter{
		Severity:     abuse.SeverityHigh,
		Category:     "ssrf",
		ActionTaken:  "blocked",
		CreatedAfter: &since,
		Limit:        10,
	})

	if got, want := len(args), 5; got != want {
		t.Fatalf("expected %d args, got %d", want, got)
	}

	if got, want := args[len(args)-1], int32(10); got != want {
		t.Fatalf("expected limit %v, got %v", want, got)
	}

	for _, want := range []string{"severity = $1", "category = $2", "action_taken = $3", "created_at >=", "ORDER BY created_at DESC LIMIT $5"} {
		if !strings.Contains(query, want) {
			t.Fatalf("expected query to contain %q, got %q", want, query)
		}
	}
}

func TestBuildBlockedTargetQueryKeepsExpiryGrouped(t *testing.T) {
	t.Parallel()

	query, args := buildBlockedTargetQuery(abuse.BlockedTargetFilter{
		TargetType:     abuse.BlockedTargetTypeIP,
		OnlyActive:     true,
		IncludeExpired: false,
		Limit:          20,
	})

	if got, want := len(args), 3; got != want {
		t.Fatalf("expected %d args, got %d", want, got)
	}

	if !strings.Contains(query, "(expires_at IS NULL OR expires_at > now())") {
		t.Fatalf("expected expiry clause to remain grouped, got %q", query)
	}

	if !strings.Contains(query, "ORDER BY is_active DESC, updated_at DESC") {
		t.Fatalf("expected deterministic blocked-target ordering, got %q", query)
	}
}
