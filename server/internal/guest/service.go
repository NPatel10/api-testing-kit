package guest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"api-testing-kit/server/internal/abuse"
	"api-testing-kit/server/internal/ratelimit"
	"api-testing-kit/server/internal/safety"
	"api-testing-kit/server/internal/templates"
	"api-testing-kit/server/internal/usage"
)

type UsageRecorder interface {
	Create(ctx context.Context, event usage.Event) (usage.Event, error)
}

type AbuseRecorder interface {
	Create(ctx context.Context, event abuse.Event) (abuse.Event, error)
}

type GuestErrorCode string

const (
	ErrorUnavailable           GuestErrorCode = "guest_unavailable"
	ErrorInvalidInput          GuestErrorCode = "invalid_guest_run"
	ErrorTemplateNotFound      GuestErrorCode = "guest_template_not_found"
	ErrorTemplateNotAllowed    GuestErrorCode = "guest_template_not_allowed"
	ErrorOverrideNotAllowed    GuestErrorCode = "guest_override_not_allowed"
	ErrorRateLimited           GuestErrorCode = "guest_rate_limited"
	ErrorRequestTooLarge       GuestErrorCode = "guest_request_too_large"
	ErrorBlockedTarget         GuestErrorCode = "blocked_target"
	ErrorUpstreamRequestFailed GuestErrorCode = "guest_request_failed"
)

type GuestError struct {
	Code    GuestErrorCode
	Message string
	Status  int
	Err     error
}

func (e *GuestError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *GuestError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

type RunRequest struct {
	TemplateSlug string        `json:"templateSlug,omitempty"`
	TemplateID   string        `json:"templateId,omitempty"`
	Overrides    *RunOverrides `json:"overrides,omitempty"`
}

type RunOverrides struct {
	QueryParams map[string]string `json:"queryParams,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Body        *BodyOverride     `json:"body,omitempty"`
}

type BodyOverride struct {
	Raw        string            `json:"raw,omitempty"`
	FormFields map[string]string `json:"formFields,omitempty"`
}

type RunResult struct {
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

type Service struct {
	client               *http.Client
	limiter              *ratelimit.Limiter
	safetyOptions        safety.Options
	usageRecorder        UsageRecorder
	abuseRecorder        AbuseRecorder
	now                  func() time.Time
	maxRequestBodyBytes  int
	maxResponseBodyBytes int
	requestTimeout       time.Duration
}

func NewService(client *http.Client, limiter *ratelimit.Limiter, usageRecorder UsageRecorder, abuseRecorder AbuseRecorder, opts safety.Options) *Service {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	if limiter == nil {
		limiter = ratelimit.NewLimiter(ratelimit.GuestConfig())
	}

	normalized := opts
	if len(normalized.AllowedSchemes) == 0 {
		normalized.AllowedSchemes = safety.DefaultOptions().AllowedSchemes
	}
	if len(normalized.AllowedPorts) == 0 {
		normalized.AllowedPorts = safety.DefaultOptions().AllowedPorts
	}
	if normalized.MaxRedirects == 0 {
		normalized.MaxRedirects = safety.DefaultOptions().MaxRedirects
	}
	if normalized.Resolver == nil {
		normalized.Resolver = safety.DefaultOptions().Resolver
	}

	return &Service{
		client:               client,
		limiter:              limiter,
		safetyOptions:        normalized,
		usageRecorder:        usageRecorder,
		abuseRecorder:        abuseRecorder,
		now:                  time.Now,
		maxRequestBodyBytes:  64 * 1024,
		maxResponseBodyBytes: 512 * 1024,
		requestTimeout:       10 * time.Second,
	}
}

func (s *Service) Execute(ctx context.Context, clientIP string, input RunRequest) (RunResult, error) {
	if s == nil || s.client == nil || s.limiter == nil {
		return RunResult{}, &GuestError{
			Code:    ErrorUnavailable,
			Message: "guest execution is temporarily unavailable",
			Status:  http.StatusServiceUnavailable,
		}
	}

	template, err := resolveTemplate(input)
	if err != nil {
		return RunResult{}, err
	}
	if !template.GuestSafe {
		return RunResult{}, &GuestError{
			Code:    ErrorTemplateNotAllowed,
			Message: "selected template is not available in guest mode",
			Status:  http.StatusForbidden,
		}
	}

	plan, err := buildRequestPlan(template, input.Overrides)
	if err != nil {
		return RunResult{}, err
	}

	if _, err := safety.ValidateURL(ctx, plan.URL, s.safetyOptions); err != nil {
		return RunResult{}, s.blockedTarget(ctx, template.Slug, plan.URL, err)
	}

	if err := s.checkRateLimit(ctx, clientIP, plan.URL, template.Slug); err != nil {
		return RunResult{}, err
	}

	if len(plan.Body) > s.maxRequestBodyBytes {
		_ = s.recordAbuse(ctx, template.Slug, plan.URL, abuse.CategorySuspicious, "body_size_limit", abuse.ActionBlocked, "guest request body exceeded limit")
		return RunResult{}, &GuestError{
			Code:    ErrorRequestTooLarge,
			Message: "guest request body is too large",
			Status:  http.StatusRequestEntityTooLarge,
		}
	}

	reqCtx, cancel := context.WithTimeout(ctx, s.requestTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(reqCtx, plan.Method, plan.URL, strings.NewReader(plan.Body))
	if err != nil {
		return RunResult{}, &GuestError{
			Code:    ErrorInvalidInput,
			Message: "template request could not be constructed",
			Status:  http.StatusBadRequest,
			Err:     err,
		}
	}
	for _, header := range plan.Headers {
		request.Header.Set(header.Name, header.Value)
	}
	if request.Header.Get("Content-Type") == "" && plan.ContentType != "" {
		request.Header.Set("Content-Type", plan.ContentType)
	}

	startedAt := s.now().UTC()
	response, redirectCount, err := s.doRequest(request)
	if err != nil {
		_ = s.recordUsage(ctx, template.Slug, clientIP, plan.URL, "failed", 0)
		return RunResult{}, s.upstreamFailure(ctx, template.Slug, plan.URL, err)
	}
	defer response.Body.Close()

	preview, parsedJSON, sizeBytes, truncated, err := readPreview(response.Body, s.maxResponseBodyBytes)
	if err != nil {
		_ = s.recordUsage(ctx, template.Slug, clientIP, plan.URL, "failed", 0)
		return RunResult{}, s.upstreamFailure(ctx, template.Slug, plan.URL, err)
	}

	completedAt := s.now().UTC()
	_ = s.recordUsage(ctx, template.Slug, clientIP, plan.URL, "succeeded", response.StatusCode)
	finalURL := plan.URL
	if response.Request != nil && response.Request.URL != nil {
		finalURL = response.Request.URL.String()
	}

	return RunResult{
		Status:              "succeeded",
		Method:              request.Method,
		URL:                 plan.URL,
		FinalURL:            finalURL,
		ResponseStatus:      response.StatusCode,
		ResponseHeaders:     response.Header,
		ResponseBody:        preview,
		ResponseJSON:        parsedJSON,
		ResponseSizeBytes:   sizeBytes,
		ResponseTimeMS:      int(completedAt.Sub(startedAt).Milliseconds()),
		ResponseContentType: response.Header.Get("Content-Type"),
		RedirectCount:       redirectCount,
		Truncated:           truncated,
	}, nil
}

type requestPlan struct {
	Method      string
	URL         string
	Headers     []headerPair
	Body        string
	ContentType string
}

type headerPair struct {
	Name  string
	Value string
}

func resolveTemplate(input RunRequest) (templates.Template, error) {
	templateRef := strings.TrimSpace(input.TemplateSlug)
	if templateRef == "" {
		templateRef = strings.TrimSpace(input.TemplateID)
	}
	if templateRef == "" {
		return templates.Template{}, &GuestError{
			Code:    ErrorInvalidInput,
			Message: "template reference is required",
			Status:  http.StatusBadRequest,
		}
	}

	template, ok := templates.Get(templateRef)
	if !ok {
		return templates.Template{}, &GuestError{
			Code:    ErrorTemplateNotFound,
			Message: "template not found",
			Status:  http.StatusNotFound,
		}
	}

	return template, nil
}

func buildRequestPlan(template templates.Template, overrides *RunOverrides) (requestPlan, error) {
	allowedQuery, allowedHeaders, bodyAllowed := parseAllowedOverrides(template.AllowedOverrides)

	queryOverrides := map[string]string{}
	headerOverrides := map[string]string{}
	var bodyOverride *BodyOverride
	if overrides != nil {
		for key, value := range overrides.QueryParams {
			key = normalizeKey(key)
			if key != "" {
				queryOverrides[key] = value
			}
		}
		for key, value := range overrides.Headers {
			key = normalizeKey(key)
			if key != "" {
				headerOverrides[key] = value
			}
		}
		bodyOverride = overrides.Body
	}

	for key := range queryOverrides {
		if _, ok := allowedQuery[key]; !ok {
			return requestPlan{}, &GuestError{
				Code:    ErrorOverrideNotAllowed,
				Message: fmt.Sprintf("query override %q is not allowed", key),
				Status:  http.StatusForbidden,
			}
		}
	}
	for key := range headerOverrides {
		if _, ok := allowedHeaders[key]; !ok {
			return requestPlan{}, &GuestError{
				Code:    ErrorOverrideNotAllowed,
				Message: fmt.Sprintf("header override %q is not allowed", key),
				Status:  http.StatusForbidden,
			}
		}
	}

	baseURL, err := url.Parse(template.Request.URL)
	if err != nil {
		return requestPlan{}, &GuestError{
			Code:    ErrorInvalidInput,
			Message: "template URL is invalid",
			Status:  http.StatusBadRequest,
			Err:     err,
		}
	}

	queryValues := baseURL.Query()
	for _, item := range template.Request.QueryParams {
		key := normalizeKey(item.Name)
		if override, ok := queryOverrides[key]; ok {
			item.Value = override
		}
		if item.Enabled {
			queryValues.Set(item.Name, item.Value)
		}
	}
	for key, value := range queryOverrides {
		if !hasQueryParam(template.Request.QueryParams, key) {
			queryValues.Set(key, value)
		}
	}
	baseURL.RawQuery = queryValues.Encode()

	headers := make([]headerPair, 0, len(template.Request.Headers)+len(headerOverrides))
	for _, item := range template.Request.Headers {
		if item.Enabled {
			headers = append(headers, headerPair{Name: item.Name, Value: item.Value})
		}
	}
	for key, value := range headerOverrides {
		replaced := false
		for i := range headers {
			if strings.EqualFold(headers[i].Name, key) {
				headers[i].Value = value
				replaced = true
				break
			}
		}
		if !replaced {
			headers = append(headers, headerPair{Name: key, Value: value})
		}
	}

	bodyMode := strings.ToLower(strings.TrimSpace(template.Request.Body.Mode))
	var body string
	var contentType string
	if bodyOverride != nil && !bodyAllowed {
		return requestPlan{}, &GuestError{
			Code:    ErrorOverrideNotAllowed,
			Message: "body overrides are not allowed for this template",
			Status:  http.StatusForbidden,
		}
	}

	switch bodyMode {
	case "", "none":
		if bodyOverride != nil && (strings.TrimSpace(bodyOverride.Raw) != "" || len(bodyOverride.FormFields) > 0) {
			return requestPlan{}, &GuestError{
				Code:    ErrorOverrideNotAllowed,
				Message: "body overrides are not allowed for this template",
				Status:  http.StatusForbidden,
			}
		}
	case "raw", "json":
		body = template.Request.Body.Raw
		if bodyOverride != nil {
			if len(bodyOverride.FormFields) > 0 {
				return requestPlan{}, &GuestError{
					Code:    ErrorOverrideNotAllowed,
					Message: "form field overrides are only supported for form templates",
					Status:  http.StatusForbidden,
				}
			}
			if strings.TrimSpace(bodyOverride.Raw) != "" {
				body = bodyOverride.Raw
			}
		}
		contentType = template.Request.Body.ContentType
		if contentType == "" && bodyMode == "json" {
			contentType = "application/json"
		}
	case "form_urlencoded":
		fields := template.Request.Body.FormFields
		if bodyOverride != nil {
			if strings.TrimSpace(bodyOverride.Raw) != "" {
				return requestPlan{}, &GuestError{
					Code:    ErrorOverrideNotAllowed,
					Message: "raw body overrides are not supported for form templates",
					Status:  http.StatusForbidden,
				}
			}
			if len(bodyOverride.FormFields) > 0 {
				fields = fields[:0]
				for key, value := range bodyOverride.FormFields {
					fields = append(fields, templates.Param{Name: key, Value: value, Enabled: true})
				}
			}
		}
		values := url.Values{}
		for _, field := range fields {
			if field.Enabled {
				values.Set(field.Name, field.Value)
			}
		}
		body = values.Encode()
		contentType = "application/x-www-form-urlencoded"
	default:
		return requestPlan{}, &GuestError{
			Code:    ErrorInvalidInput,
			Message: "unsupported template body mode",
			Status:  http.StatusBadRequest,
		}
	}

	return requestPlan{
		Method:      template.Request.Method,
		URL:         baseURL.String(),
		Headers:     headers,
		Body:        body,
		ContentType: contentType,
	}, nil
}

func parseAllowedOverrides(values []string) (map[string]struct{}, map[string]struct{}, bool) {
	allowedQuery := make(map[string]struct{})
	allowedHeaders := make(map[string]struct{})
	bodyAllowed := false

	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		switch {
		case strings.HasPrefix(normalized, "queryparams."):
			allowedQuery[strings.TrimPrefix(normalized, "queryparams.")] = struct{}{}
		case strings.HasPrefix(normalized, "headers."):
			allowedHeaders[strings.TrimPrefix(normalized, "headers.")] = struct{}{}
		case strings.HasPrefix(normalized, "body."):
			bodyAllowed = true
		}
	}

	return allowedQuery, allowedHeaders, bodyAllowed
}

func (s *Service) checkRateLimit(ctx context.Context, clientIP, finalURL, templateSlug string) error {
	if clientIP != "" {
		decision, err := s.limiter.AllowIP(clientIP)
		if err != nil {
			return err
		}
		if !decision.Allowed {
			_ = s.recordAbuse(ctx, templateSlug, finalURL, abuse.CategorySuspicious, "ip_rate_limit", abuse.ActionBlocked, decision.Reason)
			return &GuestError{
				Code:    ErrorRateLimited,
				Message: "guest request rate limit exceeded",
				Status:  http.StatusTooManyRequests,
			}
		}
	}

	domainKey, err := ratelimit.DomainKeyFromURL(finalURL)
	if err != nil {
		return &GuestError{
			Code:    ErrorInvalidInput,
			Message: "template URL is invalid",
			Status:  http.StatusBadRequest,
			Err:     err,
		}
	}
	decision, err := s.limiter.AllowDomain(domainKey)
	if err != nil {
		return err
	}
	if !decision.Allowed {
		_ = s.recordAbuse(ctx, templateSlug, finalURL, abuse.CategorySuspicious, "domain_rate_limit", abuse.ActionBlocked, decision.Reason)
		return &GuestError{
			Code:    ErrorRateLimited,
			Message: "guest request rate limit exceeded",
			Status:  http.StatusTooManyRequests,
		}
	}

	return nil
}

func (s *Service) doRequest(request *http.Request) (*http.Response, int, error) {
	redirects := 0
	client := *s.client
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirects = len(via)
		var urls []string
		urls = append(urls, request.URL.String())
		for _, item := range via {
			urls = append(urls, item.URL.String())
		}
		urls = append(urls, req.URL.String())
		_, err := safety.ValidateRedirectChain(req.Context(), urls, s.safetyOptions)
		return err
	}

	response, err := client.Do(request)
	return response, redirects, err
}

func (s *Service) blockedTarget(ctx context.Context, templateSlug, finalURL string, err error) error {
	_ = s.recordAbuse(ctx, templateSlug, finalURL, abuse.CategoryBlockedHost, "blocked_target", abuse.ActionBlocked, err.Error())
	return &GuestError{
		Code:    ErrorBlockedTarget,
		Message: err.Error(),
		Status:  http.StatusForbidden,
		Err:     err,
	}
}

func (s *Service) upstreamFailure(ctx context.Context, templateSlug, finalURL string, err error) error {
	_ = s.recordAbuse(ctx, templateSlug, finalURL, abuse.CategorySuspicious, "upstream_failure", abuse.ActionBlocked, err.Error())
	return &GuestError{
		Code:    ErrorUpstreamRequestFailed,
		Message: err.Error(),
		Status:  http.StatusBadGateway,
		Err:     err,
	}
}

func (s *Service) recordUsage(ctx context.Context, templateSlug, clientIP, targetURL, outcome string, status int) error {
	if s.usageRecorder == nil {
		return nil
	}

	event := usage.Event{
		Bucket:   "hour",
		EventKey: "guest.run." + outcome,
		Quantity: 1,
		Dimensions: mustJSON(map[string]any{
			"templateSlug": templateSlug,
			"clientIp":     clientIP,
			"target":       targetURL,
			"status":       status,
			"outcome":      outcome,
		}),
		OccurredAt: s.now().UTC(),
	}
	_, err := s.usageRecorder.Create(ctx, event)
	return err
}

func (s *Service) recordAbuse(ctx context.Context, templateSlug, targetURL string, category abuse.Category, ruleKey string, action abuse.Action, message string) error {
	if s.abuseRecorder == nil {
		return nil
	}

	target := targetURL
	event := abuse.Event{
		Severity:    abuse.SeverityMedium,
		Category:    category,
		Target:      &target,
		RuleKey:     ruleKey,
		ActionTaken: action,
		Message:     message,
		Details: mustJSON(map[string]any{
			"templateSlug": templateSlug,
			"target":       targetURL,
			"ruleKey":      ruleKey,
			"message":      message,
		}),
		CreatedAt: s.now().UTC(),
	}
	_, err := s.abuseRecorder.Create(ctx, event)
	return err
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

func normalizeKey(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func hasQueryParam(items []templates.Param, key string) bool {
	for _, item := range items {
		if strings.EqualFold(strings.TrimSpace(item.Name), key) {
			return true
		}
	}
	return false
}
