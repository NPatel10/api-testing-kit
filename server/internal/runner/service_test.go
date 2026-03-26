package runner

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"api-testing-kit/server/internal/abuse"
	"api-testing-kit/server/internal/history"
	"api-testing-kit/server/internal/ratelimit"
	"api-testing-kit/server/internal/safety"
	"api-testing-kit/server/internal/usage"
)

type fakeTransport struct {
	response *http.Response
	err      error
	request  *http.Request
}

func (t *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.request = req
	if t.err != nil {
		return nil, t.err
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

type blockingTransport struct {
	started chan struct{}
	release chan struct{}
	request *http.Request
}

func (t *blockingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.request = req
	if t.started != nil {
		select {
		case t.started <- struct{}{}:
		default:
		}
	}
	if t.release != nil {
		<-t.release
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
		Request:    req,
	}, nil
}

type fakeHistoryRepo struct {
	items []history.RunRecord
}

func (r *fakeHistoryRepo) ListByUser(ctx context.Context, userID string, limit int32) ([]history.RunRecord, error) {
	return r.items, nil
}

func (r *fakeHistoryRepo) Create(ctx context.Context, params history.CreateParams) (history.RunRecord, error) {
	record := history.RunRecord{
		ID:            "run-1",
		UserID:        &params.UserID,
		Source:        params.Source,
		Status:        params.Status,
		Method:        params.Method,
		URL:           params.URL,
		BlockedReason: params.BlockedReason,
		ErrorCode:     params.ErrorCode,
		ErrorMessage:  params.ErrorMessage,
		Metadata:      params.Metadata,
	}
	r.items = append(r.items, record)
	return record, nil
}

type fakeUsageRecorder struct {
	events []usage.Event
}

func (r *fakeUsageRecorder) Create(ctx context.Context, event usage.Event) (usage.Event, error) {
	r.events = append(r.events, event)
	return event, nil
}

type fakeAbuseRecorder struct {
	events []abuse.Event
}

func (r *fakeAbuseRecorder) Create(ctx context.Context, event abuse.Event) (abuse.Event, error) {
	r.events = append(r.events, event)
	return event, nil
}

type fakeLimiter struct {
	ipDecision     ratelimit.Decision
	userDecision   ratelimit.Decision
	domainDecision ratelimit.Decision
}

func (l fakeLimiter) AllowIP(key string) (ratelimit.Decision, error)   { return l.ipDecision, nil }
func (l fakeLimiter) AllowUser(key string) (ratelimit.Decision, error) { return l.userDecision, nil }
func (l fakeLimiter) AllowDomain(key string) (ratelimit.Decision, error) {
	return l.domainDecision, nil
}

type fakeResolver struct {
	ips []net.IPAddr
}

func (r fakeResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	return r.ips, nil
}

func TestExecuteSuccess(t *testing.T) {
	t.Parallel()

	historyRepo := &fakeHistoryRepo{}
	client := &http.Client{
		Transport: &fakeTransport{
			response: &http.Response{
				StatusCode: 200,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
			},
		},
	}
	service := NewService(client, history.NewService(historyRepo), safety.Options{
		Resolver: fakeResolver{ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}}},
	})

	result, err := service.Execute(context.Background(), "user-1", "203.0.113.8", RunInput{
		Method: "GET",
		URL:    "https://api.example.com/test",
	})
	if err != nil {
		t.Fatalf("expected successful execution, got %v", err)
	}
	if result.ResponseStatus != 200 {
		t.Fatalf("expected status 200, got %d", result.ResponseStatus)
	}
	if len(historyRepo.items) != 1 {
		t.Fatalf("expected history entry to be recorded")
	}
}

func TestExecuteRejectsOversizeRequestBody(t *testing.T) {
	t.Parallel()

	historyRepo := &fakeHistoryRepo{}
	service := NewService(
		&http.Client{Transport: &fakeTransport{}},
		history.NewService(historyRepo),
		safety.Options{Resolver: fakeResolver{ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}}}},
		WithRequestBodyLimit(4),
	)

	_, err := service.Execute(context.Background(), "user-1", "203.0.113.8", RunInput{
		Method: "POST",
		URL:    "https://api.example.com/test",
		Body: BodyInput{
			Mode: "raw",
			Raw:  "hello",
		},
	})
	if !errors.Is(err, ErrRequestTooLarge) {
		t.Fatalf("expected request too large error, got %v", err)
	}
	if len(historyRepo.items) != 1 {
		t.Fatalf("expected blocked history entry to be recorded")
	}
	if historyRepo.items[0].Status != "blocked" {
		t.Fatalf("expected blocked history status, got %q", historyRepo.items[0].Status)
	}
}

func TestExecuteRejectsConcurrentRequests(t *testing.T) {
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	historyRepo := &fakeHistoryRepo{}
	usageRecorder := &fakeUsageRecorder{}
	abuseRecorder := &fakeAbuseRecorder{}
	service := NewService(
		&http.Client{Transport: &blockingTransport{started: started, release: release}},
		history.NewService(historyRepo),
		safety.Options{Resolver: fakeResolver{ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}}}},
		WithLimiter(fakeLimiter{
			ipDecision:     ratelimit.Decision{Allowed: true, Scope: ratelimit.ScopeIP, Key: "203.0.113.8"},
			userDecision:   ratelimit.Decision{Allowed: true, Scope: ratelimit.ScopeUser, Key: "user-1"},
			domainDecision: ratelimit.Decision{Allowed: true, Scope: ratelimit.ScopeDomain, Key: "api.example.com"},
		}),
		WithUsageRecorder(usageRecorder),
		WithAbuseRecorder(abuseRecorder),
		WithConcurrencyLimits(1, 1),
	)

	firstDone := make(chan error, 1)
	go func() {
		_, err := service.Execute(context.Background(), "user-1", "203.0.113.8", RunInput{
			Method: "GET",
			URL:    "https://api.example.com/test",
		})
		firstDone <- err
	}()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("first request never reached transport")
	}

	_, err := service.Execute(context.Background(), "user-1", "203.0.113.8", RunInput{
		Method: "GET",
		URL:    "https://api.example.com/test",
	})
	var limitErr *LimitError
	if !errors.As(err, &limitErr) || limitErr.Reason != "concurrency_limit" {
		t.Fatalf("expected concurrency limit error, got %v", err)
	}

	close(release)
	if err := <-firstDone; err != nil {
		t.Fatalf("expected first request to complete, got %v", err)
	}

	if len(usageRecorder.events) < 2 {
		t.Fatalf("expected usage events to be recorded for blocked and successful runs")
	}
	if len(abuseRecorder.events) < 1 {
		t.Fatalf("expected abuse event to be recorded for blocked run")
	}
}

func TestExecuteBlocksUnsafeDestination(t *testing.T) {
	t.Parallel()

	service := NewService(&http.Client{}, history.NewService(&fakeHistoryRepo{}), safety.Options{})
	_, err := service.Execute(context.Background(), "user-1", "203.0.113.8", RunInput{
		Method: "GET",
		URL:    "http://127.0.0.1/admin",
	})
	if err == nil {
		t.Fatalf("expected validation failure")
	}
}
