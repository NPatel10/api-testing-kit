package billing

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrUnavailable = errors.New("billing repository is unavailable")
	ErrNotFound    = errors.New("billing record not found")
	ErrInvalid     = errors.New("invalid billing input")
)

type Provider string

const (
	ProviderStripe       Provider = "stripe"
	ProviderPaddle       Provider = "paddle"
	ProviderLemonSqueezy Provider = "lemonsqueezy"
	ProviderManual       Provider = "manual"
)

type SubscriptionStatus string

const (
	SubscriptionStatusTrialing   SubscriptionStatus = "trialing"
	SubscriptionStatusActive     SubscriptionStatus = "active"
	SubscriptionStatusPastDue    SubscriptionStatus = "past_due"
	SubscriptionStatusPaused     SubscriptionStatus = "paused"
	SubscriptionStatusCanceled   SubscriptionStatus = "canceled"
	SubscriptionStatusIncomplete SubscriptionStatus = "incomplete"
	SubscriptionStatusUnpaid     SubscriptionStatus = "unpaid"
)

type Customer struct {
	ID                 string          `json:"id"`
	UserID             *string         `json:"userId,omitempty"`
	Provider           Provider        `json:"provider"`
	ProviderCustomerID string          `json:"providerCustomerId"`
	ProviderEmail      string          `json:"providerEmail,omitempty"`
	TaxCountry         string          `json:"taxCountry,omitempty"`
	TaxIDLast4         string          `json:"taxIdLast4,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	CreatedAt          time.Time       `json:"createdAt"`
	UpdatedAt          time.Time       `json:"updatedAt"`
}

type Subscription struct {
	ID                     string             `json:"id"`
	UserID                 *string            `json:"userId,omitempty"`
	BillingCustomerID      *string            `json:"billingCustomerId,omitempty"`
	PlanID                 *string            `json:"planId,omitempty"`
	Provider               Provider           `json:"provider"`
	ProviderSubscriptionID string             `json:"providerSubscriptionId"`
	Status                 SubscriptionStatus `json:"status"`
	Quantity               int                `json:"quantity"`
	TrialEndsAt            *time.Time         `json:"trialEndsAt,omitempty"`
	CurrentPeriodStart     *time.Time         `json:"currentPeriodStart,omitempty"`
	CurrentPeriodEnd       *time.Time         `json:"currentPeriodEnd,omitempty"`
	CancelAtPeriodEnd      bool               `json:"cancelAtPeriodEnd"`
	CanceledAt             *time.Time         `json:"canceledAt,omitempty"`
	PausedAt               *time.Time         `json:"pausedAt,omitempty"`
	PastDueSince           *time.Time         `json:"pastDueSince,omitempty"`
	PaymentMethodHint      string             `json:"paymentMethodHint,omitempty"`
	Metadata               json.RawMessage    `json:"metadata,omitempty"`
	CreatedAt              time.Time          `json:"createdAt"`
	UpdatedAt              time.Time          `json:"updatedAt"`
}

type SubscriptionEvent struct {
	ID              string          `json:"id"`
	UserID          *string         `json:"userId,omitempty"`
	SubscriptionID  string          `json:"subscriptionId"`
	Provider        Provider        `json:"provider"`
	ProviderEventID string          `json:"providerEventId"`
	EventType       string          `json:"eventType"`
	Payload         json.RawMessage `json:"payload,omitempty"`
	ProcessedAt     *time.Time      `json:"processedAt,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
}

type Invoice struct {
	ID                string          `json:"id"`
	UserID            *string         `json:"userId,omitempty"`
	SubscriptionID    *string         `json:"subscriptionId,omitempty"`
	Provider          Provider        `json:"provider"`
	ProviderInvoiceID string          `json:"providerInvoiceId"`
	InvoiceNumber     string          `json:"invoiceNumber,omitempty"`
	Status            string          `json:"status"`
	AmountDueCents    int             `json:"amountDueCents"`
	AmountPaidCents   int             `json:"amountPaidCents"`
	Currency          string          `json:"currency"`
	HostedInvoiceURL  string          `json:"hostedInvoiceUrl,omitempty"`
	PDFURL            string          `json:"pdfUrl,omitempty"`
	IssuedAt          *time.Time      `json:"issuedAt,omitempty"`
	DueAt             *time.Time      `json:"dueAt,omitempty"`
	PaidAt            *time.Time      `json:"paidAt,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
}

type CustomerUpsertParams struct {
	UserID             *string
	Provider           string
	ProviderCustomerID string
	ProviderEmail      *string
	TaxCountry         *string
	TaxIDLast4         *string
	Metadata           json.RawMessage
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type SubscriptionUpsertParams struct {
	UserID                 *string
	BillingCustomerID      *string
	PlanID                 *string
	Provider               string
	ProviderSubscriptionID string
	Status                 string
	Quantity               int
	TrialEndsAt            *time.Time
	CurrentPeriodStart     *time.Time
	CurrentPeriodEnd       *time.Time
	CancelAtPeriodEnd      bool
	CanceledAt             *time.Time
	PausedAt               *time.Time
	PastDueSince           *time.Time
	PaymentMethodHint      string
	Metadata               json.RawMessage
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type EventUpsertParams struct {
	UserID          *string
	SubscriptionID  string
	Provider        string
	ProviderEventID string
	EventType       string
	Payload         json.RawMessage
	ProcessedAt     *time.Time
	CreatedAt       time.Time
}

type InvoiceUpsertParams struct {
	UserID            *string
	SubscriptionID    *string
	Provider          string
	ProviderInvoiceID string
	InvoiceNumber     *string
	Status            string
	AmountDueCents    int
	AmountPaidCents   int
	Currency          string
	HostedInvoiceURL  *string
	PDFURL            *string
	IssuedAt          *time.Time
	DueAt             *time.Time
	PaidAt            *time.Time
	Metadata          json.RawMessage
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Repository interface {
	UpsertCustomer(ctx context.Context, params CustomerUpsertParams) (Customer, error)
	GetCustomerByProvider(ctx context.Context, provider string, providerCustomerID string) (Customer, error)
	ListCustomersByUser(ctx context.Context, userID string, limit int32) ([]Customer, error)
	UpsertSubscription(ctx context.Context, params SubscriptionUpsertParams) (Subscription, error)
	GetSubscriptionByProvider(ctx context.Context, provider string, providerSubscriptionID string) (Subscription, error)
	ListSubscriptionsByUser(ctx context.Context, userID string, limit int32) ([]Subscription, error)
	CreateSubscriptionEvent(ctx context.Context, params EventUpsertParams) (SubscriptionEvent, error)
	ListSubscriptionEvents(ctx context.Context, subscriptionID string, limit int32) ([]SubscriptionEvent, error)
	UpsertInvoice(ctx context.Context, params InvoiceUpsertParams) (Invoice, error)
	GetInvoiceByProvider(ctx context.Context, provider string, providerInvoiceID string) (Invoice, error)
	ListInvoicesBySubscription(ctx context.Context, subscriptionID string, limit int32) ([]Invoice, error)
}

func normalizeProvider(value string) (Provider, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return "", ErrInvalid
	}

	switch Provider(value) {
	case ProviderStripe, ProviderPaddle, ProviderLemonSqueezy, ProviderManual:
		return Provider(value), nil
	default:
		return "", ErrInvalid
	}
}

func normalizeSubscriptionStatus(value string) (SubscriptionStatus, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return "", ErrInvalid
	}

	switch SubscriptionStatus(value) {
	case SubscriptionStatusTrialing, SubscriptionStatusActive, SubscriptionStatusPastDue, SubscriptionStatusPaused, SubscriptionStatusCanceled, SubscriptionStatusIncomplete, SubscriptionStatusUnpaid:
		return SubscriptionStatus(value), nil
	default:
		return "", ErrInvalid
	}
}

func normalizeJSON(value json.RawMessage) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage(`{}`)
	}

	return value
}
