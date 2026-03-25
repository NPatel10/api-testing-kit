package ratelimit

import (
	"testing"
	"time"
)

func TestAllowAppliesIndependentScopes(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 25, 10, 0, 0, 0, time.UTC)
	limiter := NewLimiter(Config{
		IP: Policy{
			Window:      time.Minute,
			Limit:       2,
			BurstWindow: 10 * time.Second,
			BurstLimit:  10,
			Cooldown:    30 * time.Second,
		},
		User: Policy{
			Window:      time.Minute,
			Limit:       2,
			BurstWindow: 10 * time.Second,
			BurstLimit:  10,
			Cooldown:    30 * time.Second,
		},
		Domain: Policy{
			Window:      time.Minute,
			Limit:       2,
			BurstWindow: 10 * time.Second,
			BurstLimit:  10,
			Cooldown:    30 * time.Second,
		},
		Clock: func() time.Time { return now },
	})

	first, err := limiter.AllowIP("192.0.2.10")
	if err != nil {
		t.Fatalf("allow ip: %v", err)
	}
	if !first.Allowed || first.Remaining != 1 {
		t.Fatalf("expected first ip request to be allowed with 1 remaining, got %+v", first)
	}

	second, err := limiter.AllowIP("192.0.2.10")
	if err != nil {
		t.Fatalf("allow ip second: %v", err)
	}
	if !second.Allowed || second.Remaining != 0 {
		t.Fatalf("expected second ip request to be allowed with 0 remaining, got %+v", second)
	}

	third, err := limiter.AllowIP("192.0.2.10")
	if err != nil {
		t.Fatalf("allow ip third: %v", err)
	}
	if third.Allowed || third.Reason != "quota_limit" {
		t.Fatalf("expected quota limit block, got %+v", third)
	}

	user, err := limiter.AllowUser("user-123")
	if err != nil {
		t.Fatalf("allow user: %v", err)
	}
	if !user.Allowed {
		t.Fatalf("expected user scope to remain independent, got %+v", user)
	}
}

func TestBurstLimitTriggersCooldown(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 25, 10, 0, 0, 0, time.UTC)
	limiter := NewLimiter(Config{
		IP: Policy{
			Window:      time.Minute,
			Limit:       10,
			BurstWindow: 5 * time.Second,
			BurstLimit:  2,
			Cooldown:    20 * time.Second,
		},
		Clock: func() time.Time { return now },
	})

	for i := 0; i < 2; i++ {
		decision, err := limiter.AllowIP("198.51.100.1")
		if err != nil {
			t.Fatalf("allow ip %d: %v", i, err)
		}
		if !decision.Allowed {
			t.Fatalf("expected request %d to be allowed, got %+v", i, decision)
		}
		now = now.Add(time.Second)
	}

	blocked, err := limiter.AllowIP("198.51.100.1")
	if err != nil {
		t.Fatalf("allow ip blocked: %v", err)
	}
	if blocked.Allowed || blocked.Reason != "burst_limit" {
		t.Fatalf("expected burst limit block, got %+v", blocked)
	}
	if blocked.CooldownUntil.Sub(now) < 20*time.Second {
		t.Fatalf("expected cooldown to extend beyond configured penalty, got %+v", blocked)
	}

	now = now.Add(21 * time.Second)
	allowed, err := limiter.AllowIP("198.51.100.1")
	if err != nil {
		t.Fatalf("allow ip after cooldown: %v", err)
	}
	if !allowed.Allowed {
		t.Fatalf("expected request after cooldown to be allowed, got %+v", allowed)
	}
}

func TestQuotaLimitUsesWindowReset(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 25, 10, 0, 0, 0, time.UTC)
	limiter := NewLimiter(Config{
		User: Policy{
			Window:      time.Minute,
			Limit:       2,
			BurstWindow: 10 * time.Second,
			BurstLimit:  10,
			Cooldown:    15 * time.Second,
		},
		Clock: func() time.Time { return now },
	})

	for i := 0; i < 2; i++ {
		decision, err := limiter.AllowUser("user-123")
		if err != nil {
			t.Fatalf("allow user %d: %v", i, err)
		}
		if !decision.Allowed {
			t.Fatalf("expected request %d to be allowed, got %+v", i, decision)
		}
		now = now.Add(5 * time.Second)
	}

	blocked, err := limiter.AllowUser("user-123")
	if err != nil {
		t.Fatalf("allow user blocked: %v", err)
	}
	if blocked.Allowed || blocked.Reason != "quota_limit" {
		t.Fatalf("expected quota limit block, got %+v", blocked)
	}
	if want := now.Add(50 * time.Second); !blocked.CooldownUntil.Equal(want) {
		t.Fatalf("expected retry window at %s, got %s", want, blocked.CooldownUntil)
	}

	now = now.Add(51 * time.Second)
	allowed, err := limiter.AllowUser("user-123")
	if err != nil {
		t.Fatalf("allow user after reset: %v", err)
	}
	if !allowed.Allowed {
		t.Fatalf("expected request after reset to be allowed, got %+v", allowed)
	}
}

func TestDomainKeyHelpers(t *testing.T) {
	t.Parallel()

	key, err := DomainKeyFromURL("https://Example.com/path?q=1")
	if err != nil {
		t.Fatalf("domain key from url: %v", err)
	}
	if key != "example.com" {
		t.Fatalf("expected normalized domain, got %q", key)
	}

	limiter := NewLimiter(Config{
		Domain: Policy{
			Window:      time.Minute,
			Limit:       1,
			BurstWindow: 10 * time.Second,
			BurstLimit:  1,
			Cooldown:    10 * time.Second,
		},
		Clock: func() time.Time {
			return time.Date(2026, time.March, 25, 10, 0, 0, 0, time.UTC)
		},
	})

	first, err := limiter.AllowDomain("Example.com")
	if err != nil {
		t.Fatalf("allow domain: %v", err)
	}
	if !first.Allowed {
		t.Fatalf("expected first domain request to be allowed, got %+v", first)
	}

	second, err := limiter.AllowDomain("example.com")
	if err != nil {
		t.Fatalf("allow domain second: %v", err)
	}
	if second.Allowed || second.Reason != "burst_limit" && second.Reason != "quota_limit" {
		t.Fatalf("expected repeated domain request to be throttled, got %+v", second)
	}
}
