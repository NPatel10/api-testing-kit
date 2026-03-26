package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"api-testing-kit/server/internal/billing"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionEventRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionEventRepository(pool *pgxpool.Pool) *SubscriptionEventRepository {
	return &SubscriptionEventRepository{pool: pool}
}

func (r *SubscriptionEventRepository) CreateSubscriptionEvent(ctx context.Context, params billing.EventUpsertParams) (billing.SubscriptionEvent, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO subscription_events (
			user_id,
			subscription_id,
			provider,
			provider_event_id,
			event_type,
			payload,
			processed_at,
			created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, user_id, subscription_id, provider, provider_event_id, event_type, payload, processed_at, created_at
	`,
		params.UserID,
		params.SubscriptionID,
		params.Provider,
		params.ProviderEventID,
		params.EventType,
		params.Payload,
		params.ProcessedAt,
		params.CreatedAt,
	)

	item, err := scanSubscriptionEvent(row.Scan)
	if err != nil {
		return billing.SubscriptionEvent{}, err
	}

	return item, nil
}

func (r *SubscriptionEventRepository) ListSubscriptionEvents(ctx context.Context, subscriptionID string, limit int32) ([]billing.SubscriptionEvent, error) {
	query, args := buildSubscriptionEventListQuery(subscriptionID, limit)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]billing.SubscriptionEvent, 0)
	for rows.Next() {
		item, err := scanSubscriptionEvent(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func buildSubscriptionEventListQuery(subscriptionID string, limit int32) (string, []any) {
	limit = normalizeLimit(limit, 25, 100)
	args := []any{strings.TrimSpace(subscriptionID), limit}

	return fmt.Sprintf(`
		SELECT id, user_id, subscription_id, provider, provider_event_id, event_type, payload, processed_at, created_at
		FROM subscription_events
		WHERE subscription_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`), args
}

func scanSubscriptionEvent(scan func(dest ...any) error) (billing.SubscriptionEvent, error) {
	var item billing.SubscriptionEvent
	var userID sql.NullString
	var processedAt sql.NullTime

	if err := scan(
		&item.ID,
		&userID,
		&item.SubscriptionID,
		&item.Provider,
		&item.ProviderEventID,
		&item.EventType,
		&item.Payload,
		&processedAt,
		&item.CreatedAt,
	); err != nil {
		return billing.SubscriptionEvent{}, err
	}

	if userID.Valid {
		value := userID.String
		item.UserID = &value
	}
	if processedAt.Valid {
		item.ProcessedAt = &processedAt.Time
	}

	return item, nil
}
