package runner

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"api-testing-kit/server/internal/history"
	"api-testing-kit/server/internal/safety"
)

var (
	ErrUnavailable = errors.New("runner is unavailable")
	ErrInvalid     = errors.New("invalid run payload")
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

type Service struct {
	client          *http.Client
	history         *history.Service
	safetyOptions   safety.Options
	maxPreviewBytes int
}

func NewService(client *http.Client, historyService *history.Service, opts safety.Options) *Service {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &Service{
		client:          client,
		history:         historyService,
		safetyOptions:   opts,
		maxPreviewBytes: 64 * 1024,
	}
}

func (s *Service) Execute(ctx context.Context, userID string, input RunInput) (RunResult, error) {
	if s == nil || s.client == nil || s.history == nil {
		return RunResult{}, ErrUnavailable
	}

	request, rawURL, requestBody, err := s.buildRequest(ctx, input)
	if err != nil {
		var validationErr *safety.ValidationError
		if errors.As(err, &validationErr) {
			_ = s.persistFailure(ctx, userID, input, mustJSON(input.Body), time.Now().UTC(), rawURL, "blocked", string(validationErr.Code), string(validationErr.Code), validationErr.Message)
		}
		return RunResult{}, err
	}

	startedAt := time.Now().UTC()
	redirects := []string{rawURL}
	client := *s.client
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirects = redirects[:0]
		for _, item := range via {
			redirects = append(redirects, item.URL.String())
		}
		redirects = append(redirects, req.URL.String())
		_, err := safety.ValidateRedirectChain(req.Context(), redirects, s.safetyOptions)
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		runErr := s.persistFailure(ctx, userID, input, requestBody, startedAt, rawURL, "failed", "", "upstream_request_failed", err.Error())
		if runErr != nil {
			return RunResult{}, runErr
		}
		return RunResult{}, err
	}
	defer response.Body.Close()

	bodyPreview, bodyJSON, sizeBytes, truncated, readErr := readPreview(response.Body, s.maxPreviewBytes)
	if readErr != nil {
		return RunResult{}, readErr
	}

	finalURL := response.Request.URL.String()
	completedAt := time.Now().UTC()
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
		RedirectCount:       len(redirects) - 1,
		StartedAt:           &startedAt,
		CompletedAt:         &completedAt,
		Metadata:            mustJSON(map[string]any{"truncated": truncated}),
	})
	if err != nil {
		return RunResult{}, err
	}

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
		RedirectCount:       len(redirects) - 1,
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

	bodyReader, requestBody, err := encodeBody(input.Body)
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

func encodeBody(body BodyInput) (io.Reader, json.RawMessage, error) {
	switch strings.TrimSpace(body.Mode) {
	case "", "none":
		return nil, json.RawMessage(`{}`), nil
	case "raw", "json":
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
		payload := mustJSON(body)
		return strings.NewReader(values.Encode()), payload, nil
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
