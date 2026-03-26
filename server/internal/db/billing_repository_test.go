package db

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"api-testing-kit/server/internal/billing"
)

func TestBuildBillingCustomerListQueryDefaultsLimit(t *testing.T) {
	t.Parallel()

	query, args := buildBillingCustomerListQuery(" user-123 ", 0)

	if got, want := len(args), 2; got != want {
		t.Fatalf("expected %d args, got %d", want, got)
	}
	if got, want := args[0], "user-123"; got != want {
		t.Fatalf("expected trimmed user id %v, got %v", want, got)
	}
	if got, want := args[1], int32(25); got != want {
		t.Fatalf("expected default limit %v, got %v", want, got)
	}
	if !strings.Contains(query, "FROM billing_customers") {
		t.Fatalf("expected billing_customers table in query")
	}
}

func TestBuildSubscriptionListQueryUsesUserFilter(t *testing.T) {
	t.Parallel()

	query, args := buildSubscriptionListQuery("user-123", 15)

	if got, want := len(args), 2; got != want {
		t.Fatalf("expected %d args, got %d", want, got)
	}
	if !strings.Contains(query, "WHERE user_id = $1") {
		t.Fatalf("expected user filter in query, got %q", query)
	}
	if !strings.Contains(query, "LIMIT $2") {
		t.Fatalf("expected limit placeholder in query, got %q", query)
	}
}

func TestBuildSubscriptionEventListQueryTrimsID(t *testing.T) {
	t.Parallel()

	query, args := buildSubscriptionEventListQuery(" sub-1 ", 0)

	if got, want := args[0], "sub-1"; got != want {
		t.Fatalf("expected trimmed subscription id %v, got %v", want, got)
	}
	if !strings.Contains(query, "FROM subscription_events") {
		t.Fatalf("expected subscription_events table in query")
	}
}

func TestBuildInvoiceListQueryOrdersByIssueDate(t *testing.T) {
	t.Parallel()

	query, args := buildInvoiceListQuery("sub-1", 50)

	if got, want := len(args), 2; got != want {
		t.Fatalf("expected %d args, got %d", want, got)
	}
	if !strings.Contains(query, "COALESCE(issued_at, created_at) DESC") {
		t.Fatalf("expected issue date ordering, got %q", query)
	}
}

func TestScanBillingCustomer(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 26, 10, 0, 0, 0, time.UTC)
	item, err := scanBillingCustomer(func(dest ...any) error {
		*dest[0].(*string) = "cust-1"
		*dest[2].(*billing.Provider) = billing.ProviderStripe
		*dest[3].(*string) = "cus_123"
		*dest[7].(*json.RawMessage) = json.RawMessage(`{"tier":"pro"}`)
		*dest[8].(*time.Time) = now
		*dest[9].(*time.Time) = now
		return nil
	})
	if err != nil {
		t.Fatalf("scan customer: %v", err)
	}
	if item.ProviderCustomerID != "cus_123" {
		t.Fatalf("unexpected provider customer id: %q", item.ProviderCustomerID)
	}
}
