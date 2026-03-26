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

type BillingCustomerRepository struct {
	pool *pgxpool.Pool
}

func NewBillingCustomerRepository(pool *pgxpool.Pool) *BillingCustomerRepository {
	return &BillingCustomerRepository{pool: pool}
}

func (r *BillingCustomerRepository) UpsertCustomer(ctx context.Context, params billing.CustomerUpsertParams) (billing.Customer, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO billing_customers (
			user_id,
			provider,
			provider_customer_id,
			provider_email,
			tax_country,
			tax_id_last4,
			metadata,
			created_at,
			updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (provider, provider_customer_id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			provider_email = EXCLUDED.provider_email,
			tax_country = EXCLUDED.tax_country,
			tax_id_last4 = EXCLUDED.tax_id_last4,
			metadata = EXCLUDED.metadata,
			updated_at = EXCLUDED.updated_at
		RETURNING id, user_id, provider, provider_customer_id, provider_email, tax_country, tax_id_last4, metadata, created_at, updated_at
	`,
		params.UserID,
		params.Provider,
		params.ProviderCustomerID,
		params.ProviderEmail,
		params.TaxCountry,
		params.TaxIDLast4,
		params.Metadata,
		params.CreatedAt,
		params.UpdatedAt,
	)

	item, err := scanBillingCustomer(row.Scan)
	if err != nil {
		return billing.Customer{}, err
	}

	return item, nil
}

func (r *BillingCustomerRepository) GetCustomerByProvider(ctx context.Context, provider string, providerCustomerID string) (billing.Customer, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, provider, provider_customer_id, provider_email, tax_country, tax_id_last4, metadata, created_at, updated_at
		FROM billing_customers
		WHERE provider = $1 AND provider_customer_id = $2
	`, provider, providerCustomerID)

	item, err := scanBillingCustomer(row.Scan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return billing.Customer{}, billing.ErrNotFound
		}
		return billing.Customer{}, err
	}

	return item, nil
}

func (r *BillingCustomerRepository) ListCustomersByUser(ctx context.Context, userID string, limit int32) ([]billing.Customer, error) {
	query, args := buildBillingCustomerListQuery(userID, limit)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]billing.Customer, 0)
	for rows.Next() {
		item, err := scanBillingCustomer(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func buildBillingCustomerListQuery(userID string, limit int32) (string, []any) {
	limit = normalizeLimit(limit, 25, 100)
	args := []any{strings.TrimSpace(userID), limit}

	return fmt.Sprintf(`
		SELECT id, user_id, provider, provider_customer_id, provider_email, tax_country, tax_id_last4, metadata, created_at, updated_at
		FROM billing_customers
		WHERE user_id = $1
		ORDER BY updated_at DESC, created_at DESC
		LIMIT $2
	`), args
}

func scanBillingCustomer(scan func(dest ...any) error) (billing.Customer, error) {
	var item billing.Customer
	var userID sql.NullString
	var providerEmail sql.NullString
	var taxCountry sql.NullString
	var taxIDLast4 sql.NullString

	if err := scan(
		&item.ID,
		&userID,
		&item.Provider,
		&item.ProviderCustomerID,
		&providerEmail,
		&taxCountry,
		&taxIDLast4,
		&item.Metadata,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return billing.Customer{}, err
	}

	if userID.Valid {
		value := userID.String
		item.UserID = &value
	}
	if providerEmail.Valid {
		item.ProviderEmail = providerEmail.String
	}
	if taxCountry.Valid {
		item.TaxCountry = taxCountry.String
	}
	if taxIDLast4.Valid {
		item.TaxIDLast4 = taxIDLast4.String
	}

	return item, nil
}
