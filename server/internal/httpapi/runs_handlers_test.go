package httpapi

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api-testing-kit/server/internal/auth"
	"api-testing-kit/server/internal/history"
	"api-testing-kit/server/internal/ratelimit"
	"api-testing-kit/server/internal/runner"
	"api-testing-kit/server/internal/safety"
)

type recordingLimiter struct {
	ipKeys     []string
	userKeys   []string
	domainKeys []string
}

func (l *recordingLimiter) AllowIP(key string) (ratelimit.Decision, error) {
	l.ipKeys = append(l.ipKeys, key)
	return ratelimit.Decision{Allowed: true, Scope: ratelimit.ScopeIP, Key: key}, nil
}

func (l *recordingLimiter) AllowUser(key string) (ratelimit.Decision, error) {
	l.userKeys = append(l.userKeys, key)
	return ratelimit.Decision{Allowed: true, Scope: ratelimit.ScopeUser, Key: key}, nil
}

func (l *recordingLimiter) AllowDomain(key string) (ratelimit.Decision, error) {
	l.domainKeys = append(l.domainKeys, key)
	return ratelimit.Decision{Allowed: true, Scope: ratelimit.ScopeDomain, Key: key}, nil
}

type runsRouteTransport struct {
	body string
}

func (t *runsRouteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		payload, _ := io.ReadAll(req.Body)
		t.body = string(payload)
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
		Request:    req,
	}, nil
}

type runsRouteResolver struct {
	ips []net.IPAddr
}

func (r runsRouteResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	return r.ips, nil
}

func TestRunRoutePassesClientIPToAuthenticatedRunner(t *testing.T) {
	authRepo := newFakeAuthRepo()
	authService := auth.NewService(authRepo)
	authResult, err := authService.Signup(context.Background(), auth.SignupInput{
		Email:    "runs@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("signup failed: %v", err)
	}

	limiter := &recordingLimiter{}
	transport := &runsRouteTransport{}
	service := runner.NewService(
		&http.Client{Transport: transport},
		history.NewService(newFakeHistoryRepo()),
		safety.Options{
			Resolver: runsRouteResolver{
				ips: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}},
			},
		},
		runner.WithLimiter(limiter),
	)

	mux := http.NewServeMux()
	NewRunsHandler(service, authService).Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/runs", strings.NewReader(`{"method":"GET","url":"https://api.example.com/test"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: authResult.Token})
	req.RemoteAddr = "203.0.113.8:4321"
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if len(limiter.ipKeys) != 1 || limiter.ipKeys[0] != "203.0.113.8" {
		t.Fatalf("expected client IP to be forwarded, got %#v", limiter.ipKeys)
	}
	if len(limiter.userKeys) != 1 {
		t.Fatalf("expected user limiter to be checked, got %#v", limiter.userKeys)
	}
	if len(limiter.domainKeys) != 1 || limiter.domainKeys[0] != "api.example.com" {
		t.Fatalf("expected domain limiter to be checked, got %#v", limiter.domainKeys)
	}
	if transport.body != "" {
		t.Fatalf("expected empty request body for GET, got %q", transport.body)
	}
}
