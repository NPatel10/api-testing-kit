package db

import (
	"context"
	"fmt"
	"strings"

	"api-testing-kit/server/internal/abuse"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlockedTargetRepository struct {
	pool *pgxpool.Pool
}

func NewBlockedTargetRepository(pool *pgxpool.Pool) *BlockedTargetRepository {
	return &BlockedTargetRepository{pool: pool}
}

func (r *BlockedTargetRepository) List(ctx context.Context, filter abuse.BlockedTargetFilter) ([]abuse.BlockedTarget, error) {
	query, args := buildBlockedTargetQuery(filter)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]abuse.BlockedTarget, 0)
	for rows.Next() {
		var item abuse.BlockedTarget
		if err := rows.Scan(
			&item.ID,
			&item.TargetType,
			&item.TargetValue,
			&item.Reason,
			&item.Source,
			&item.IsActive,
			&item.ExpiresAt,
			&item.CreatedByUserID,
			&item.Metadata,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *BlockedTargetRepository) Upsert(ctx context.Context, target abuse.BlockedTarget) (abuse.BlockedTarget, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO blocked_targets (
			target_type,
			target_value,
			reason,
			source,
			is_active,
			expires_at,
			created_by_user_id,
			metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (target_type, target_value)
		DO UPDATE SET
			reason = EXCLUDED.reason,
			source = EXCLUDED.source,
			is_active = EXCLUDED.is_active,
			expires_at = EXCLUDED.expires_at,
			created_by_user_id = EXCLUDED.created_by_user_id,
			metadata = EXCLUDED.metadata,
			updated_at = now()
		RETURNING id, target_type, target_value, reason, source, is_active, expires_at, created_by_user_id, metadata, created_at, updated_at
	`,
		target.TargetType,
		target.TargetValue,
		target.Reason,
		target.Source,
		target.IsActive,
		target.ExpiresAt,
		target.CreatedByUserID,
		target.Metadata,
	)

	var saved abuse.BlockedTarget
	if err := row.Scan(
		&saved.ID,
		&saved.TargetType,
		&saved.TargetValue,
		&saved.Reason,
		&saved.Source,
		&saved.IsActive,
		&saved.ExpiresAt,
		&saved.CreatedByUserID,
		&saved.Metadata,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	); err != nil {
		return abuse.BlockedTarget{}, err
	}

	return saved, nil
}

func buildBlockedTargetQuery(filter abuse.BlockedTargetFilter) (string, []any) {
	var clauses []string
	args := make([]any, 0, 4)

	addClause := func(expr string, value any) {
		args = append(args, value)
		clauses = append(clauses, fmt.Sprintf(expr, len(args)))
	}

	if filter.TargetType != "" {
		addClause("target_type = $%d", filter.TargetType)
	}
	if filter.TargetValue != "" {
		addClause("target_value = $%d", filter.TargetValue)
	}
	if filter.OnlyActive {
		addClause("is_active = $%d", true)
	}
	if !filter.IncludeExpired {
		clauses = append(clauses, "(expires_at IS NULL OR expires_at > now())")
	}

	limit := normalizeLimit(filter.Limit, 100, 500)
	args = append(args, limit)

	var builder strings.Builder
	builder.WriteString(`
		SELECT id, target_type, target_value, reason, source, is_active, expires_at, created_by_user_id, metadata, created_at, updated_at
		FROM blocked_targets
	`)
	if len(clauses) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(strings.Join(clauses, " AND "))
	}
	builder.WriteString(fmt.Sprintf(" ORDER BY is_active DESC, updated_at DESC, target_type ASC, target_value ASC LIMIT $%d", len(args)))

	return builder.String(), args
}
