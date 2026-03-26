package guest

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"api-testing-kit/server/internal/abuse"
	"api-testing-kit/server/internal/ratelimit"
	"api-testing-kit/server/internal/safety"
	"api-testing-kit/server/internal/usage"
)

type guestFakeTransport struct {
	response *http.Response
	request  *http.Request
	body     string
}

func (t *guestFakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.request = req
	if req.Body != nil {
		payload, _ := io.ReadAll(req.Body)
		t.body = string(payload)
	}

	if t.response == nil {
		t.response = &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
		}
	}

	copyResponse := *t.response
	copyResponse.Request = req
	return &copyResponse, nil
}

type guestFakeUsageRecorder struct {
	events []usage.Event
}

func (r *guestFakeUsageRecorder) Create(ctx context.Context, event usage.Event) (usage.Event, error) {
	r.events = append(r.events, event)
	return event, nil
}

type guestFakeAbuseRecorder struct {
	events []abuse.Event
}

func (r *guestFakeAbuseRecorder) Create(ctx context.Context, event abuse.Event) (abuse.Event, error) {
	r.events = append(r.events, event)
	return event, nil
}

type guestFakeResolver struct {
	ips []net.IPAddr
}

func (r guestFakeResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	return r.ips, nil
}

func TestExecuteGuestTemplateSuccess(t *testing.T) {
	t.Parallel()

	transport := &guestFakeTransport{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"id":1,"title":"demo"}`)),
		},
	}
	usageRecorder := &guestFakeUsageRecorder{}
	abuseRecorder := &guestFakeAbuseRecorder{}
	service := NewService(
		&http.Client{Transport: transport},
		ratelimit.NewLimiter(ratelimit.GuestConfig()),
		usageRecorder,
		abuseRecorder,
		safety.Options{
			Resolver: guestFakeResolver{
				ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}},
			},
		},
	)

	result, err := service.Execute(context.Background(), "203.0.113.8", RunRequest{
		TemplateSlug: "jsonplaceholder-posts",
		Overrides: &RunOverrides{
			QueryParams: map[string]string{"userid": "2"},
			Headers:     map[string]string{"accept": "application/json"},
		},
	})
	if err != nil {
		t.Fatalf("expected guest run to succeed, got %v", err)
	}

	if got, want := result.ResponseStatus, http.StatusOK; got != want {
		t.Fatalf("expected status %d, got %d", want, got)
	}
	if !strings.Contains(result.URL, "userId=2") {
		t.Fatalf("expected overridden query string in final url, got %q", result.URL)
	}
	if len(usageRecorder.events) != 1 {
		t.Fatalf("expected one usage event, got %d", len(usageRecorder.events))
	}
	if len(abuseRecorder.events) != 0 {
		t.Fatalf("expected no abuse events, got %d", len(abuseRecorder.events))
	}
	if got := transport.request.Method; got != http.MethodGet {
		t.Fatalf("expected GET request, got %s", got)
	}
}

func TestExecuteGuestTemplateAllowsBodyOverride(t *testing.T) {
	t.Parallel()

	transport := &guestFakeTransport{
		response: &http.Response{
			StatusCode: http.StatusCreated,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"status":"created"}`)),
		},
	}
	service := NewService(
		&http.Client{Transport: transport},
		ratelimit.NewLimiter(ratelimit.GuestConfig()),
		nil,
		nil,
		safety.Options{
			Resolver: guestFakeResolver{
				ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}},
			},
		},
	)

	result, err := service.Execute(context.Background(), "203.0.113.8", RunRequest{
		TemplateSlug: "auth-flow-mock",
		Overrides: &RunOverrides{
			Body: &BodyOverride{
				Raw: `{"email":"guest@example.dev"}`,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected body override to be allowed, got %v", err)
	}

	if got, want := result.ResponseStatus, http.StatusCreated; got != want {
		t.Fatalf("expected status %d, got %d", want, got)
	}
	if got, want := transport.body, `{"email":"guest@example.dev"}`; got != want {
		t.Fatalf("expected body %q, got %q", want, got)
	}
}

func TestExecuteRejectsDisallowedOverride(t *testing.T) {
	t.Parallel()

	service := NewService(
		&http.Client{},
		ratelimit.NewLimiter(ratelimit.GuestConfig()),
		nil,
		nil,
		safety.Options{
			Resolver: guestFakeResolver{
				ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}},
			},
		},
	)

	_, err := service.Execute(context.Background(), "203.0.113.8", RunRequest{
		TemplateSlug: "jsonplaceholder-posts",
		Overrides: &RunOverrides{
			QueryParams: map[string]string{"page": "2"},
		},
	})
	if err == nil {
		t.Fatalf("expected disallowed override to fail")
	}

	var guestErr *GuestError
	if !errors.As(err, &guestErr) {
		t.Fatalf("expected GuestError, got %T: %v", err, err)
	}
	if guestErr.Code != ErrorOverrideNotAllowed {
		t.Fatalf("expected override error, got %s", guestErr.Code)
	}
}
