package ratelimit

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Scope string

const (
	ScopeIP     Scope = "ip"
	ScopeUser   Scope = "user"
	ScopeDomain Scope = "domain"
)

type Policy struct {
	Window      time.Duration
	Limit       int
	BurstWindow time.Duration
	BurstLimit  int
	Cooldown    time.Duration
}

type Config struct {
	IP     Policy
	User   Policy
	Domain Policy
	Clock  func() time.Time
}

type Decision struct {
	Allowed       bool
	Scope         Scope
	Key           string
	Reason        string
	Remaining     int
	ResetAt       time.Time
	RetryAfter    time.Duration
	CooldownUntil time.Time
}

type Snapshot struct {
	Scope         Scope
	Key           string
	Events        int
	BurstEvents   int
	Remaining     int
	ResetAt       time.Time
	CooldownUntil time.Time
}

type Limiter struct {
	mu       sync.Mutex
	now      func() time.Time
	policies map[Scope]Policy
	states   map[stateKey]*state
}

type stateKey struct {
	scope Scope
	key   string
}

type state struct {
	events        []time.Time
	burstEvents   []time.Time
	cooldownUntil time.Time
}

func DefaultConfig() Config {
	return Config{
		IP: Policy{
			Window:      10 * time.Minute,
			Limit:       10,
			BurstWindow: 15 * time.Second,
			BurstLimit:  3,
			Cooldown:    2 * time.Minute,
		},
		User: Policy{
			Window:      time.Hour,
			Limit:       60,
			BurstWindow: 30 * time.Second,
			BurstLimit:  6,
			Cooldown:    5 * time.Minute,
		},
		Domain: Policy{
			Window:      time.Minute,
			Limit:       12,
			BurstWindow: 15 * time.Second,
			BurstLimit:  4,
			Cooldown:    30 * time.Second,
		},
		Clock: time.Now,
	}
}

func GuestConfig() Config {
	cfg := DefaultConfig()
	cfg.IP = Policy{
		Window:      10 * time.Minute,
		Limit:       10,
		BurstWindow: 10 * time.Second,
		BurstLimit:  2,
		Cooldown:    2 * time.Minute,
	}
	cfg.Domain = Policy{
		Window:      time.Minute,
		Limit:       8,
		BurstWindow: 10 * time.Second,
		BurstLimit:  2,
		Cooldown:    30 * time.Second,
	}
	return cfg
}

func AuthenticatedConfig() Config {
	cfg := DefaultConfig()
	cfg.User = Policy{
		Window:      time.Hour,
		Limit:       60,
		BurstWindow: 30 * time.Second,
		BurstLimit:  5,
		Cooldown:    5 * time.Minute,
	}
	cfg.Domain = Policy{
		Window:      time.Minute,
		Limit:       12,
		BurstWindow: 15 * time.Second,
		BurstLimit:  4,
		Cooldown:    30 * time.Second,
	}
	return cfg
}

func NewLimiter(cfg Config) *Limiter {
	normalized := cfg
	if normalized.Clock == nil {
		normalized.Clock = time.Now
	}

	return &Limiter{
		now: normalized.Clock,
		policies: map[Scope]Policy{
			ScopeIP:     normalized.IP,
			ScopeUser:   normalized.User,
			ScopeDomain: normalized.Domain,
		},
		states: make(map[stateKey]*state),
	}
}

func (l *Limiter) AllowIP(key string) (Decision, error) {
	return l.Allow(ScopeIP, key)
}

func (l *Limiter) AllowUser(key string) (Decision, error) {
	return l.Allow(ScopeUser, key)
}

func (l *Limiter) AllowDomain(key string) (Decision, error) {
	return l.Allow(ScopeDomain, NormalizeDomain(key))
}

func (l *Limiter) Allow(scope Scope, key string) (Decision, error) {
	if l == nil {
		return Decision{}, fmt.Errorf("ratelimiter is nil")
	}

	policy, ok := l.policies[scope]
	if !ok {
		return Decision{}, fmt.Errorf("unsupported scope %q", scope)
	}

	normalizedKey := NormalizeKey(scope, key)
	if normalizedKey == "" {
		return Decision{}, fmt.Errorf("rate limit key is empty")
	}

	now := l.now().UTC()

	l.mu.Lock()
	defer l.mu.Unlock()

	st := l.stateFor(scope, normalizedKey)
	l.prune(stateTrimPolicy{
		now:         now,
		quotaWindow: policy.Window,
		burstWindow: policy.BurstWindow,
	}, st)

	if !st.cooldownUntil.IsZero() && now.Before(st.cooldownUntil) {
		return buildBlockedDecision(scope, normalizedKey, "cooldown", st.cooldownUntil, st.cooldownUntil, st.cooldownUntil.Sub(now), policy), nil
	}

	if policy.BurstLimit > 0 && policy.BurstWindow > 0 && len(st.burstEvents) >= policy.BurstLimit {
		blockUntil := maxTime(now.Add(policy.Cooldown), st.burstEvents[0].Add(policy.BurstWindow))
		st.cooldownUntil = blockUntil
		return buildBlockedDecision(scope, normalizedKey, "burst_limit", blockUntil, blockUntil, blockUntil.Sub(now), policy), nil
	}

	if policy.Limit > 0 && policy.Window > 0 && len(st.events) >= policy.Limit {
		blockUntil := maxTime(now.Add(policy.Cooldown), st.events[0].Add(policy.Window))
		st.cooldownUntil = blockUntil
		return buildBlockedDecision(scope, normalizedKey, "quota_limit", blockUntil, blockUntil, blockUntil.Sub(now), policy), nil
	}

	st.events = append(st.events, now)
	st.burstEvents = append(st.burstEvents, now)

	decision := Decision{
		Allowed:    true,
		Scope:      scope,
		Key:        normalizedKey,
		Remaining:  remainingAfter(policy.Limit, len(st.events)),
		ResetAt:    resetAt(st.events, policy.Window),
		RetryAfter: 0,
	}

	return decision, nil
}

func (l *Limiter) Snapshot(scope Scope, key string) (Snapshot, error) {
	if l == nil {
		return Snapshot{}, fmt.Errorf("ratelimiter is nil")
	}

	policy, ok := l.policies[scope]
	if !ok {
		return Snapshot{}, fmt.Errorf("unsupported scope %q", scope)
	}

	normalizedKey := NormalizeKey(scope, key)
	if normalizedKey == "" {
		return Snapshot{}, fmt.Errorf("rate limit key is empty")
	}

	now := l.now().UTC()

	l.mu.Lock()
	defer l.mu.Unlock()

	st := l.stateFor(scope, normalizedKey)
	l.prune(stateTrimPolicy{
		now:         now,
		quotaWindow: policy.Window,
		burstWindow: policy.BurstWindow,
	}, st)

	return Snapshot{
		Scope:         scope,
		Key:           normalizedKey,
		Events:        len(st.events),
		BurstEvents:   len(st.burstEvents),
		Remaining:     remainingAfter(policy.Limit, len(st.events)),
		ResetAt:       resetAt(st.events, policy.Window),
		CooldownUntil: st.cooldownUntil,
	}, nil
}

func DomainKeyFromURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := parsed.Hostname()
	if host == "" {
		return "", fmt.Errorf("missing host")
	}

	return NormalizeDomain(host), nil
}

func NormalizeDomain(domain string) string {
	return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(domain)), ".")
}

func NormalizeKey(scope Scope, key string) string {
	trimmed := strings.TrimSpace(key)
	if trimmed == "" {
		return ""
	}

	if scope == ScopeDomain {
		return NormalizeDomain(trimmed)
	}

	return trimmed
}

func (l *Limiter) stateFor(scope Scope, key string) *state {
	id := stateKey{scope: scope, key: key}
	if existing := l.states[id]; existing != nil {
		return existing
	}

	st := &state{}
	l.states[id] = st
	return st
}

type stateTrimPolicy struct {
	now         time.Time
	quotaWindow time.Duration
	burstWindow time.Duration
}

func (l *Limiter) prune(policy stateTrimPolicy, st *state) {
	if st == nil {
		return
	}

	if !st.cooldownUntil.IsZero() && !policy.now.Before(st.cooldownUntil) {
		st.cooldownUntil = time.Time{}
	}

	if policy.quotaWindow > 0 {
		cutoff := policy.now.Add(-policy.quotaWindow)
		st.events = trimTimes(st.events, cutoff)
	}

	if policy.burstWindow > 0 {
		cutoff := policy.now.Add(-policy.burstWindow)
		st.burstEvents = trimTimes(st.burstEvents, cutoff)
	}
}

func trimTimes(values []time.Time, cutoff time.Time) []time.Time {
	if len(values) == 0 {
		return nil
	}

	idx := 0
	for idx < len(values) && values[idx].Before(cutoff) {
		idx++
	}

	if idx == 0 {
		return values
	}

	if idx >= len(values) {
		return nil
	}

	return append([]time.Time(nil), values[idx:]...)
}

func buildBlockedDecision(scope Scope, key, reason string, resetAt, cooldownUntil time.Time, retryAfter time.Duration, policy Policy) Decision {
	remaining := 0
	if policy.Limit > 0 {
		remaining = 0
	}

	return Decision{
		Allowed:       false,
		Scope:         scope,
		Key:           key,
		Reason:        reason,
		Remaining:     remaining,
		ResetAt:       resetAt,
		RetryAfter:    retryAfter,
		CooldownUntil: cooldownUntil,
	}
}

func remainingAfter(limit, used int) int {
	if limit <= 0 {
		return -1
	}

	remaining := limit - used
	if remaining < 0 {
		return 0
	}

	return remaining
}

func resetAt(values []time.Time, window time.Duration) time.Time {
	if len(values) == 0 || window <= 0 {
		return time.Time{}
	}

	return values[0].Add(window)
}

func maxTime(a, b time.Time) time.Time {
	if a.IsZero() {
		return b
	}
	if b.IsZero() {
		return a
	}
	if a.After(b) {
		return a
	}
	return b
}
