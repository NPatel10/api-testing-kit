package httpapi

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api-testing-kit/server/internal/abuse"
	"api-testing-kit/server/internal/guest"
	"api-testing-kit/server/internal/ratelimit"
	"api-testing-kit/server/internal/safety"
	"api-testing-kit/server/internal/usage"
)

type httpGuestUsageRecorder struct {
	events []usage.Event
}

func (r *httpGuestUsageRecorder) Create(ctx context.Context, event usage.Event) (usage.Event, error) {
	r.events = append(r.events, event)
	return event, nil
}

type httpGuestAbuseRecorder struct {
	events []abuse.Event
}

func (r *httpGuestAbuseRecorder) Create(ctx context.Context, event abuse.Event) (abuse.Event, error) {
	r.events = append(r.events, event)
	return event, nil
}

type httpGuestTransport struct {
	body string
}

func (t *httpGuestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		payload, _ := io.ReadAll(req.Body)
		t.body = string(payload)
	}

	return &http.Response{
		StatusCode: http.StatusCreated,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"status":"created"}`)),
		Request:    req,
	}, nil
}

type httpGuestResolver struct {
	ips []net.IPAddr
}

func (r httpGuestResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	return r.ips, nil
}

func TestGuestRunRouteReturnsResult(t *testing.T) {
	t.Parallel()

	transport := &httpGuestTransport{}
	service := guest.NewService(
		&http.Client{Transport: transport},
		ratelimit.NewLimiter(ratelimit.GuestConfig()),
		&httpGuestUsageRecorder{},
		&httpGuestAbuseRecorder{},
		safety.Options{
			Resolver: httpGuestResolver{
				ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}},
			},
		},
	)

	mux := http.NewServeMux()
	NewGuestRunsHandler(service).Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/guest-runs", strings.NewReader(`{"templateSlug":"auth-flow-mock","overrides":{"body":{"raw":"{\"email\":\"guest@example.dev\"}"}}}`))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "203.0.113.8:1234"
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if got, want := transport.body, `{"email":"guest@example.dev"}`; got != want {
		t.Fatalf("expected request body %q, got %q", want, got)
	}
}

func TestGuestRunRouteRejectsDisallowedOverride(t *testing.T) {
	t.Parallel()

	service := guest.NewService(
		&http.Client{},
		ratelimit.NewLimiter(ratelimit.GuestConfig()),
		nil,
		nil,
		safety.Options{
			Resolver: httpGuestResolver{
				ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}},
			},
		},
	)

	mux := http.NewServeMux()
	NewGuestRunsHandler(service).Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/guest-runs", strings.NewReader(`{"templateSlug":"jsonplaceholder-posts","overrides":{"queryParams":{"page":"2"}}}`))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "203.0.113.8:1234"
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	if got := rr.Body.String(); !strings.Contains(got, "guest_override_not_allowed") {
		t.Fatalf("expected override error payload, got %s", got)
	}
}
