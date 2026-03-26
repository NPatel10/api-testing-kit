package entitlements

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrUnavailable = errors.New("entitlements repository is unavailable")
	ErrNotFound    = errors.New("plan not found")
	ErrInvalid     = errors.New("invalid entitlement input")
)

const (
	EntitlementKeyCustomURLExecution     = "custom_url_execution"
	EntitlementKeySavedHistoryDepth      = "saved_history_depth"
	EntitlementKeyEnvironmentVariables   = "environment_variables"
	EntitlementKeySharedLinks            = "shared_links"
	EntitlementKeyRequestTimeoutSeconds  = "request_timeout_seconds"
	EntitlementKeyMaxConcurrentRequests  = "max_concurrent_requests"
	EntitlementKeyRequestsPerHour        = "requests_per_hour"
	EntitlementKeyRequestsPerDay         = "requests_per_day"
	EntitlementKeyResponsePreviewBytes   = "response_preview_bytes"
	EntitlementKeyGuestRequestBodyBytes  = "guest_request_body_bytes"
	EntitlementKeyAuthenticatedBodyBytes = "authenticated_request_body_bytes"
	EntitlementKeyAuthenticatedRedirects = "authenticated_redirects"
)

type Plan struct {
	ID          string          `json:"id"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	IsActive    bool            `json:"isActive"`
	SortOrder   int             `json:"sortOrder"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type PlanEntitlement struct {
	ID        string          `json:"id"`
	PlanID    string          `json:"planId"`
	Key       string          `json:"key"`
	Value     json.RawMessage `json:"value"`
	CreatedAt time.Time       `json:"createdAt"`
}

type Snapshot struct {
	Plan         Plan                       `json:"plan"`
	Entitlements map[string]json.RawMessage `json:"entitlements"`
}

type AccessPolicy struct {
	CanUseCustomURLs       bool `json:"canUseCustomUrls"`
	CanUseEnvironmentVars  bool `json:"canUseEnvironmentVariables"`
	CanUseSharedLinks      bool `json:"canUseSharedLinks"`
	SavedHistoryDepth      int  `json:"savedHistoryDepth"`
	RequestTimeoutSeconds  int  `json:"requestTimeoutSeconds"`
	MaxConcurrentRequests  int  `json:"maxConcurrentRequests"`
	RequestsPerHour        int  `json:"requestsPerHour"`
	RequestsPerDay         int  `json:"requestsPerDay"`
	ResponsePreviewBytes   int  `json:"responsePreviewBytes"`
	GuestRequestBodyBytes  int  `json:"guestRequestBodyBytes"`
	AuthenticatedBodyBytes int  `json:"authenticatedRequestBodyBytes"`
	AuthenticatedRedirects int  `json:"authenticatedRedirects"`
}

type ResolvedPlan struct {
	Snapshot Snapshot     `json:"snapshot"`
	Policy   AccessPolicy `json:"policy"`
}

type Repository interface {
	GetPlanByCode(ctx context.Context, code string) (Plan, error)
	GetPlanByID(ctx context.Context, id string) (Plan, error)
	ListActivePlans(ctx context.Context) ([]Plan, error)
	ListEntitlementsByPlanID(ctx context.Context, planID string) ([]PlanEntitlement, error)
}

type Service struct {
	repo     Repository
	baseline AccessPolicy
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:     repo,
		baseline: DefaultAccessPolicy(),
	}
}

func (s *Service) WithBaseline(policy AccessPolicy) *Service {
	if s == nil {
		return nil
	}
	s.baseline = policy
	return s
}

func (s *Service) ResolveByCode(ctx context.Context, code string) (ResolvedPlan, error) {
	if s == nil || s.repo == nil {
		return ResolvedPlan{}, ErrUnavailable
	}

	code = strings.TrimSpace(code)
	if code == "" {
		return ResolvedPlan{}, ErrInvalid
	}

	plan, err := s.repo.GetPlanByCode(ctx, code)
	if err != nil {
		return ResolvedPlan{}, err
	}

	return s.resolve(ctx, plan)
}

func (s *Service) ResolveByID(ctx context.Context, id string) (ResolvedPlan, error) {
	if s == nil || s.repo == nil {
		return ResolvedPlan{}, ErrUnavailable
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return ResolvedPlan{}, ErrInvalid
	}

	plan, err := s.repo.GetPlanByID(ctx, id)
	if err != nil {
		return ResolvedPlan{}, err
	}

	return s.resolve(ctx, plan)
}

func (s *Service) ListActivePlans(ctx context.Context) ([]Plan, error) {
	if s == nil || s.repo == nil {
		return nil, ErrUnavailable
	}

	return s.repo.ListActivePlans(ctx)
}

func (s *Service) resolve(ctx context.Context, plan Plan) (ResolvedPlan, error) {
	ents, err := s.repo.ListEntitlementsByPlanID(ctx, plan.ID)
	if err != nil {
		return ResolvedPlan{}, err
	}

	snapshot := Snapshot{
		Plan:         plan,
		Entitlements: normalizeEntitlements(ents),
	}

	return ResolvedPlan{
		Snapshot: snapshot,
		Policy:   ResolveAccessPolicy(snapshot, s.baseline),
	}, nil
}

func DefaultAccessPolicy() AccessPolicy {
	return AccessPolicy{}
}

func ResolveAccessPolicy(snapshot Snapshot, baseline AccessPolicy) AccessPolicy {
	policy := baseline

	if value, ok := snapshot.Bool(EntitlementKeyCustomURLExecution); ok {
		policy.CanUseCustomURLs = value
	}
	if value, ok := snapshot.Bool(EntitlementKeyEnvironmentVariables); ok {
		policy.CanUseEnvironmentVars = value
	}
	if value, ok := snapshot.Bool(EntitlementKeySharedLinks); ok {
		policy.CanUseSharedLinks = value
	}
	if value, ok := snapshot.Int(EntitlementKeySavedHistoryDepth); ok {
		policy.SavedHistoryDepth = value
	}
	if value, ok := snapshot.Int(EntitlementKeyRequestTimeoutSeconds); ok {
		policy.RequestTimeoutSeconds = value
	}
	if value, ok := snapshot.Int(EntitlementKeyMaxConcurrentRequests); ok {
		policy.MaxConcurrentRequests = value
	}
	if value, ok := snapshot.Int(EntitlementKeyRequestsPerHour); ok {
		policy.RequestsPerHour = value
	}
	if value, ok := snapshot.Int(EntitlementKeyRequestsPerDay); ok {
		policy.RequestsPerDay = value
	}
	if value, ok := snapshot.Int(EntitlementKeyResponsePreviewBytes); ok {
		policy.ResponsePreviewBytes = value
	}
	if value, ok := snapshot.Int(EntitlementKeyGuestRequestBodyBytes); ok {
		policy.GuestRequestBodyBytes = value
	}
	if value, ok := snapshot.Int(EntitlementKeyAuthenticatedBodyBytes); ok {
		policy.AuthenticatedBodyBytes = value
	}
	if value, ok := snapshot.Int(EntitlementKeyAuthenticatedRedirects); ok {
		policy.AuthenticatedRedirects = value
	}

	return policy
}

func normalizeEntitlements(items []PlanEntitlement) map[string]json.RawMessage {
	if len(items) == 0 {
		return map[string]json.RawMessage{}
	}

	entitlements := make(map[string]json.RawMessage, len(items))
	for _, item := range items {
		key := strings.TrimSpace(item.Key)
		if key == "" || len(item.Value) == 0 {
			continue
		}

		copyValue := json.RawMessage(append([]byte(nil), item.Value...))
		entitlements[key] = copyValue
	}

	return entitlements
}

func (s Snapshot) Bool(key string) (bool, bool) {
	raw, ok := s.Entitlements[strings.TrimSpace(key)]
	if !ok {
		return false, false
	}

	var value bool
	if err := json.Unmarshal(raw, &value); err != nil {
		return false, false
	}

	return value, true
}

func (s Snapshot) Int(key string) (int, bool) {
	raw, ok := s.Entitlements[strings.TrimSpace(key)]
	if !ok {
		return 0, false
	}

	var value int
	if err := json.Unmarshal(raw, &value); err != nil {
		return 0, false
	}

	return value, true
}
