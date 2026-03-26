package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"api-testing-kit/server/internal/billing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{pool: pool}
}

func (r *SubscriptionRepository) UpsertSubscription(ctx context.Context, params billing.SubscriptionUpsertParams) (billing.Subscription, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO subscriptions (
			user_id,
			billing_customer_id,
			plan_id,
			provider,
			provider_subscription_id,
			status,
			quantity,
			trial_ends_at,
			current_period_start,
			current_period_end,
			cancel_at_period_end,
			canceled_at,
			paused_at,
			past_due_since,
			payment_method_hint,
			metadata,
			created_at,
			updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
		ON CONFLICT (provider, provider_subscription_id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			billing_customer_id = EXCLUDED.billing_customer_id,
			plan_id = EXCLUDED.plan_id,
			status = EXCLUDED.status,
			quantity = EXCLUDED.quantity,
			trial_ends_at = EXCLUDED.trial_ends_at,
			current_period_start = EXCLUDED.current_period_start,
			current_period_end = EXCLUDED.current_period_end,
			cancel_at_period_end = EXCLUDED.cancel_at_period_end,
			canceled_at = EXCLUDED.canceled_at,
			paused_at = EXCLUDED.paused_at,
			past_due_since = EXCLUDED.past_due_since,
			payment_method_hint = EXCLUDED.payment_method_hint,
			metadata = EXCLUDED.metadata,
			updated_at = EXCLUDED.updated_at
		RETURNING id, user_id, billing_customer_id, plan_id, provider, provider_subscription_id, status, quantity, trial_ends_at, current_period_start, current_period_end, cancel_at_period_end, canceled_at, paused_at, past_due_since, COALESCE(payment_method_hint, ''), metadata, created_at, updated_at
	`,
		params.UserID,
		params.BillingCustomerID,
		params.PlanID,
		params.Provider,
		params.ProviderSubscriptionID,
		params.Status,
		params.Quantity,
		params.TrialEndsAt,
		params.CurrentPeriodStart,
		params.CurrentPeriodEnd,
		params.CancelAtPeriodEnd,
		params.CanceledAt,
		params.PausedAt,
		params.PastDueSince,
		params.PaymentMethodHint,
		params.Metadata,
		params.CreatedAt,
		params.UpdatedAt,
	)

	item, err := scanSubscription(row.Scan)
	if err != nil {
		return billing.Subscription{}, err
	}

	return item, nil
}

func (r *SubscriptionRepository) GetSubscriptionByProvider(ctx context.Context, provider string, providerSubscriptionID string) (billing.Subscription, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, billing_customer_id, plan_id, provider, provider_subscription_id, status, quantity, trial_ends_at, current_period_start, current_period_end, cancel_at_period_end, canceled_at, paused_at, past_due_since, COALESCE(payment_method_hint, ''), metadata, created_at, updated_at
		FROM subscriptions
		WHERE provider = $1 AND provider_subscription_id = $2
	`, provider, providerSubscriptionID)

	item, err := scanSubscription(row.Scan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return billing.Subscription{}, billing.ErrNotFound
		}
		return billing.Subscription{}, err
	}

	return item, nil
}

func (r *SubscriptionRepository) ListSubscriptionsByUser(ctx context.Context, userID string, limit int32) ([]billing.Subscription, error) {
	query, args := buildSubscriptionListQuery(userID, limit)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]billing.Subscription, 0)
	for rows.Next() {
		item, err := scanSubscription(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func buildSubscriptionListQuery(userID string, limit int32) (string, []any) {
	limit = normalizeLimit(limit, 25, 100)
	args := []any{strings.TrimSpace(userID), limit}

	return fmt.Sprintf(`
		SELECT id, user_id, billing_customer_id, plan_id, provider, provider_subscription_id, status, quantity, trial_ends_at, current_period_start, current_period_end, cancel_at_period_end, canceled_at, paused_at, past_due_since, COALESCE(payment_method_hint, ''), metadata, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY updated_at DESC, created_at DESC
		LIMIT $2
	`), args
}

func scanSubscription(scan func(dest ...any) error) (billing.Subscription, error) {
	var item billing.Subscription
	var userID sql.NullString
	var billingCustomerID sql.NullString
	var planID sql.NullString
	var trialEndsAt sql.NullTime
	var currentPeriodStart sql.NullTime
	var currentPeriodEnd sql.NullTime
	var canceledAt sql.NullTime
	var pausedAt sql.NullTime
	var pastDueSince sql.NullTime

	if err := scan(
		&item.ID,
		&userID,
		&billingCustomerID,
		&planID,
		&item.Provider,
		&item.ProviderSubscriptionID,
		&item.Status,
		&item.Quantity,
		&trialEndsAt,
		&currentPeriodStart,
		&currentPeriodEnd,
		&item.CancelAtPeriodEnd,
		&canceledAt,
		&pausedAt,
		&pastDueSince,
		&item.PaymentMethodHint,
		&item.Metadata,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return billing.Subscription{}, err
	}

	if userID.Valid {
		value := userID.String
		item.UserID = &value
	}
	if billingCustomerID.Valid {
		value := billingCustomerID.String
		item.BillingCustomerID = &value
	}
	if planID.Valid {
		value := planID.String
		item.PlanID = &value
	}
	if trialEndsAt.Valid {
		item.TrialEndsAt = &trialEndsAt.Time
	}
	if currentPeriodStart.Valid {
		item.CurrentPeriodStart = &currentPeriodStart.Time
	}
	if currentPeriodEnd.Valid {
		item.CurrentPeriodEnd = &currentPeriodEnd.Time
	}
	if canceledAt.Valid {
		item.CanceledAt = &canceledAt.Time
	}
	if pausedAt.Valid {
		item.PausedAt = &pausedAt.Time
	}
	if pastDueSince.Valid {
		item.PastDueSince = &pastDueSince.Time
	}

	return item, nil
}
