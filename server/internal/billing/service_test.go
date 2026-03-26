package billing

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

type fakeRepo struct {
	customerParams     CustomerUpsertParams
	subscriptionParams SubscriptionUpsertParams
	eventParams        EventUpsertParams
	invoiceParams      InvoiceUpsertParams
}

func (f *fakeRepo) UpsertCustomer(_ context.Context, params CustomerUpsertParams) (Customer, error) {
	f.customerParams = params
	return Customer{ID: "cust-1", Provider: Provider(params.Provider), ProviderCustomerID: params.ProviderCustomerID, CreatedAt: params.CreatedAt, UpdatedAt: params.UpdatedAt}, nil
}

func (f *fakeRepo) GetCustomerByProvider(context.Context, string, string) (Customer, error) {
	return Customer{}, nil
}

func (f *fakeRepo) ListCustomersByUser(context.Context, string, int32) ([]Customer, error) {
	return nil, nil
}

func (f *fakeRepo) UpsertSubscription(_ context.Context, params SubscriptionUpsertParams) (Subscription, error) {
	f.subscriptionParams = params
	return Subscription{ID: "sub-1", Provider: Provider(params.Provider), ProviderSubscriptionID: params.ProviderSubscriptionID, Status: SubscriptionStatus(params.Status), CreatedAt: params.CreatedAt, UpdatedAt: params.UpdatedAt}, nil
}

func (f *fakeRepo) GetSubscriptionByProvider(context.Context, string, string) (Subscription, error) {
	return Subscription{}, nil
}

func (f *fakeRepo) ListSubscriptionsByUser(context.Context, string, int32) ([]Subscription, error) {
	return nil, nil
}

func (f *fakeRepo) CreateSubscriptionEvent(_ context.Context, params EventUpsertParams) (SubscriptionEvent, error) {
	f.eventParams = params
	return SubscriptionEvent{ID: "evt-1", Provider: Provider(params.Provider), ProviderEventID: params.ProviderEventID, EventType: params.EventType, CreatedAt: params.CreatedAt}, nil
}

func (f *fakeRepo) ListSubscriptionEvents(context.Context, string, int32) ([]SubscriptionEvent, error) {
	return nil, nil
}

func (f *fakeRepo) UpsertInvoice(_ context.Context, params InvoiceUpsertParams) (Invoice, error) {
	f.invoiceParams = params
	return Invoice{ID: "inv-1", Provider: Provider(params.Provider), ProviderInvoiceID: params.ProviderInvoiceID, Status: params.Status, Currency: params.Currency, CreatedAt: params.CreatedAt, UpdatedAt: params.UpdatedAt}, nil
}

func (f *fakeRepo) GetInvoiceByProvider(context.Context, string, string) (Invoice, error) {
	return Invoice{}, nil
}

func (f *fakeRepo) ListInvoicesBySubscription(context.Context, string, int32) ([]Invoice, error) {
	return nil, nil
}

func TestServiceNormalizesAndDefaultsBillingInputs(t *testing.T) {
	t.Parallel()

	repo := &fakeRepo{}
	svc := NewService(repo)
	now := time.Date(2026, 3, 26, 10, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	email := "  billing@example.dev  "
	customer, err := svc.UpsertCustomer(context.Background(), CustomerUpsertParams{
		Provider:           "  STRIPE ",
		ProviderCustomerID: "  cus_123  ",
		ProviderEmail:      &email,
	})
	if err != nil {
		t.Fatalf("upsert customer: %v", err)
	}

	if customer.Provider != ProviderStripe {
		t.Fatalf("expected provider %q, got %q", ProviderStripe, customer.Provider)
	}

	if repo.customerParams.ProviderEmail == nil || *repo.customerParams.ProviderEmail != "billing@example.dev" {
		t.Fatalf("expected trimmed provider email, got %#v", repo.customerParams.ProviderEmail)
	}

	if repo.customerParams.CreatedAt != now || repo.customerParams.UpdatedAt != now {
		t.Fatalf("expected timestamp defaults to use now")
	}

	subscription, err := svc.UpsertSubscription(context.Background(), SubscriptionUpsertParams{
		Provider:               "paddle",
		ProviderSubscriptionID: " sub_123 ",
		Status:                 "active",
		PaymentMethodHint:      " card ",
	})
	if err != nil {
		t.Fatalf("upsert subscription: %v", err)
	}

	if subscription.Status != SubscriptionStatusActive {
		t.Fatalf("expected active status, got %q", subscription.Status)
	}

	if repo.subscriptionParams.PaymentMethodHint != "card" {
		t.Fatalf("expected trimmed payment method hint, got %q", repo.subscriptionParams.PaymentMethodHint)
	}

	event, err := svc.RecordSubscriptionEvent(context.Background(), EventUpsertParams{
		SubscriptionID:  "sub-1",
		Provider:        "manual",
		ProviderEventID: "evt-1",
		EventType:       "invoice.paid",
	})
	if err != nil {
		t.Fatalf("record event: %v", err)
	}

	if event.Provider != ProviderManual {
		t.Fatalf("expected manual provider, got %q", event.Provider)
	}

	invoice, err := svc.UpsertInvoice(context.Background(), InvoiceUpsertParams{
		Provider:          "lemonsqueezy",
		ProviderInvoiceID: "inv_123",
		Status:            "paid",
		Currency:          "USD",
	})
	if err != nil {
		t.Fatalf("upsert invoice: %v", err)
	}

	if invoice.Currency != "usd" {
		t.Fatalf("expected normalized currency, got %q", invoice.Currency)
	}
}

func TestServiceRejectsInvalidBillingInputs(t *testing.T) {
	t.Parallel()

	svc := NewService(&fakeRepo{})
	if _, err := svc.UpsertCustomer(context.Background(), CustomerUpsertParams{}); err == nil {
		t.Fatalf("expected customer validation error")
	}
	if _, err := svc.UpsertSubscription(context.Background(), SubscriptionUpsertParams{Provider: "stripe"}); err == nil {
		t.Fatalf("expected subscription validation error")
	}
	if _, err := svc.RecordSubscriptionEvent(context.Background(), EventUpsertParams{Provider: "stripe"}); err == nil {
		t.Fatalf("expected event validation error")
	}
	if _, err := svc.UpsertInvoice(context.Background(), InvoiceUpsertParams{Provider: "stripe", Currency: "usd"}); err == nil {
		t.Fatalf("expected invoice validation error")
	}
}

func TestNormalizeJSONDefault(t *testing.T) {
	t.Parallel()

	if got := normalizeJSON(nil); string(got) != "{}" {
		t.Fatalf("expected empty metadata default, got %s", string(got))
	}

	if got := normalizeJSON(json.RawMessage(`{"a":1}`)); string(got) != `{"a":1}` {
		t.Fatalf("expected payload to pass through, got %s", string(got))
	}
}
