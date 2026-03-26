package billing

import (
	"context"
	"strings"
	"time"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
}

func (s *Service) UpsertCustomer(ctx context.Context, params CustomerUpsertParams) (Customer, error) {
	if s == nil || s.repo == nil {
		return Customer{}, ErrUnavailable
	}

	normalized, err := normalizeCustomerParams(params, s.now())
	if err != nil {
		return Customer{}, err
	}

	return s.repo.UpsertCustomer(ctx, normalized)
}

func (s *Service) UpsertSubscription(ctx context.Context, params SubscriptionUpsertParams) (Subscription, error) {
	if s == nil || s.repo == nil {
		return Subscription{}, ErrUnavailable
	}

	normalized, err := normalizeSubscriptionParams(params, s.now())
	if err != nil {
		return Subscription{}, err
	}

	return s.repo.UpsertSubscription(ctx, normalized)
}

func (s *Service) RecordSubscriptionEvent(ctx context.Context, params EventUpsertParams) (SubscriptionEvent, error) {
	if s == nil || s.repo == nil {
		return SubscriptionEvent{}, ErrUnavailable
	}

	normalized, err := normalizeEventParams(params, s.now())
	if err != nil {
		return SubscriptionEvent{}, err
	}

	return s.repo.CreateSubscriptionEvent(ctx, normalized)
}

func (s *Service) UpsertInvoice(ctx context.Context, params InvoiceUpsertParams) (Invoice, error) {
	if s == nil || s.repo == nil {
		return Invoice{}, ErrUnavailable
	}

	normalized, err := normalizeInvoiceParams(params, s.now())
	if err != nil {
		return Invoice{}, err
	}

	return s.repo.UpsertInvoice(ctx, normalized)
}

func (s *Service) ListSubscriptionInvoices(ctx context.Context, subscriptionID string, limit int32) ([]Invoice, error) {
	if s == nil || s.repo == nil {
		return nil, ErrUnavailable
	}

	subscriptionID = strings.TrimSpace(subscriptionID)
	if subscriptionID == "" {
		return nil, ErrInvalid
	}
	if limit <= 0 {
		limit = 25
	}

	return s.repo.ListInvoicesBySubscription(ctx, subscriptionID, limit)
}

func normalizeCustomerParams(params CustomerUpsertParams, now time.Time) (CustomerUpsertParams, error) {
	params.Provider = strings.TrimSpace(params.Provider)
	params.ProviderCustomerID = strings.TrimSpace(params.ProviderCustomerID)
	if params.ProviderCustomerID == "" {
		return CustomerUpsertParams{}, ErrInvalid
	}

	provider, err := normalizeProvider(params.Provider)
	if err != nil {
		return CustomerUpsertParams{}, err
	}
	params.Provider = string(provider)

	if params.ProviderEmail != nil {
		value := strings.TrimSpace(*params.ProviderEmail)
		params.ProviderEmail = &value
	}
	if params.TaxCountry != nil {
		value := strings.TrimSpace(*params.TaxCountry)
		params.TaxCountry = &value
	}
	if params.TaxIDLast4 != nil {
		value := strings.TrimSpace(*params.TaxIDLast4)
		params.TaxIDLast4 = &value
	}
	if params.CreatedAt.IsZero() {
		params.CreatedAt = now
	}
	if params.UpdatedAt.IsZero() {
		params.UpdatedAt = now
	}
	params.Metadata = normalizeJSON(params.Metadata)

	return params, nil
}

func normalizeSubscriptionParams(params SubscriptionUpsertParams, now time.Time) (SubscriptionUpsertParams, error) {
	params.Provider = strings.TrimSpace(params.Provider)
	params.ProviderSubscriptionID = strings.TrimSpace(params.ProviderSubscriptionID)
	params.Status = strings.TrimSpace(params.Status)
	params.PaymentMethodHint = strings.TrimSpace(params.PaymentMethodHint)
	if params.ProviderSubscriptionID == "" {
		return SubscriptionUpsertParams{}, ErrInvalid
	}

	provider, err := normalizeProvider(params.Provider)
	if err != nil {
		return SubscriptionUpsertParams{}, err
	}
	params.Provider = string(provider)

	status, err := normalizeSubscriptionStatus(params.Status)
	if err != nil {
		return SubscriptionUpsertParams{}, err
	}
	params.Status = string(status)

	if params.CreatedAt.IsZero() {
		params.CreatedAt = now
	}
	if params.UpdatedAt.IsZero() {
		params.UpdatedAt = now
	}
	params.Metadata = normalizeJSON(params.Metadata)

	return params, nil
}

func normalizeEventParams(params EventUpsertParams, now time.Time) (EventUpsertParams, error) {
	params.Provider = strings.TrimSpace(params.Provider)
	params.SubscriptionID = strings.TrimSpace(params.SubscriptionID)
	params.ProviderEventID = strings.TrimSpace(params.ProviderEventID)
	params.EventType = strings.TrimSpace(params.EventType)
	if params.SubscriptionID == "" || params.ProviderEventID == "" || params.EventType == "" {
		return EventUpsertParams{}, ErrInvalid
	}

	provider, err := normalizeProvider(params.Provider)
	if err != nil {
		return EventUpsertParams{}, err
	}
	params.Provider = string(provider)

	if params.CreatedAt.IsZero() {
		params.CreatedAt = now
	}
	params.Payload = normalizeJSON(params.Payload)

	return params, nil
}

func normalizeInvoiceParams(params InvoiceUpsertParams, now time.Time) (InvoiceUpsertParams, error) {
	params.Provider = strings.TrimSpace(params.Provider)
	params.ProviderInvoiceID = strings.TrimSpace(params.ProviderInvoiceID)
	params.Status = strings.TrimSpace(params.Status)
	params.Currency = strings.ToLower(strings.TrimSpace(params.Currency))
	if params.ProviderInvoiceID == "" || params.Status == "" || params.Currency == "" {
		return InvoiceUpsertParams{}, ErrInvalid
	}

	provider, err := normalizeProvider(params.Provider)
	if err != nil {
		return InvoiceUpsertParams{}, err
	}
	params.Provider = string(provider)

	if params.InvoiceNumber != nil {
		value := strings.TrimSpace(*params.InvoiceNumber)
		params.InvoiceNumber = &value
	}
	if params.HostedInvoiceURL != nil {
		value := strings.TrimSpace(*params.HostedInvoiceURL)
		params.HostedInvoiceURL = &value
	}
	if params.PDFURL != nil {
		value := strings.TrimSpace(*params.PDFURL)
		params.PDFURL = &value
	}
	if params.CreatedAt.IsZero() {
		params.CreatedAt = now
	}
	if params.UpdatedAt.IsZero() {
		params.UpdatedAt = now
	}
	params.Metadata = normalizeJSON(params.Metadata)

	return params, nil
}

const ProviderStatusUnfinalized = "provider_choice_not_finalized"

type StubService struct {
	ProviderStatus string
	SentAt         func() time.Time
}

type CheckoutHook struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	NextAction   string `json:"nextAction"`
	ProviderNote string `json:"providerNote"`
}

type WebhookReceipt struct {
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	ProviderNote string    `json:"providerNote"`
	ReceivedAt   time.Time `json:"receivedAt"`
}

func NewStubService() *StubService {
	return &StubService{
		ProviderStatus: ProviderStatusUnfinalized,
		SentAt:         time.Now,
	}
}

func (s *StubService) CheckoutSuccess() CheckoutHook {
	return CheckoutHook{
		Status:       "checkout_success_stub",
		Message:      "Checkout success handling is stubbed until a billing provider is chosen.",
		NextAction:   "return_to_app",
		ProviderNote: s.providerNote(),
	}
}

func (s *StubService) CheckoutCancel() CheckoutHook {
	return CheckoutHook{
		Status:       "checkout_cancel_stub",
		Message:      "Checkout cancel handling is stubbed until a billing provider is chosen.",
		NextAction:   "return_to_pricing",
		ProviderNote: s.providerNote(),
	}
}

func (s *StubService) WebhookReceived() WebhookReceipt {
	return WebhookReceipt{
		Status:       "webhook_received_stub",
		Message:      "Billing webhooks are accepted as a stub while the provider remains undecided.",
		ProviderNote: s.providerNote(),
		ReceivedAt:   s.now(),
	}
}

func (s *StubService) providerNote() string {
	if s == nil || s.ProviderStatus == "" {
		return ProviderStatusUnfinalized
	}

	return s.ProviderStatus
}

func (s *StubService) now() time.Time {
	if s != nil && s.SentAt != nil {
		return s.SentAt().UTC()
	}

	return time.Now().UTC()
}
