package requests

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrUnavailable = errors.New("requests repository is unavailable")
	ErrNotFound    = errors.New("saved request not found")
	ErrInvalid     = errors.New("invalid saved request input")
)

type SavedRequest struct {
	ID              string          `json:"id"`
	CollectionID    *string         `json:"collectionId,omitempty"`
	OwnerUserID     *string         `json:"ownerUserId,omitempty"`
	Name            string          `json:"name"`
	Description     string          `json:"description,omitempty"`
	Method          string          `json:"method"`
	URL             string          `json:"url"`
	QueryParams     json.RawMessage `json:"queryParams,omitempty"`
	Headers         json.RawMessage `json:"headers,omitempty"`
	AuthScheme      string          `json:"authScheme"`
	AuthConfig      json.RawMessage `json:"authConfig,omitempty"`
	BodyMode        string          `json:"bodyMode"`
	BodyConfig      json.RawMessage `json:"bodyConfig,omitempty"`
	ExampleResponse json.RawMessage `json:"exampleResponse,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	DeletedAt       *time.Time      `json:"deletedAt,omitempty"`
}

type CreateParams struct {
	CollectionID    *string
	OwnerUserID     string
	Name            string
	Description     string
	Method          string
	URL             string
	QueryParams     json.RawMessage
	Headers         json.RawMessage
	AuthScheme      string
	AuthConfig      json.RawMessage
	BodyMode        string
	BodyConfig      json.RawMessage
	ExampleResponse json.RawMessage
	Metadata        json.RawMessage
}

type UpdateParams struct {
	ID              string
	OwnerUserID     string
	CollectionID    **string
	Name            *string
	Description     *string
	Method          *string
	URL             *string
	QueryParams     *json.RawMessage
	Headers         *json.RawMessage
	AuthScheme      *string
	AuthConfig      *json.RawMessage
	BodyMode        *string
	BodyConfig      *json.RawMessage
	ExampleResponse *json.RawMessage
	Metadata        *json.RawMessage
}

type Repository interface {
	GetByID(ctx context.Context, id string, ownerUserID string) (SavedRequest, error)
	Create(ctx context.Context, params CreateParams) (SavedRequest, error)
	Update(ctx context.Context, params UpdateParams) (SavedRequest, error)
	Delete(ctx context.Context, id string, ownerUserID string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, id string, ownerUserID string) (SavedRequest, error) {
	if s == nil || s.repo == nil {
		return SavedRequest{}, ErrUnavailable
	}
	if strings.TrimSpace(id) == "" || strings.TrimSpace(ownerUserID) == "" {
		return SavedRequest{}, ErrInvalid
	}
	return s.repo.GetByID(ctx, strings.TrimSpace(id), strings.TrimSpace(ownerUserID))
}

func (s *Service) Create(ctx context.Context, params CreateParams) (SavedRequest, error) {
	if s == nil || s.repo == nil {
		return SavedRequest{}, ErrUnavailable
	}
	normalized, err := normalizeCreate(params)
	if err != nil {
		return SavedRequest{}, err
	}
	return s.repo.Create(ctx, normalized)
}

func (s *Service) Update(ctx context.Context, params UpdateParams) (SavedRequest, error) {
	if s == nil || s.repo == nil {
		return SavedRequest{}, ErrUnavailable
	}
	normalized, err := normalizeUpdate(params)
	if err != nil {
		return SavedRequest{}, err
	}
	return s.repo.Update(ctx, normalized)
}

func (s *Service) Delete(ctx context.Context, id string, ownerUserID string) error {
	if s == nil || s.repo == nil {
		return ErrUnavailable
	}
	if strings.TrimSpace(id) == "" || strings.TrimSpace(ownerUserID) == "" {
		return ErrInvalid
	}
	return s.repo.Delete(ctx, strings.TrimSpace(id), strings.TrimSpace(ownerUserID))
}

func normalizeCreate(params CreateParams) (CreateParams, error) {
	params.OwnerUserID = strings.TrimSpace(params.OwnerUserID)
	params.Name = strings.TrimSpace(params.Name)
	params.Description = strings.TrimSpace(params.Description)
	params.Method = strings.ToUpper(strings.TrimSpace(params.Method))
	params.URL = strings.TrimSpace(params.URL)
	params.AuthScheme = strings.TrimSpace(params.AuthScheme)
	params.BodyMode = strings.TrimSpace(params.BodyMode)

	if params.OwnerUserID == "" || params.Name == "" || params.Method == "" || params.URL == "" {
		return CreateParams{}, ErrInvalid
	}
	if params.AuthScheme == "" {
		params.AuthScheme = "none"
	}
	if params.BodyMode == "" {
		params.BodyMode = "none"
	}
	for _, value := range []*json.RawMessage{&params.QueryParams, &params.Headers, &params.AuthConfig, &params.BodyConfig, &params.ExampleResponse, &params.Metadata} {
		if len(*value) == 0 {
			*value = json.RawMessage(`{}`)
		}
	}
	return params, nil
}

func normalizeUpdate(params UpdateParams) (UpdateParams, error) {
	params.ID = strings.TrimSpace(params.ID)
	params.OwnerUserID = strings.TrimSpace(params.OwnerUserID)
	if params.ID == "" || params.OwnerUserID == "" {
		return UpdateParams{}, ErrInvalid
	}
	if params.Name != nil {
		value := strings.TrimSpace(*params.Name)
		if value == "" {
			return UpdateParams{}, ErrInvalid
		}
		params.Name = &value
	}
	if params.Description != nil {
		value := strings.TrimSpace(*params.Description)
		params.Description = &value
	}
	if params.Method != nil {
		value := strings.ToUpper(strings.TrimSpace(*params.Method))
		if value == "" {
			return UpdateParams{}, ErrInvalid
		}
		params.Method = &value
	}
	if params.URL != nil {
		value := strings.TrimSpace(*params.URL)
		if value == "" {
			return UpdateParams{}, ErrInvalid
		}
		params.URL = &value
	}
	if params.AuthScheme != nil {
		value := strings.TrimSpace(*params.AuthScheme)
		params.AuthScheme = &value
	}
	if params.BodyMode != nil {
		value := strings.TrimSpace(*params.BodyMode)
		params.BodyMode = &value
	}
	return params, nil
}
