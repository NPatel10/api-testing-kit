package db

import (
	"context"
	"database/sql"
	"errors"

	"api-testing-kit/server/internal/entitlements"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlanRepository struct {
	pool *pgxpool.Pool
}

func NewPlanRepository(pool *pgxpool.Pool) *PlanRepository {
	return &PlanRepository{pool: pool}
}

func (r *PlanRepository) ListActivePlans(ctx context.Context) ([]entitlements.Plan, error) {
	rows, err := r.pool.Query(ctx, buildPlanListQuery())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entitlements.Plan, 0)
	for rows.Next() {
		item, err := scanPlan(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *PlanRepository) GetPlanByCode(ctx context.Context, code string) (entitlements.Plan, error) {
	row := r.pool.QueryRow(ctx, buildPlanByCodeQuery(), code)
	item, err := scanPlan(row.Scan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entitlements.Plan{}, entitlements.ErrNotFound
		}
		return entitlements.Plan{}, err
	}

	return item, nil
}

func (r *PlanRepository) GetPlanByID(ctx context.Context, id string) (entitlements.Plan, error) {
	row := r.pool.QueryRow(ctx, buildPlanByIDQuery(), id)
	item, err := scanPlan(row.Scan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entitlements.Plan{}, entitlements.ErrNotFound
		}
		return entitlements.Plan{}, err
	}

	return item, nil
}

func buildPlanListQuery() string {
	return `
		SELECT
			id,
			code,
			name,
			COALESCE(description, ''),
			is_active,
			sort_order,
			metadata,
			created_at,
			updated_at
		FROM plans
		WHERE is_active = TRUE
		ORDER BY sort_order ASC, code ASC
	`
}

func buildPlanByCodeQuery() string {
	return `
		SELECT
			id,
			code,
			name,
			COALESCE(description, ''),
			is_active,
			sort_order,
			metadata,
			created_at,
			updated_at
		FROM plans
		WHERE code = $1
		  AND is_active = TRUE
	`
}

func buildPlanByIDQuery() string {
	return `
		SELECT
			id,
			code,
			name,
			COALESCE(description, ''),
			is_active,
			sort_order,
			metadata,
			created_at,
			updated_at
		FROM plans
		WHERE id = $1
		  AND is_active = TRUE
	`
}

func scanPlan(scan func(dest ...any) error) (entitlements.Plan, error) {
	var item entitlements.Plan
	var description sql.NullString

	if err := scan(
		&item.ID,
		&item.Code,
		&item.Name,
		&description,
		&item.IsActive,
		&item.SortOrder,
		&item.Metadata,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return entitlements.Plan{}, err
	}

	if description.Valid {
		item.Description = description.String
	}

	return item, nil
}
