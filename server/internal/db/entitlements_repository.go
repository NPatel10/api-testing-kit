package db

import (
	"context"

	"api-testing-kit/server/internal/entitlements"
)

type EntitlementRepository struct {
	plans            *PlanRepository
	planEntitlements *PlanEntitlementRepository
}

func NewEntitlementRepository(plans *PlanRepository, planEntitlements *PlanEntitlementRepository) *EntitlementRepository {
	return &EntitlementRepository{
		plans:            plans,
		planEntitlements: planEntitlements,
	}
}

func (r *EntitlementRepository) GetPlanByCode(ctx context.Context, code string) (entitlements.Plan, error) {
	return r.plans.GetPlanByCode(ctx, code)
}

func (r *EntitlementRepository) GetPlanByID(ctx context.Context, id string) (entitlements.Plan, error) {
	return r.plans.GetPlanByID(ctx, id)
}

func (r *EntitlementRepository) ListActivePlans(ctx context.Context) ([]entitlements.Plan, error) {
	return r.plans.ListActivePlans(ctx)
}

func (r *EntitlementRepository) ListEntitlementsByPlanID(ctx context.Context, planID string) ([]entitlements.PlanEntitlement, error) {
	return r.planEntitlements.ListByPlanID(ctx, planID)
}
