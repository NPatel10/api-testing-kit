package collections

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrUnavailable = errors.New("collections repository is unavailable")
	ErrNotFound    = errors.New("collection not found")
	ErrInvalid     = errors.New("invalid collection input")
)

type Visibility string

const (
	VisibilityPrivate        Visibility = "private"
	VisibilitySharedReadonly Visibility = "shared_readonly"
	VisibilityInternal       Visibility = "internal"
)

type Collection struct {
	ID          string          `json:"id"`
	OwnerUserID *string         `json:"ownerUserId,omitempty"`
	Name        string          `json:"name"`
	Slug        *string         `json:"slug,omitempty"`
	Description string          `json:"description,omitempty"`
	Visibility  Visibility      `json:"visibility"`
	Color       string          `json:"color,omitempty"`
	SortOrder   int             `json:"sortOrder"`
	SharedToken *string         `json:"sharedToken,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	DeletedAt   *time.Time      `json:"deletedAt,omitempty"`
}

type CreateParams struct {
	OwnerUserID string
	Name        string
	Slug        *string
	Description string
	Visibility  Visibility
	Color       string
	SortOrder   int
	Metadata    json.RawMessage
}

type UpdateParams struct {
	ID          string
	OwnerUserID string
	Name        *string
	Slug        **string
	Description *string
	Visibility  *Visibility
	Color       *string
	SortOrder   *int
	Metadata    *json.RawMessage
}

type Repository interface {
	ListByOwner(ctx context.Context, ownerUserID string) ([]Collection, error)
	Create(ctx context.Context, params CreateParams) (Collection, error)
	Update(ctx context.Context, params UpdateParams) (Collection, error)
	Delete(ctx context.Context, id string, ownerUserID string) error
}

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
}

func (s *Service) List(ctx context.Context, ownerUserID string) ([]Collection, error) {
	if s == nil || s.repo == nil {
		return nil, ErrUnavailable
	}

	if strings.TrimSpace(ownerUserID) == "" {
		return nil, ErrInvalid
	}

	return s.repo.ListByOwner(ctx, ownerUserID)
}

func (s *Service) Create(ctx context.Context, params CreateParams) (Collection, error) {
	if s == nil || s.repo == nil {
		return Collection{}, ErrUnavailable
	}

	normalized, err := normalizeCreate(params)
	if err != nil {
		return Collection{}, err
	}

	return s.repo.Create(ctx, normalized)
}

func (s *Service) Update(ctx context.Context, params UpdateParams) (Collection, error) {
	if s == nil || s.repo == nil {
		return Collection{}, ErrUnavailable
	}

	normalized, err := normalizeUpdate(params)
	if err != nil {
		return Collection{}, err
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
	params.Color = strings.TrimSpace(params.Color)

	if params.OwnerUserID == "" || params.Name == "" {
		return CreateParams{}, ErrInvalid
	}

	if params.Slug != nil {
		slug := strings.TrimSpace(*params.Slug)
		params.Slug = &slug
		if slug == "" {
			params.Slug = nil
		}
	}

	if params.Visibility == "" {
		params.Visibility = VisibilityPrivate
	}

	switch params.Visibility {
	case VisibilityPrivate, VisibilitySharedReadonly, VisibilityInternal:
	default:
		return CreateParams{}, ErrInvalid
	}

	if len(params.Metadata) == 0 {
		params.Metadata = json.RawMessage(`{}`)
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

	if params.Slug != nil && *params.Slug != nil {
		value := strings.TrimSpace(**params.Slug)
		*params.Slug = &value
		if value == "" {
			*params.Slug = nil
		}
	}

	if params.Description != nil {
		value := strings.TrimSpace(*params.Description)
		params.Description = &value
	}

	if params.Color != nil {
		value := strings.TrimSpace(*params.Color)
		params.Color = &value
	}

	if params.Visibility != nil {
		switch *params.Visibility {
		case VisibilityPrivate, VisibilitySharedReadonly, VisibilityInternal:
		default:
			return UpdateParams{}, ErrInvalid
		}
	}

	return params, nil
}
