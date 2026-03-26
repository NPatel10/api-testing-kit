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

type InvoiceRepository struct {
	pool *pgxpool.Pool
}

func NewInvoiceRepository(pool *pgxpool.Pool) *InvoiceRepository {
	return &InvoiceRepository{pool: pool}
}

func (r *InvoiceRepository) UpsertInvoice(ctx context.Context, params billing.InvoiceUpsertParams) (billing.Invoice, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO invoices (
			user_id,
			subscription_id,
			provider,
			provider_invoice_id,
			invoice_number,
			status,
			amount_due_cents,
			amount_paid_cents,
			currency,
			hosted_invoice_url,
			pdf_url,
			issued_at,
			due_at,
			paid_at,
			metadata,
			created_at,
			updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		ON CONFLICT (provider, provider_invoice_id)
		DO UPDATE SET
			subscription_id = EXCLUDED.subscription_id,
			invoice_number = EXCLUDED.invoice_number,
			status = EXCLUDED.status,
			amount_due_cents = EXCLUDED.amount_due_cents,
			amount_paid_cents = EXCLUDED.amount_paid_cents,
			currency = EXCLUDED.currency,
			hosted_invoice_url = EXCLUDED.hosted_invoice_url,
			pdf_url = EXCLUDED.pdf_url,
			issued_at = EXCLUDED.issued_at,
			due_at = EXCLUDED.due_at,
			paid_at = EXCLUDED.paid_at,
			metadata = EXCLUDED.metadata,
			updated_at = EXCLUDED.updated_at
		RETURNING id, user_id, subscription_id, provider, provider_invoice_id, invoice_number, status, amount_due_cents, amount_paid_cents, currency, hosted_invoice_url, pdf_url, issued_at, due_at, paid_at, metadata, created_at, updated_at
	`,
		params.UserID,
		params.SubscriptionID,
		params.Provider,
		params.ProviderInvoiceID,
		params.InvoiceNumber,
		params.Status,
		params.AmountDueCents,
		params.AmountPaidCents,
		params.Currency,
		params.HostedInvoiceURL,
		params.PDFURL,
		params.IssuedAt,
		params.DueAt,
		params.PaidAt,
		params.Metadata,
		params.CreatedAt,
		params.UpdatedAt,
	)

	item, err := scanInvoice(row.Scan)
	if err != nil {
		return billing.Invoice{}, err
	}

	return item, nil
}

func (r *InvoiceRepository) GetInvoiceByProvider(ctx context.Context, provider string, providerInvoiceID string) (billing.Invoice, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, subscription_id, provider, provider_invoice_id, invoice_number, status, amount_due_cents, amount_paid_cents, currency, hosted_invoice_url, pdf_url, issued_at, due_at, paid_at, metadata, created_at, updated_at
		FROM invoices
		WHERE provider = $1 AND provider_invoice_id = $2
	`, provider, providerInvoiceID)

	item, err := scanInvoice(row.Scan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return billing.Invoice{}, billing.ErrNotFound
		}
		return billing.Invoice{}, err
	}

	return item, nil
}

func (r *InvoiceRepository) ListInvoicesBySubscription(ctx context.Context, subscriptionID string, limit int32) ([]billing.Invoice, error) {
	query, args := buildInvoiceListQuery(subscriptionID, limit)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]billing.Invoice, 0)
	for rows.Next() {
		item, err := scanInvoice(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func buildInvoiceListQuery(subscriptionID string, limit int32) (string, []any) {
	limit = normalizeLimit(limit, 25, 100)
	args := []any{strings.TrimSpace(subscriptionID), limit}

	return fmt.Sprintf(`
		SELECT id, user_id, subscription_id, provider, provider_invoice_id, invoice_number, status, amount_due_cents, amount_paid_cents, currency, hosted_invoice_url, pdf_url, issued_at, due_at, paid_at, metadata, created_at, updated_at
		FROM invoices
		WHERE subscription_id = $1
		ORDER BY COALESCE(issued_at, created_at) DESC
		LIMIT $2
	`), args
}

func scanInvoice(scan func(dest ...any) error) (billing.Invoice, error) {
	var item billing.Invoice
	var userID sql.NullString
	var subscriptionID sql.NullString
	var invoiceNumber sql.NullString
	var hostedInvoiceURL sql.NullString
	var pdfURL sql.NullString
	var issuedAt sql.NullTime
	var dueAt sql.NullTime
	var paidAt sql.NullTime

	if err := scan(
		&item.ID,
		&userID,
		&subscriptionID,
		&item.Provider,
		&item.ProviderInvoiceID,
		&invoiceNumber,
		&item.Status,
		&item.AmountDueCents,
		&item.AmountPaidCents,
		&item.Currency,
		&hostedInvoiceURL,
		&pdfURL,
		&issuedAt,
		&dueAt,
		&paidAt,
		&item.Metadata,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return billing.Invoice{}, err
	}

	if userID.Valid {
		value := userID.String
		item.UserID = &value
	}
	if subscriptionID.Valid {
		value := subscriptionID.String
		item.SubscriptionID = &value
	}
	if invoiceNumber.Valid {
		item.InvoiceNumber = invoiceNumber.String
	}
	if hostedInvoiceURL.Valid {
		item.HostedInvoiceURL = hostedInvoiceURL.String
	}
	if pdfURL.Valid {
		item.PDFURL = pdfURL.String
	}
	if issuedAt.Valid {
		item.IssuedAt = &issuedAt.Time
	}
	if dueAt.Valid {
		item.DueAt = &dueAt.Time
	}
	if paidAt.Valid {
		item.PaidAt = &paidAt.Time
	}

	return item, nil
}
