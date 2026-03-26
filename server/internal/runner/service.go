package runner

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"api-testing-kit/server/internal/abuse"
	"api-testing-kit/server/internal/history"
	"api-testing-kit/server/internal/ratelimit"
	"api-testing-kit/server/internal/safety"
	"api-testing-kit/server/internal/usage"
)

var (
	ErrUnavailable      = errors.New("runner is unavailable")
	ErrInvalid          = errors.New("invalid run payload")
	ErrRequestTooLarge  = errors.New("request body too large")
	ErrTimedOut         = errors.New("request timed out")
	ErrConcurrencyLimit = errors.New("runner concurrency limit exceeded")
)

type KeyValue struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type AuthInput struct {
	Scheme   string `json:"scheme"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

type BodyInput struct {
	Mode       string     `json:"mode"`
	Raw        string     `json:"raw,omitempty"`
	FormFields []KeyValue `json:"formFields,omitempty"`
}

type RunInput struct {
	Method      string     `json:"method"`
	URL         string     `json:"url"`
	QueryParams []KeyValue `json:"queryParams,omitempty"`
	Headers     []KeyValue `json:"headers,omitempty"`
	Auth        AuthInput  `json:"auth"`
	Body        BodyInput  `json:"body"`
}

type RunResult struct {
	RunID               string              `json:"runId,omitempty"`
	Status              string              `json:"status"`
	Method              string              `json:"method"`
	URL                 string              `json:"url"`
	FinalURL            string              `json:"finalUrl,omitempty"`
	ResponseStatus      int                 `json:"responseStatus,omitempty"`
	ResponseHeaders     map[string][]string `json:"responseHeaders,omitempty"`
	ResponseBody        string              `json:"responseBodyPreview,omitempty"`
	ResponseJSON        any                 `json:"responseJson,omitempty"`
	ResponseSizeBytes   int                 `json:"responseSizeBytes,omitempty"`
	ResponseTimeMS      int                 `json:"responseTimeMs,omitempty"`
	ResponseContentType string              `json:"responseContentType,omitempty"`
	RedirectCount       int                 `json:"redirectCount"`
	BlockedReason       string              `json:"blockedReason,omitempty"`
	ErrorCode           string              `json:"errorCode,omitempty"`
	ErrorMessage        string              `json:"errorMessage,omitempty"`
	Truncated           bool                `json:"truncated"`
}

type UsageRecorder interface {
	Create(ctx context.Context, event usage.Event) (usage.Event, error)
}

type AbuseRecorder interface {
	Create(ctx context.Context, event abuse.Event) (abuse.Event, error)
}

type RateLimiter interface {
	AllowIP(key string) (ratelimit.Decision, error)
	AllowUser(key string) (ratelimit.Decision, error)
	AllowDomain(key string) (ratelimit.Decision, error)
}

type Option func(*Service)

func WithLimiter(limiter RateLimiter) Option {
	return func(s *Service) {
		s.limiter = limiter
	}
}

func WithUsageRecorder(recorder UsageRecorder) Option {
	return func(s *Service) {
		s.usageRecorder = recorder
	}
}

func WithAbuseRecorder(recorder AbuseRecorder) Option {
	return func(s *Service) {
		s.abuseRecorder = recorder
	}
}

func WithRequestTimeout(timeout time.Duration) Option {
	return func(s *Service) {
		if timeout > 0 {
			s.requestTimeout = timeout
		}
	}
}

func WithRequestBodyLimit(limit int) Option {
	return func(s *Service) {
		if limit > 0 {
			s.maxRequestBodyBytes = limit
		}
	}
}

func WithResponsePreviewLimit(limit int) Option {
	return func(s *Service) {
		if limit > 0 {
			s.maxPreviewBytes = limit
		}
	}
}

func WithConcurrencyLimits(userLimit, ipLimit int) Option {
	return func(s *Service) {
		if userLimit > 0 {
			s.maxConcurrentPerUser = userLimit
		}
		if ipLimit > 0 {
			s.maxConcurrentPerIP = ipLimit
		}
	}
}

type Service struct {
	client               *http.Client
	history              *history.Service
	limiter              RateLimiter
	usageRecorder        UsageRecorder
	abuseRecorder        AbuseRecorder
	safetyOptions        safety.Options
	requestTimeout       time.Duration
	maxPreviewBytes      int
	maxRequestBodyBytes  int
	maxConcurrentPerUser int
	maxConcurrentPerIP   int
	mu                   sync.Mutex
	activeByUser         map[string]int
	activeByIP           map[string]int
}

type LimitError struct {
	Scope      ratelimit.Scope
	Key        string
	Reason     string
	RetryAfter time.Duration
	Message    string
}

func (e *LimitError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Message
}

func NewService(client *http.Client, historyService *history.Service, opts safety.Options, options ...Option) *Service {
	if client == nil {
		client = &http.Client{}
	}

	defaults := safety.DefaultOptions()
	normalized := opts
	if len(normalized.AllowedSchemes) == 0 {
		normalized.AllowedSchemes = defaults.AllowedSchemes
	}
	if len(normalized.AllowedPorts) == 0 {
		normalized.AllowedPorts = defaults.AllowedPorts
	}
	if normalized.MaxRedirects == 0 {
		normalized.MaxRedirects = defaults.MaxRedirects
	}
	if normalized.Resolver == nil {
		normalized.Resolver = defaults.Resolver
	}

	service := &Service{
		client:               client,
		history:              historyService,
		limiter:              ratelimit.NewLimiter(ratelimit.AuthenticatedConfig()),
		safetyOptions:        normalized,
		requestTimeout:       15 * time.Second,
		maxPreviewBytes:      1024 * 1024,
		maxRequestBodyBytes:  256 * 1024,
		maxConcurrentPerUser: 5,
		maxConcurrentPerIP:   1,
		activeByUser:         make(map[string]int),
		activeByIP:           make(map[string]int),
	}

	for _, option := range options {
		if option != nil {
			option(service)
		}
	}

	if service.limiter == nil {
		service.limiter = ratelimit.NewLimiter(ratelimit.AuthenticatedConfig())
	}
	if service.requestTimeout <= 0 {
		service.requestTimeout = 15 * time.Second
	}
	if service.maxPreviewBytes <= 0 {
		service.maxPreviewBytes = 1024 * 1024
	}
	if service.maxRequestBodyBytes <= 0 {
		service.maxRequestBodyBytes = 256 * 1024
	}
	if service.maxConcurrentPerUser <= 0 {
		service.maxConcurrentPerUser = 5
	}
	if service.maxConcurrentPerIP <= 0 {
		service.maxConcurrentPerIP = 1
	}
	if client.Timeout <= 0 {
		client.Timeout = service.requestTimeout
	}

	return service
}

func (s *Service) Execute(ctx context.Context, userID string, clientIP string, input RunInput) (RunResult, error) {
	if s == nil || s.client == nil || s.history == nil || s.limiter == nil {
		return RunResult{}, ErrUnavailable
	}

	userID = strings.TrimSpace(userID)
	clientIP = strings.TrimSpace(clientIP)
	if userID == "" {
		return RunResult{}, ErrInvalid
	}

	request, rawURL, requestBody, err := s.buildRequest(ctx, input)
	if err != nil {
		now := time.Now().UTC()
		switch {
		case errors.Is(err, ErrRequestTooLarge):
			_ = s.persistFailure(ctx, userID, input, requestBody, now, rawURL, "blocked", "request_body_too_large", "request_body_too_large", "request body exceeds the authenticated limit")
		case isValidationError(err):
			_ = s.persistFailure(ctx, userID, input, mustJSON(input.Body), now, rawURL, "blocked", validationErrorCode(err), validationErrorCode(err), validationErrorMessage(err))
		}
		return RunResult{}, err
	}

	domainKey, err := ratelimit.DomainKeyFromURL(rawURL)
	if err != nil {
		return RunResult{}, err
	}

	startedAt := time.Now().UTC()
	release, limitErr, err := s.acquireLimits(ctx, userID, clientIP, domainKey, input, rawURL, requestBody, startedAt)
	if err != nil {
		return RunResult{}, err
	}
	if limitErr != nil {
		return RunResult{}, limitErr
	}
	if release != nil {
		defer release()
	}

	requestCtx, cancel := context.WithTimeout(ctx, s.requestTimeout)
	defer cancel()
	request = request.WithContext(requestCtx)

	redirects := 0
	client := *s.client
	client.Timeout = s.requestTimeout
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirects = len(via)
		chain := make([]string, 0, len(via)+2)
		chain = append(chain, rawURL)
		for _, item := range via {
			chain = append(chain, item.URL.String())
		}
		chain = append(chain, req.URL.String())
		_, err := safety.ValidateRedirectChain(req.Context(), chain, s.safetyOptions)
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		if isTimeoutError(err, requestCtx) {
			completedAt := time.Now().UTC()
			_ = s.persistFailure(ctx, userID, input, requestBody, startedAt, rawURL, "failed", "timeout", "request_timeout", "request timed out")
			_ = s.recordUsage(ctx, userID, nil, input, rawURL, "timeout", "request_timeout", 0, startedAt, completedAt)
			return RunResult{}, ErrTimedOut
		}

		_ = s.persistFailure(ctx, userID, input, requestBody, startedAt, rawURL, "failed", "", "upstream_request_failed", err.Error())
		_ = s.recordUsage(ctx, userID, nil, input, rawURL, "failed", "upstream_request_failed", 0, startedAt, time.Now().UTC())
		return RunResult{}, err
	}
	defer response.Body.Close()

	bodyPreview, bodyJSON, sizeBytes, truncated, readErr := readPreview(response.Body, s.maxPreviewBytes)
	if readErr != nil {
		if isTimeoutError(readErr, requestCtx) {
			completedAt := time.Now().UTC()
			_ = s.persistFailure(ctx, userID, input, requestBody, startedAt, rawURL, "failed", "timeout", "request_timeout", "request timed out")
			_ = s.recordUsage(ctx, userID, nil, input, rawURL, "timeout", "request_timeout", 0, startedAt, completedAt)
			return RunResult{}, ErrTimedOut
		}

		return RunResult{}, readErr
	}

	completedAt := time.Now().UTC()
	finalURL := response.Request.URL.String()
	historyRecord, err := s.history.Create(ctx, history.CreateParams{
		UserID:              userID,
		Source:              "authenticated",
		Status:              "succeeded",
		Method:              request.Method,
		URL:                 rawURL,
		FinalURL:            &finalURL,
		TargetHost:          response.Request.URL.Hostname(),
		RequestHeaders:      mustJSON(request.Header),
		RequestQueryParams:  mustJSON(input.QueryParams),
		RequestAuth:         mustJSON(input.Auth),
		RequestBody:         requestBody,
		ResponseStatus:      intPtr(response.StatusCode),
		ResponseHeaders:     mustJSON(response.Header),
		ResponseBodyPreview: bodyPreview,
		ResponseSizeBytes:   intPtr(sizeBytes),
		ResponseTimeMS:      intPtr(int(completedAt.Sub(startedAt).Milliseconds())),
		ResponseContentType: response.Header.Get("Content-Type"),
		RedirectCount:       redirects,
		StartedAt:           &startedAt,
		CompletedAt:         &completedAt,
		Metadata:            mustJSON(map[string]any{"truncated": truncated}),
	})
	if err != nil {
		return RunResult{}, err
	}

	_ = s.recordUsage(ctx, userID, &historyRecord.ID, input, rawURL, "succeeded", "", response.StatusCode, startedAt, completedAt)

	return RunResult{
		RunID:               historyRecord.ID,
		Status:              "succeeded",
		Method:              request.Method,
		URL:                 rawURL,
		FinalURL:            finalURL,
		ResponseStatus:      response.StatusCode,
		ResponseHeaders:     response.Header,
		ResponseBody:        bodyPreview,
		ResponseJSON:        bodyJSON,
		ResponseSizeBytes:   sizeBytes,
		ResponseTimeMS:      int(completedAt.Sub(startedAt).Milliseconds()),
		ResponseContentType: response.Header.Get("Content-Type"),
		RedirectCount:       redirects,
		Truncated:           truncated,
	}, nil
}

func (s *Service) buildRequest(ctx context.Context, input RunInput) (*http.Request, string, json.RawMessage, error) {
	method := strings.ToUpper(strings.TrimSpace(input.Method))
	rawURL := strings.TrimSpace(input.URL)
	if method == "" || rawURL == "" {
		return nil, "", nil, ErrInvalid
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, "", nil, err
	}
	values := parsed.Query()
	for _, item := range input.QueryParams {
		if !item.Enabled {
			continue
		}
		values.Set(item.Name, item.Value)
	}
	parsed.RawQuery = values.Encode()
	rawURL = parsed.String()

	if _, err := safety.ValidateURL(ctx, rawURL, s.safetyOptions); err != nil {
		blockCode := "blocked_target"
		blockMessage := err.Error()
		var validationErr *safety.ValidationError
		if errors.As(err, &validationErr) {
			blockCode = string(validationErr.Code)
			blockMessage = validationErr.Message
		}
		return nil, rawURL, nil, &safety.ValidationError{Code: safety.ErrorCode(blockCode), Message: blockMessage, URL: rawURL}
	}

	bodyReader, requestBody, err := s.encodeBody(input.Body)
	if err != nil {
		return nil, rawURL, nil, err
	}

	request, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		return nil, rawURL, nil, err
	}
	for _, item := range input.Headers {
		if !item.Enabled {
			continue
		}
		request.Header.Set(item.Name, item.Value)
	}
	applyAuth(request, input.Auth)
	if input.Body.Mode == "json" && request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}
	if input.Body.Mode == "form_urlencoded" && request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return request, rawURL, requestBody, nil
}

func (s *Service) encodeBody(body BodyInput) (io.Reader, json.RawMessage, error) {
	switch strings.TrimSpace(body.Mode) {
	case "", "none":
		return nil, json.RawMessage(`{}`), nil
	case "raw", "json":
		if len(body.Raw) > s.maxRequestBodyBytes {
			return nil, nil, ErrRequestTooLarge
		}
		payload := json.RawMessage(`{"mode":"` + body.Mode + `","raw":` + strconvQuote(body.Raw) + `}`)
		return strings.NewReader(body.Raw), payload, nil
	case "form_urlencoded":
		values := url.Values{}
		for _, field := range body.FormFields {
			if !field.Enabled {
				continue
			}
			values.Set(field.Name, field.Value)
		}
		encoded := values.Encode()
		if len(encoded) > s.maxRequestBodyBytes {
			return nil, nil, ErrRequestTooLarge
		}
		payload := mustJSON(body)
		return strings.NewReader(encoded), payload, nil
	default:
		return nil, nil, ErrInvalid
	}
}

func applyAuth(request *http.Request, authInput AuthInput) {
	switch strings.TrimSpace(authInput.Scheme) {
	case "basic":
		credentials := authInput.Username + ":" + authInput.Password
		request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(credentials)))
	case "bearer":
		request.Header.Set("Authorization", "Bearer "+authInput.Token)
	}
}

func (s *Service) acquireLimits(ctx context.Context, userID, clientIP, domainKey string, input RunInput, rawURL string, requestBody json.RawMessage, startedAt time.Time) (func(), *LimitError, error) {
	method := strings.ToUpper(strings.TrimSpace(input.Method))
	scopeChecks := []struct {
		scope ratelimit.Scope
		key   string
		allow func(string) (ratelimit.Decision, error)
	}{
		{scope: ratelimit.ScopeIP, key: clientIP, allow: s.limiter.AllowIP},
		{scope: ratelimit.ScopeUser, key: userID, allow: s.limiter.AllowUser},
		{scope: ratelimit.ScopeDomain, key: domainKey, allow: s.limiter.AllowDomain},
	}

	for _, check := range scopeChecks {
		if strings.TrimSpace(check.key) == "" {
			continue
		}
		decision, err := check.allow(check.key)
		if err != nil {
			return nil, nil, err
		}
		if !decision.Allowed {
			completedAt := time.Now().UTC()
			runID := s.persistBlockedRun(ctx, userID, input, requestBody, rawURL, decision, startedAt, completedAt)
			_ = s.recordUsage(ctx, userID, runID, input, rawURL, "limited", "rate_limited", http.StatusTooManyRequests, startedAt, completedAt)
			_ = s.recordAbuse(ctx, userID, clientIP, runID, input, rawURL, decision, method, requestBody, completedAt)
			return nil, &LimitError{
				Scope:      decision.Scope,
				Key:        decision.Key,
				Reason:     decision.Reason,
				RetryAfter: decision.RetryAfter,
				Message:    limitMessage(decision),
			}, nil
		}
	}

	release, err := s.acquireConcurrency(userID, clientIP)
	if err != nil {
		if errors.Is(err, ErrConcurrencyLimit) {
			decision := ratelimit.Decision{
				Allowed:    false,
				Scope:      ratelimit.ScopeUser,
				Key:        userID,
				Reason:     "concurrency_limit",
				RetryAfter: 0,
			}
			completedAt := time.Now().UTC()
			runID := s.persistBlockedRun(ctx, userID, input, requestBody, rawURL, decision, startedAt, completedAt)
			_ = s.recordUsage(ctx, userID, runID, input, rawURL, "limited", "concurrency_limit", http.StatusTooManyRequests, startedAt, completedAt)
			_ = s.recordAbuse(ctx, userID, clientIP, runID, input, rawURL, decision, method, requestBody, completedAt)
			return nil, &LimitError{
				Scope:      ratelimit.ScopeUser,
				Key:        userID,
				Reason:     "concurrency_limit",
				RetryAfter: 0,
				Message:    "too many active authenticated requests",
			}, nil
		}
		return nil, nil, err
	}

	return release, nil, nil
}

func (s *Service) acquireConcurrency(userID, clientIP string) (func(), error) {
	userID = strings.TrimSpace(userID)
	clientIP = strings.TrimSpace(clientIP)
	if userID == "" || clientIP == "" {
		return func() {}, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeByUser[userID] >= s.maxConcurrentPerUser || s.activeByIP[clientIP] >= s.maxConcurrentPerIP {
		return nil, ErrConcurrencyLimit
	}

	s.activeByUser[userID]++
	s.activeByIP[clientIP]++

	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		if s.activeByUser[userID] > 0 {
			s.activeByUser[userID]--
		}
		if s.activeByIP[clientIP] > 0 {
			s.activeByIP[clientIP]--
		}
	}, nil
}

func (s *Service) persistFailure(ctx context.Context, userID string, input RunInput, requestBody json.RawMessage, startedAt time.Time, rawURL, status, code, errorCode, errorMessage string) error {
	_, err := s.history.Create(ctx, history.CreateParams{
		UserID:             userID,
		Source:             "authenticated",
		Status:             status,
		Method:             strings.ToUpper(strings.TrimSpace(input.Method)),
		URL:                firstNonEmpty(rawURL, strings.TrimSpace(input.URL)),
		RequestHeaders:     mustJSON(input.Headers),
		RequestQueryParams: mustJSON(input.QueryParams),
		RequestAuth:        mustJSON(input.Auth),
		RequestBody:        requestBody,
		BlockedReason:      code,
		ErrorCode:          errorCode,
		ErrorMessage:       errorMessage,
		StartedAt:          &startedAt,
		CompletedAt:        timePtr(time.Now().UTC()),
	})
	return err
}

func (s *Service) persistBlockedRun(ctx context.Context, userID string, input RunInput, requestBody json.RawMessage, rawURL string, decision ratelimit.Decision, startedAt, completedAt time.Time) *string {
	record, err := s.history.Create(ctx, history.CreateParams{
		UserID:             userID,
		Source:             "authenticated",
		Status:             "blocked",
		Method:             strings.ToUpper(strings.TrimSpace(input.Method)),
		URL:                firstNonEmpty(rawURL, strings.TrimSpace(input.URL)),
		RequestHeaders:     mustJSON(input.Headers),
		RequestQueryParams: mustJSON(input.QueryParams),
		RequestAuth:        mustJSON(input.Auth),
		RequestBody:        requestBody,
		BlockedReason:      decision.Reason,
		ErrorCode:          "rate_limited",
		ErrorMessage:       limitMessage(decision),
		StartedAt:          &startedAt,
		CompletedAt:        &completedAt,
		Metadata:           mustJSON(map[string]any{"scope": decision.Scope, "key": decision.Key, "retryAfterMs": int(decision.RetryAfter.Milliseconds())}),
	})
	if err != nil {
		return nil
	}

	return &record.ID
}

func (s *Service) recordUsage(ctx context.Context, userID string, runID *string, input RunInput, rawURL, outcome, errorCode string, status int, startedAt, completedAt time.Time) error {
	if s.usageRecorder == nil {
		return nil
	}

	event := usage.Event{
		UserID:       stringPtr(userID),
		RequestRunID: runID,
		Bucket:       "hour",
		EventKey:     "authenticated.run." + outcome,
		Quantity:     1,
		Dimensions: mustJSON(map[string]any{
			"method":     strings.ToUpper(strings.TrimSpace(input.Method)),
			"url":        firstNonEmpty(rawURL, strings.TrimSpace(input.URL)),
			"outcome":    outcome,
			"errorCode":  errorCode,
			"status":     status,
			"durationMs": int(completedAt.Sub(startedAt).Milliseconds()),
			"bodyBytes":  len(input.Body.Raw),
			"clientType": "authenticated",
		}),
		OccurredAt: completedAt,
	}

	_, err := s.usageRecorder.Create(ctx, event)
	return err
}

func (s *Service) recordAbuse(ctx context.Context, userID, clientIP string, runID *string, input RunInput, rawURL string, decision ratelimit.Decision, method string, requestBody json.RawMessage, createdAt time.Time) error {
	if s.abuseRecorder == nil {
		return nil
	}

	target := firstNonEmpty(rawURL, strings.TrimSpace(input.URL))
	event := abuse.Event{
		UserID:      stringPtr(userID),
		RequestID:   runID,
		SourceIP:    stringPtr(clientIP),
		Target:      stringPtr(target),
		RuleKey:     "authenticated-rate-limit",
		Category:    abuse.CategorySuspicious,
		Severity:    abuse.SeverityMedium,
		ActionTaken: abuse.ActionBlocked,
		Message:     limitMessage(decision),
		Details: mustJSON(map[string]any{
			"scope":        decision.Scope,
			"key":          decision.Key,
			"reason":       decision.Reason,
			"retryAfterMs": int(decision.RetryAfter.Milliseconds()),
			"method":       method,
			"url":          target,
			"clientIp":     clientIP,
			"bodyBytes":    len(requestBody),
		}),
		CreatedAt: createdAt,
	}

	_, err := s.abuseRecorder.Create(ctx, event)
	return err
}

func limitMessage(decision ratelimit.Decision) string {
	switch decision.Reason {
	case "cooldown":
		return "request rate limit cooldown in effect"
	case "burst_limit":
		return "request burst limit exceeded"
	case "quota_limit":
		return "request quota limit exceeded"
	case "concurrency_limit":
		return "too many active authenticated requests"
	default:
		return "authenticated request limit exceeded"
	}
}

func isValidationError(err error) bool {
	var validationErr *safety.ValidationError
	return errors.As(err, &validationErr)
}

func validationErrorCode(err error) string {
	var validationErr *safety.ValidationError
	if errors.As(err, &validationErr) {
		return string(validationErr.Code)
	}
	return "blocked_target"
}

func validationErrorMessage(err error) string {
	var validationErr *safety.ValidationError
	if errors.As(err, &validationErr) {
		return validationErr.Message
	}
	return err.Error()
}

func isTimeoutError(err error, ctx context.Context) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrTimedOut) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	if ctx != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return true
	}
	return false
}

func readPreview(body io.Reader, maxBytes int) (string, any, int, bool, error) {
	limited := io.LimitReader(body, int64(maxBytes+1))
	payload, err := io.ReadAll(limited)
	if err != nil {
		return "", nil, 0, false, err
	}
	truncated := len(payload) > maxBytes
	if truncated {
		payload = payload[:maxBytes]
	}
	text := string(payload)
	var parsed any
	if json.Unmarshal(payload, &parsed) != nil {
		parsed = nil
	}
	return text, parsed, len(payload), truncated, nil
}

func mustJSON(value any) json.RawMessage {
	payload, _ := json.Marshal(value)
	if len(payload) == 0 {
		return json.RawMessage(`{}`)
	}
	return payload
}

func strconvQuote(value string) string {
	payload, _ := json.Marshal(value)
	return string(payload)
}

func intPtr(value int) *int {
	return &value
}

func timePtr(value time.Time) *time.Time {
	return &value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func stringPtr(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
