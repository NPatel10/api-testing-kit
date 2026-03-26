package db

import (
	"context"

	"api-testing-kit/server/internal/entitlements"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PlanEntitlementRepository struct {
	pool *pgxpool.Pool
}

func NewPlanEntitlementRepository(pool *pgxpool.Pool) *PlanEntitlementRepository {
	return &PlanEntitlementRepository{pool: pool}
}

func (r *PlanEntitlementRepository) ListByPlanID(ctx context.Context, planID string) ([]entitlements.PlanEntitlement, error) {
	rows, err := r.pool.Query(ctx, buildPlanEntitlementsByPlanIDQuery(), planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entitlements.PlanEntitlement, 0)
	for rows.Next() {
		var item entitlements.PlanEntitlement
		if err := rows.Scan(&item.ID, &item.PlanID, &item.Key, &item.Value, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func buildPlanEntitlementsByPlanIDQuery() string {
	return `
		SELECT
			id,
			plan_id,
			key,
			value,
			created_at
		FROM plan_entitlements
		WHERE plan_id = $1
		ORDER BY key ASC
	`
}
