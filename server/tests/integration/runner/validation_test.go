package runner_test

import (
	"context"
	"net"
	"testing"

	"api-testing-kit/server/internal/safety"
)

type staticResolver struct {
	ips []net.IP
}

func (r staticResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	addrs := make([]net.IPAddr, 0, len(r.ips))
	for _, ip := range r.ips {
		addrs = append(addrs, net.IPAddr{IP: append(net.IP(nil), ip...)})
	}

	return addrs, nil
}

func TestRunnerValidationAllowsPublicHTTPS(t *testing.T) {
	t.Parallel()

	result, err := safety.ValidateURL(context.Background(), "https://api.example.com/status", safety.Options{
		AllowedSchemes: []string{"https"},
		AllowedPorts:   []int{443},
		Resolver:       staticResolver{ips: []net.IP{net.ParseIP("93.184.216.34")}},
	})
	if err != nil {
		t.Fatalf("ValidateURL returned error: %v", err)
	}

	if result.URL == nil || result.URL.Scheme != "https" {
		t.Fatalf("unexpected parsed URL: %+v", result.URL)
	}
}

func TestRunnerValidationRejectsRedirectChainsPastLimit(t *testing.T) {
	t.Parallel()

	_, err := safety.ValidateRedirectChain(context.Background(), []string{
		"https://api.example.com/step-1",
		"https://api.example.com/step-2",
		"https://api.example.com/step-3",
		"https://api.example.com/step-4",
	}, safety.Options{
		AllowedSchemes: []string{"https"},
		AllowedPorts:   []int{443},
		Resolver:       staticResolver{ips: []net.IP{net.ParseIP("93.184.216.34")}},
		MaxRedirects:   2,
	})
	if err == nil {
		t.Fatal("expected redirect chain validation to fail")
	}

	var validationErr *safety.ValidationError
	if !safety.AsValidationError(err, &validationErr) {
		t.Fatalf("expected validation error, got %T", err)
	}

	if validationErr.Code != safety.ErrorTooManyRedirects {
		t.Fatalf("unexpected validation code: %s", validationErr.Code)
	}
}

func TestRunnerValidationRejectsInvalidURL(t *testing.T) {
	t.Parallel()

	_, err := safety.ValidateURL(context.Background(), "://broken", safety.DefaultOptions())
	if err == nil {
		t.Fatal("expected invalid URL to fail")
	}

	var validationErr *safety.ValidationError
	if !safety.AsValidationError(err, &validationErr) {
		t.Fatalf("expected validation error, got %T", err)
	}

	if validationErr.Code != safety.ErrorInvalidURL {
		t.Fatalf("unexpected validation code: %s", validationErr.Code)
	}
}
