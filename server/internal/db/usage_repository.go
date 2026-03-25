package db

import (
	"context"
	"fmt"
	"strings"

	"api-testing-kit/server/internal/usage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsageRepository struct {
	pool *pgxpool.Pool
}

func NewUsageRepository(pool *pgxpool.Pool) *UsageRepository {
	return &UsageRepository{pool: pool}
}

func (r *UsageRepository) Create(ctx context.Context, event usage.Event) (usage.Event, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO usage_events (
			user_id,
			session_id,
			request_run_id,
			bucket,
			event_key,
			quantity,
			dimensions,
			occurred_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, session_id, request_run_id, bucket, event_key, quantity, dimensions, occurred_at
	`,
		event.UserID,
		event.SessionID,
		event.RequestRunID,
		event.Bucket,
		event.EventKey,
		event.Quantity,
		event.Dimensions,
		event.OccurredAt,
	)

	var created usage.Event
	if err := row.Scan(
		&created.ID,
		&created.UserID,
		&created.SessionID,
		&created.RequestRunID,
		&created.Bucket,
		&created.EventKey,
		&created.Quantity,
		&created.Dimensions,
		&created.OccurredAt,
	); err != nil {
		return usage.Event{}, err
	}

	return created, nil
}

func (r *UsageRepository) ListRecent(ctx context.Context, filter usage.RecentFilter) ([]usage.Event, error) {
	query, args := buildUsageRecentQuery(filter)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]usage.Event, 0)
	for rows.Next() {
		var item usage.Event
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.SessionID,
			&item.RequestRunID,
			&item.Bucket,
			&item.EventKey,
			&item.Quantity,
			&item.Dimensions,
			&item.OccurredAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *UsageRepository) SummarizeByEventKey(ctx context.Context, filter usage.SummaryFilter) ([]usage.SummaryRow, error) {
	query, args := buildUsageSummaryQuery(filter)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]usage.SummaryRow, 0)
	for rows.Next() {
		var item usage.SummaryRow
		if err := rows.Scan(
			&item.Bucket,
			&item.EventKey,
			&item.EventCount,
			&item.TotalQuantity,
			&item.FirstOccurred,
			&item.LastOccurred,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func buildUsageRecentQuery(filter usage.RecentFilter) (string, []any) {
	var clauses []string
	args := make([]any, 0, 8)

	addClause := func(expr string, value any) {
		args = append(args, value)
		clauses = append(clauses, fmt.Sprintf(expr, len(args)))
	}

	if filter.UserID != nil {
		addClause("user_id = $%d", *filter.UserID)
	}
	if filter.SessionID != nil {
		addClause("session_id = $%d", *filter.SessionID)
	}
	if filter.RequestRunID != nil {
		addClause("request_run_id = $%d", *filter.RequestRunID)
	}
	if filter.Bucket != "" {
		addClause("bucket = $%d", filter.Bucket)
	}
	if filter.EventKey != "" {
		addClause("event_key = $%d", filter.EventKey)
	}
	if filter.OccurredAfter != nil {
		addClause("occurred_at >= $%d", *filter.OccurredAfter)
	}
	if filter.OccurredBefore != nil {
		addClause("occurred_at <= $%d", *filter.OccurredBefore)
	}

	limit := normalizeLimit(filter.Limit, 100, 500)
	args = append(args, limit)

	var builder strings.Builder
	builder.WriteString(`
		SELECT id, user_id, session_id, request_run_id, bucket, event_key, quantity, dimensions, occurred_at
		FROM usage_events
	`)
	if len(clauses) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(strings.Join(clauses, " AND "))
	}
	builder.WriteString(fmt.Sprintf(" ORDER BY occurred_at DESC LIMIT $%d", len(args)))

	return builder.String(), args
}

func buildUsageSummaryQuery(filter usage.SummaryFilter) (string, []any) {
	var clauses []string
	args := make([]any, 0, 8)

	addClause := func(expr string, value any) {
		args = append(args, value)
		clauses = append(clauses, fmt.Sprintf(expr, len(args)))
	}

	if filter.UserID != nil {
		addClause("user_id = $%d", *filter.UserID)
	}
	if filter.SessionID != nil {
		addClause("session_id = $%d", *filter.SessionID)
	}
	if filter.RequestRunID != nil {
		addClause("request_run_id = $%d", *filter.RequestRunID)
	}
	if filter.Bucket != "" {
		addClause("bucket = $%d", filter.Bucket)
	}
	if filter.EventKey != "" {
		addClause("event_key = $%d", filter.EventKey)
	}
	if filter.OccurredAfter != nil {
		addClause("occurred_at >= $%d", *filter.OccurredAfter)
	}
	if filter.OccurredBefore != nil {
		addClause("occurred_at <= $%d", *filter.OccurredBefore)
	}

	limit := normalizeLimit(filter.Limit, 50, 200)
	args = append(args, limit)

	var builder strings.Builder
	builder.WriteString(`
		SELECT bucket, event_key, COUNT(*) AS event_count, COALESCE(SUM(quantity), 0) AS total_quantity, MIN(occurred_at), MAX(occurred_at)
		FROM usage_events
	`)
	if len(clauses) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(strings.Join(clauses, " AND "))
	}
	builder.WriteString(`
		GROUP BY bucket, event_key
		ORDER BY MAX(occurred_at) DESC, bucket ASC, event_key ASC
		LIMIT $`)
	builder.WriteString(fmt.Sprint(len(args)))

	return builder.String(), args
}
