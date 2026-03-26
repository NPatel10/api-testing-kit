package history

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrUnavailable = errors.New("history repository is unavailable")
	ErrInvalid     = errors.New("invalid history input")
)

type RunRecord struct {
	ID                  string          `json:"id"`
	UserID              *string         `json:"userId,omitempty"`
	CollectionID        *string         `json:"collectionId,omitempty"`
	SavedRequestID      *string         `json:"savedRequestId,omitempty"`
	Source              string          `json:"source"`
	Status              string          `json:"status"`
	Method              string          `json:"method"`
	URL                 string          `json:"url"`
	FinalURL            *string         `json:"finalUrl,omitempty"`
	TargetHost          string          `json:"targetHost,omitempty"`
	RequestHeaders      json.RawMessage `json:"requestHeaders,omitempty"`
	RequestQueryParams  json.RawMessage `json:"requestQueryParams,omitempty"`
	RequestAuth         json.RawMessage `json:"requestAuth,omitempty"`
	RequestBody         json.RawMessage `json:"requestBody,omitempty"`
	ResponseStatus      *int            `json:"responseStatus,omitempty"`
	ResponseHeaders     json.RawMessage `json:"responseHeaders,omitempty"`
	ResponseBodyPreview string          `json:"responseBodyPreview,omitempty"`
	ResponseSizeBytes   *int            `json:"responseSizeBytes,omitempty"`
	ResponseTimeMS      *int            `json:"responseTimeMs,omitempty"`
	ResponseContentType string          `json:"responseContentType,omitempty"`
	RedirectCount       int             `json:"redirectCount"`
	BlockedReason       string          `json:"blockedReason,omitempty"`
	ErrorCode           string          `json:"errorCode,omitempty"`
	ErrorMessage        string          `json:"errorMessage,omitempty"`
	StartedAt           *time.Time      `json:"startedAt,omitempty"`
	CompletedAt         *time.Time      `json:"completedAt,omitempty"`
	CreatedAt           time.Time       `json:"createdAt"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
}

type ListQuery struct {
	UserID string
	Status string
	Method string
	Domain string
	Date   *time.Time
	Page   int32
	Limit  int32
}

type ListPagination struct {
	Page    int32 `json:"page"`
	Limit   int32 `json:"limit"`
	HasMore bool  `json:"hasMore"`
}

type ListResult struct {
	Items      []RunRecord
	Pagination ListPagination
}

type CreateParams struct {
	UserID              string
	CollectionID        *string
	SavedRequestID      *string
	Source              string
	Status              string
	Method              string
	URL                 string
	FinalURL            *string
	TargetHost          string
	RequestHeaders      json.RawMessage
	RequestQueryParams  json.RawMessage
	RequestAuth         json.RawMessage
	RequestBody         json.RawMessage
	ResponseStatus      *int
	ResponseHeaders     json.RawMessage
	ResponseBodyPreview string
	ResponseSizeBytes   *int
	ResponseTimeMS      *int
	ResponseContentType string
	RedirectCount       int
	BlockedReason       string
	ErrorCode           string
	ErrorMessage        string
	StartedAt           *time.Time
	CompletedAt         *time.Time
	Metadata            json.RawMessage
}

type Repository interface {
	ListByUser(ctx context.Context, userID string, limit int32) ([]RunRecord, error)
	ListByUserWithFilters(ctx context.Context, query ListQuery) ([]RunRecord, error)
	Create(ctx context.Context, params CreateParams) (RunRecord, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, userID string, limit int32) ([]RunRecord, error) {
	if s == nil || s.repo == nil {
		return nil, ErrUnavailable
	}
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrInvalid
	}
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListByUser(ctx, userID, limit)
}

func (s *Service) ListWithFilters(ctx context.Context, query ListQuery) (ListResult, error) {
	if s == nil || s.repo == nil {
		return ListResult{}, ErrUnavailable
	}

	query.UserID = strings.TrimSpace(query.UserID)
	if query.UserID == "" {
		return ListResult{}, ErrInvalid
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}
	query.Status = strings.TrimSpace(query.Status)
	switch method := strings.TrimSpace(query.Method); {
	case method == "":
		query.Method = ""
	case strings.EqualFold(method, "all"):
		query.Method = ""
	default:
		query.Method = strings.ToUpper(method)
	}
	query.Domain = strings.ToLower(strings.TrimSpace(query.Domain))

	fetchQuery := query
	fetchQuery.Limit = query.Limit + 1

	items, err := s.repo.ListByUserWithFilters(ctx, fetchQuery)
	if err != nil {
		return ListResult{}, err
	}

	hasMore := int32(len(items)) > query.Limit
	if hasMore {
		items = items[:query.Limit]
	}

	return ListResult{
		Items: items,
		Pagination: ListPagination{
			Page:    query.Page,
			Limit:   query.Limit,
			HasMore: hasMore,
		},
	}, nil
}

func (s *Service) Create(ctx context.Context, params CreateParams) (RunRecord, error) {
	if s == nil || s.repo == nil {
		return RunRecord{}, ErrUnavailable
	}
	params.UserID = strings.TrimSpace(params.UserID)
	params.Method = strings.ToUpper(strings.TrimSpace(params.Method))
	params.URL = strings.TrimSpace(params.URL)
	params.Source = strings.TrimSpace(params.Source)
	params.Status = strings.TrimSpace(params.Status)
	if params.UserID == "" || params.Method == "" || params.URL == "" {
		return RunRecord{}, ErrInvalid
	}
	if params.Source == "" {
		params.Source = "authenticated"
	}
	if params.Status == "" {
		params.Status = "succeeded"
	}
	for _, value := range []*json.RawMessage{&params.RequestHeaders, &params.RequestQueryParams, &params.RequestAuth, &params.RequestBody, &params.ResponseHeaders, &params.Metadata} {
		if len(*value) == 0 {
			*value = json.RawMessage(`{}`)
		}
	}
	return s.repo.Create(ctx, params)
}
