package runner

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"api-testing-kit/server/internal/history"
	"api-testing-kit/server/internal/safety"
)

type fakeTransport struct {
	response *http.Response
	err      error
}

func (t fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.response, t.err
}

type fakeHistoryRepo struct {
	items []history.RunRecord
}

func (r *fakeHistoryRepo) ListByUser(ctx context.Context, userID string, limit int32) ([]history.RunRecord, error) {
	return r.items, nil
}

func (r *fakeHistoryRepo) Create(ctx context.Context, params history.CreateParams) (history.RunRecord, error) {
	record := history.RunRecord{ID: "run-1", Status: params.Status, Method: params.Method, URL: params.URL}
	r.items = append(r.items, record)
	return record, nil
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
		Transport: fakeTransport{
			response: &http.Response{
				StatusCode: 200,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
				Request:    mustRequest("GET", "https://api.example.com/test"),
			},
		},
	}
	service := NewService(client, history.NewService(historyRepo), safety.Options{
		Resolver: fakeResolver{ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}}},
	})

	result, err := service.Execute(context.Background(), "user-1", RunInput{
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

func TestExecuteBlocksUnsafeDestination(t *testing.T) {
	t.Parallel()

	service := NewService(&http.Client{}, history.NewService(&fakeHistoryRepo{}), safety.Options{})
	_, err := service.Execute(context.Background(), "user-1", RunInput{
		Method: "GET",
		URL:    "http://127.0.0.1/admin",
	})
	if err == nil {
		t.Fatalf("expected validation failure")
	}
}

func mustRequest(method string, rawURL string) *http.Request {
	req, _ := http.NewRequest(method, rawURL, nil)
	return req
}
