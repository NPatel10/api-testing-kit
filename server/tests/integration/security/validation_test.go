package security_test

import (
	"context"
	"net"
	"testing"

	"api-testing-kit/server/internal/safety"
)

type resolverWithIP struct {
	ip net.IP
}

func (r resolverWithIP) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	return []net.IPAddr{{IP: append(net.IP(nil), r.ip...)}}, nil
}

func TestSecurityValidationBlocksLocalhost(t *testing.T) {
	t.Parallel()

	_, err := safety.ValidateURL(context.Background(), "http://localhost:8080/healthz", safety.DefaultOptions())
	if err == nil {
		t.Fatal("expected localhost URL to be blocked")
	}

	var validationErr *safety.ValidationError
	if !safety.AsValidationError(err, &validationErr) {
		t.Fatalf("expected validation error, got %T", err)
	}

	if validationErr.Code != safety.ErrorBlockedHost {
		t.Fatalf("unexpected validation code: %s", validationErr.Code)
	}
}

func TestSecurityValidationBlocksPrivateIPs(t *testing.T) {
	t.Parallel()

	_, err := safety.ValidateURL(context.Background(), "https://internal.example.com/status", safety.Options{
		AllowedSchemes: []string{"https"},
		AllowedPorts:   []int{443},
		Resolver:       resolverWithIP{ip: net.ParseIP("10.0.0.5")},
	})
	if err == nil {
		t.Fatal("expected private IP to be blocked")
	}

	var validationErr *safety.ValidationError
	if !safety.AsValidationError(err, &validationErr) {
		t.Fatalf("expected validation error, got %T", err)
	}

	if validationErr.Code != safety.ErrorBlockedIP {
		t.Fatalf("unexpected validation code: %s", validationErr.Code)
	}
}

func TestSecurityValidationBlocksUnsupportedSchemes(t *testing.T) {
	t.Parallel()

	_, err := safety.ValidateURL(context.Background(), "file:///etc/passwd", safety.DefaultOptions())
	if err == nil {
		t.Fatal("expected unsupported scheme to be blocked")
	}

	var validationErr *safety.ValidationError
	if !safety.AsValidationError(err, &validationErr) {
		t.Fatalf("expected validation error, got %T", err)
	}

	if validationErr.Code != safety.ErrorUnsupportedScheme {
		t.Fatalf("unexpected validation code: %s", validationErr.Code)
	}
}
