package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"api-testing-kit/server/internal/abuse"
	"api-testing-kit/server/internal/auth"
)

func TestAdminAbuseRoutesForAdmins(t *testing.T) {
	t.Parallel()

	authRepo := newFakeAuthRepo()
	authService := auth.NewService(authRepo)
	adminResult := mustSignUpAndPromote(t, authService, authRepo, "admin@example.com", "admin")
	userResult := mustSignUpAndPromote(t, authService, authRepo, "user@example.com", "user")

	handler := NewAdminHandler(
		&fakeAdminAbuseRepo{
			summary: []abuse.SummaryRow{
				{
					Severity:      abuse.SeverityHigh,
					Category:      abuse.CategoryBlockedHost,
					ActionTaken:   abuse.ActionBlocked,
					Count:         4,
					LastCreatedAt: time.Date(2026, time.March, 26, 8, 0, 0, 0, time.UTC),
				},
			},
			recent: []abuse.Event{
				{
					ID:          "abuse-1",
					RuleKey:     "blocked-target",
					Category:    abuse.CategoryBlockedHost,
					Severity:    abuse.SeverityHigh,
					ActionTaken: abuse.ActionBlocked,
					Message:     "blocked target surfaced",
					CreatedAt:   time.Date(2026, time.March, 26, 8, 1, 0, 0, time.UTC),
				},
			},
		},
		&fakeBlockedTargetRepo{
			items: []abuse.BlockedTarget{
				{
					ID:          "blocked-1",
					TargetType:  abuse.BlockedTargetTypeIP,
					TargetValue: "169.254.169.254",
					Reason:      "metadata IP",
					Source:      "manual",
					IsActive:    true,
					CreatedAt:   time.Date(2026, time.March, 26, 7, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2026, time.March, 26, 7, 0, 0, 0, time.UTC),
				},
			},
		},
		authService,
	)

	mux := http.NewServeMux()
	handler.Register(mux)

	adminReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/abuse", nil)
	adminReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: adminResult.Token})
	adminRR := httptest.NewRecorder()
	mux.ServeHTTP(adminRR, adminReq)

	if adminRR.Code != http.StatusOK {
		t.Fatalf("expected admin abuse status %d, got %d", http.StatusOK, adminRR.Code)
	}

	var abusePayload struct {
		Summary        []abuse.SummaryRow    `json:"summary"`
		Recent         []abuse.Event         `json:"recent"`
		BlockedTargets []abuse.BlockedTarget `json:"blockedTargets"`
	}
	if err := json.Unmarshal(adminRR.Body.Bytes(), &abusePayload); err != nil {
		t.Fatalf("failed to decode abuse payload: %v", err)
	}
	if len(abusePayload.Summary) != 1 || len(abusePayload.Recent) != 1 || len(abusePayload.BlockedTargets) != 1 {
		t.Fatalf("unexpected abuse payload shape: %+v", abusePayload)
	}

	blockedReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/blocked-targets", nil)
	blockedReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: adminResult.Token})
	blockedRR := httptest.NewRecorder()
	mux.ServeHTTP(blockedRR, blockedReq)

	if blockedRR.Code != http.StatusOK {
		t.Fatalf("expected blocked targets status %d, got %d", http.StatusOK, blockedRR.Code)
	}

	var blockedPayload struct {
		BlockedTargets []abuse.BlockedTarget `json:"blockedTargets"`
	}
	if err := json.Unmarshal(blockedRR.Body.Bytes(), &blockedPayload); err != nil {
		t.Fatalf("failed to decode blocked targets payload: %v", err)
	}
	if len(blockedPayload.BlockedTargets) != 1 {
		t.Fatalf("expected one blocked target, got %d", len(blockedPayload.BlockedTargets))
	}

	userReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/abuse", nil)
	userReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: userResult.Token})
	userRR := httptest.NewRecorder()
	mux.ServeHTTP(userRR, userReq)

	if userRR.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden status for non-admin user, got %d", userRR.Code)
	}

	if !strings.Contains(userRR.Body.String(), "forbidden") {
		t.Fatalf("expected forbidden error payload, got %s", userRR.Body.String())
	}
}

func TestAdminAbuseRoutesWithoutServiceAreUnavailable(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	NewAdminHandler(nil, nil, nil).Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/abuse", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected unavailable status %d, got %d", http.StatusServiceUnavailable, rr.Code)
	}
}

func mustSignUpAndPromote(t *testing.T, service *auth.Service, repo *fakeAuthRepo, email, role string) auth.AuthResult {
	t.Helper()

	result, err := service.Signup(context.Background(), auth.SignupInput{
		Email:    email,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("signup failed for %s: %v", email, err)
	}

	record := repo.users[result.User.Email]
	record.Role = role
	repo.users[result.User.Email] = record

	return result
}

type fakeAdminAbuseRepo struct {
	summary []abuse.SummaryRow
	recent  []abuse.Event
}

func (r *fakeAdminAbuseRepo) SummarizeByCategory(ctx context.Context, filter abuse.SummaryFilter) ([]abuse.SummaryRow, error) {
	return append([]abuse.SummaryRow(nil), r.summary...), nil
}

func (r *fakeAdminAbuseRepo) ListRecent(ctx context.Context, filter abuse.RecentFilter) ([]abuse.Event, error) {
	return append([]abuse.Event(nil), r.recent...), nil
}

type fakeBlockedTargetRepo struct {
	items []abuse.BlockedTarget
}

func (r *fakeBlockedTargetRepo) List(ctx context.Context, filter abuse.BlockedTargetFilter) ([]abuse.BlockedTarget, error) {
	return append([]abuse.BlockedTarget(nil), r.items...), nil
}
