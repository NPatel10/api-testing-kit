package safety

import (
	"context"
	"net"
	"testing"
)

type fakeResolver struct {
	ips []net.IPAddr
	err error
}

func (r fakeResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	return r.ips, r.err
}

func TestValidateURL_AllowsPublicHTTPSDestination(t *testing.T) {
	t.Parallel()

	result, err := ValidateURL(context.Background(), "https://example.com/path", Options{
		Resolver: fakeResolver{
			ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}},
		},
	})
	if err != nil {
		t.Fatalf("expected allowed URL, got error: %v", err)
	}

	if result.ResolvedHost != "example.com" {
		t.Fatalf("expected host example.com, got %q", result.ResolvedHost)
	}

	if got := result.Port; got != 443 {
		t.Fatalf("expected default HTTPS port 443, got %d", got)
	}
}

func TestValidateURL_BlocksLocalhost(t *testing.T) {
	t.Parallel()

	_, err := ValidateURL(context.Background(), "http://localhost:8080", Options{})
	assertValidationError(t, err, ErrorBlockedHost)
}

func TestValidateURL_BlocksPrivateIP(t *testing.T) {
	t.Parallel()

	_, err := ValidateURL(context.Background(), "https://192.168.1.25", Options{})
	assertValidationError(t, err, ErrorBlockedIP)
}

func TestValidateURL_BlocksMetadataIP(t *testing.T) {
	t.Parallel()

	_, err := ValidateURL(context.Background(), "http://169.254.169.254/latest/meta-data", Options{})
	assertValidationError(t, err, ErrorBlockedIP)
}

func TestValidateURL_BlocksUnsupportedProtocol(t *testing.T) {
	t.Parallel()

	_, err := ValidateURL(context.Background(), "file:///etc/passwd", Options{})
	assertValidationError(t, err, ErrorUnsupportedScheme)
}

func TestValidateURL_BlocksNonStandardPort(t *testing.T) {
	t.Parallel()

	_, err := ValidateURL(context.Background(), "https://example.com:8443", Options{
		Resolver: fakeResolver{
			ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}},
		},
	})
	assertValidationError(t, err, ErrorPortNotAllowed)
}

func TestValidateURL_BlocksDNSResolvedPrivateIP(t *testing.T) {
	t.Parallel()

	_, err := ValidateURL(context.Background(), "https://api.internal.example", Options{
		Resolver: fakeResolver{
			ips: []net.IPAddr{{IP: net.ParseIP("10.0.0.5")}},
		},
	})
	assertValidationError(t, err, ErrorBlockedIP)
}

func TestValidateRedirectChain_BlocksTooManyRedirects(t *testing.T) {
	t.Parallel()

	_, err := ValidateRedirectChain(context.Background(), []string{
		"https://example.com/one",
		"https://example.com/two",
		"https://example.com/three",
		"https://example.com/four",
	}, Options{
		MaxRedirects: 2,
		Resolver: fakeResolver{
			ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}},
		},
	})
	assertValidationError(t, err, ErrorTooManyRedirects)
}

func TestValidateRedirectChain_BlocksUnsafeRedirectHop(t *testing.T) {
	t.Parallel()

	_, err := ValidateRedirectChain(context.Background(), []string{
		"https://example.com/start",
		"http://127.0.0.1/admin",
	}, Options{
		Resolver: fakeResolver{
			ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}},
		},
	})
	assertValidationError(t, err, ErrorBlockedIP)
}

func assertValidationError(t *testing.T, err error, code ErrorCode) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error %q, got nil", code)
	}

	var validationErr *ValidationError
	if !AsValidationError(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T: %v", err, err)
	}

	if validationErr.Code != code {
		t.Fatalf("expected error code %q, got %q", code, validationErr.Code)
	}
}
