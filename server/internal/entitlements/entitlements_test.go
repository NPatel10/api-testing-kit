package entitlements

import (
	"testing"

	"api-testing-kit/server/internal/auth"
)

func TestResolveGuestState(t *testing.T) {
	t.Parallel()

	state := Resolve(auth.UserRecord{})

	if state.Plan.Code != "guest" {
		t.Fatalf("expected guest plan, got %q", state.Plan.Code)
	}

	if state.CanUseCustomURLs() {
		t.Fatalf("guest sessions should not allow custom URLs")
	}

	if state.HistoryDepthLimit() != 0 {
		t.Fatalf("expected zero history depth for guests, got %d", state.HistoryDepthLimit())
	}
}

func TestResolveAuthenticatedState(t *testing.T) {
	t.Parallel()

	state := Resolve(auth.UserRecord{ID: "user-1", Status: "active", Role: "user"})

	if state.Plan.Code != "starter" {
		t.Fatalf("expected starter plan, got %q", state.Plan.Code)
	}

	if !state.CanUseCustomURLs() {
		t.Fatalf("authenticated sessions should allow custom URLs")
	}

	if state.HistoryDepthLimit() != 50 {
		t.Fatalf("expected history depth of 50, got %d", state.HistoryDepthLimit())
	}

	if state.CanUseEnvironmentVariables() {
		t.Fatalf("starter plan should keep environment variables locked")
	}
}
