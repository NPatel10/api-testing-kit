package entitlements

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

type stubRepository struct {
	plan         Plan
	entitlements []PlanEntitlement
	getByCodeErr error
	getByIDErr   error
	listPlans    []Plan
	listPlansErr error
	listEntErr   error
}

func (s stubRepository) GetPlanByCode(context.Context, string) (Plan, error) {
	return s.plan, s.getByCodeErr
}

func (s stubRepository) GetPlanByID(context.Context, string) (Plan, error) {
	return s.plan, s.getByIDErr
}

func (s stubRepository) ListActivePlans(context.Context) ([]Plan, error) {
	return s.listPlans, s.listPlansErr
}

func (s stubRepository) ListEntitlementsByPlanID(context.Context, string) ([]PlanEntitlement, error) {
	return s.entitlements, s.listEntErr
}

func TestResolveAccessPolicy(t *testing.T) {
	t.Parallel()

	snapshot := Snapshot{
		Plan: Plan{Code: "pro"},
		Entitlements: map[string]json.RawMessage{
			EntitlementKeyCustomURLExecution:    json.RawMessage("true"),
			EntitlementKeyEnvironmentVariables:  json.RawMessage("false"),
			EntitlementKeySharedLinks:           json.RawMessage("true"),
			EntitlementKeySavedHistoryDepth:     json.RawMessage("25"),
			EntitlementKeyRequestTimeoutSeconds: json.RawMessage("30"),
		},
	}

	policy := ResolveAccessPolicy(snapshot, AccessPolicy{
		CanUseEnvironmentVars: true,
		RequestTimeoutSeconds: 15,
	})

	if !policy.CanUseCustomURLs {
		t.Fatalf("expected custom URLs to be enabled")
	}
	if policy.CanUseEnvironmentVars {
		t.Fatalf("expected environment variables to be disabled by entitlement override")
	}
	if !policy.CanUseSharedLinks {
		t.Fatalf("expected shared links to be enabled")
	}
	if policy.SavedHistoryDepth != 25 {
		t.Fatalf("expected history depth 25, got %d", policy.SavedHistoryDepth)
	}
	if policy.RequestTimeoutSeconds != 30 {
		t.Fatalf("expected timeout 30, got %d", policy.RequestTimeoutSeconds)
	}
}

func TestServiceResolveByCode(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 3, 26, 10, 0, 0, 0, time.UTC)
	repo := stubRepository{
		plan: Plan{
			ID:        "plan-pro",
			Code:      "pro",
			Name:      "Pro",
			IsActive:  true,
			SortOrder: 10,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		},
		entitlements: []PlanEntitlement{
			{Key: EntitlementKeyCustomURLExecution, Value: json.RawMessage("true")},
			{Key: EntitlementKeySavedHistoryDepth, Value: json.RawMessage("100")},
		},
	}

	service := NewService(repo)
	resolved, err := service.ResolveByCode(context.Background(), "pro")
	if err != nil {
		t.Fatalf("resolve by code failed: %v", err)
	}

	if resolved.Snapshot.Plan.Code != "pro" {
		t.Fatalf("expected plan code pro, got %q", resolved.Snapshot.Plan.Code)
	}
	if !resolved.Policy.CanUseCustomURLs {
		t.Fatalf("expected custom URL access to be enabled")
	}
	if resolved.Policy.SavedHistoryDepth != 100 {
		t.Fatalf("expected history depth 100, got %d", resolved.Policy.SavedHistoryDepth)
	}
}

func TestServiceRejectsMissingRepo(t *testing.T) {
	t.Parallel()

	if _, err := NewService(nil).ResolveByCode(context.Background(), "pro"); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("expected unavailable error, got %v", err)
	}
}
