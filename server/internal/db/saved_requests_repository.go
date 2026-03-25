package db

import (
	"context"
	"database/sql"
	"errors"

	"api-testing-kit/server/internal/requests"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SavedRequestRepository struct {
	pool *pgxpool.Pool
}

func NewSavedRequestRepository(pool *pgxpool.Pool) *SavedRequestRepository {
	return &SavedRequestRepository{pool: pool}
}

func (r *SavedRequestRepository) GetByID(ctx context.Context, id string, ownerUserID string) (requests.SavedRequest, error) {
	row := r.pool.QueryRow(ctx, savedRequestSelect+` WHERE id = $1 AND owner_user_id = $2 AND deleted_at IS NULL`, id, ownerUserID)
	item, err := scanSavedRequest(row.Scan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return requests.SavedRequest{}, requests.ErrNotFound
		}
		return requests.SavedRequest{}, err
	}
	return item, nil
}

func (r *SavedRequestRepository) Create(ctx context.Context, params requests.CreateParams) (requests.SavedRequest, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO saved_requests (
			collection_id, owner_user_id, name, description, method, url, query_params, headers, auth_scheme, auth_config, body_mode, body_config, example_response, metadata
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING id, collection_id, owner_user_id, name, COALESCE(description, ''), method, url, query_params, headers, auth_scheme, auth_config, body_mode, body_config, example_response, metadata, created_at, updated_at, deleted_at
	`,
		params.CollectionID, params.OwnerUserID, params.Name, params.Description, params.Method, params.URL, params.QueryParams, params.Headers, params.AuthScheme, params.AuthConfig, params.BodyMode, params.BodyConfig, params.ExampleResponse, params.Metadata,
	)
	return scanSavedRequest(row.Scan)
}

func (r *SavedRequestRepository) Update(ctx context.Context, params requests.UpdateParams) (requests.SavedRequest, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE saved_requests
		SET
			collection_id = CASE WHEN $3::uuid IS NULL THEN collection_id ELSE $3::uuid END,
			name = COALESCE($4, name),
			description = COALESCE($5, description),
			method = COALESCE($6, method),
			url = COALESCE($7, url),
			query_params = COALESCE($8, query_params),
			headers = COALESCE($9, headers),
			auth_scheme = COALESCE($10, auth_scheme),
			auth_config = COALESCE($11, auth_config),
			body_mode = COALESCE($12, body_mode),
			body_config = COALESCE($13, body_config),
			example_response = COALESCE($14, example_response),
			metadata = COALESCE($15, metadata)
		WHERE id = $1 AND owner_user_id = $2 AND deleted_at IS NULL
		RETURNING id, collection_id, owner_user_id, name, COALESCE(description, ''), method, url, query_params, headers, auth_scheme, auth_config, body_mode, body_config, example_response, metadata, created_at, updated_at, deleted_at
	`,
		params.ID, params.OwnerUserID, nullableUUID(params.CollectionID), params.Name, params.Description, params.Method, params.URL, params.QueryParams, params.Headers, params.AuthScheme, params.AuthConfig, params.BodyMode, params.BodyConfig, params.ExampleResponse, params.Metadata,
	)
	item, err := scanSavedRequest(row.Scan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return requests.SavedRequest{}, requests.ErrNotFound
		}
		return requests.SavedRequest{}, err
	}
	return item, nil
}

func (r *SavedRequestRepository) Delete(ctx context.Context, id string, ownerUserID string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE saved_requests
		SET deleted_at = now()
		WHERE id = $1 AND owner_user_id = $2 AND deleted_at IS NULL
	`, id, ownerUserID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return requests.ErrNotFound
	}
	return nil
}

const savedRequestSelect = `
	SELECT id, collection_id, owner_user_id, name, COALESCE(description, ''), method, url, query_params, headers, auth_scheme, auth_config, body_mode, body_config, example_response, metadata, created_at, updated_at, deleted_at
	FROM saved_requests
`

func scanSavedRequest(scan func(dest ...any) error) (requests.SavedRequest, error) {
	var item requests.SavedRequest
	var collectionID sql.NullString
	var ownerUserID sql.NullString
	var deletedAt sql.NullTime

	if err := scan(
		&item.ID,
		&collectionID,
		&ownerUserID,
		&item.Name,
		&item.Description,
		&item.Method,
		&item.URL,
		&item.QueryParams,
		&item.Headers,
		&item.AuthScheme,
		&item.AuthConfig,
		&item.BodyMode,
		&item.BodyConfig,
		&item.ExampleResponse,
		&item.Metadata,
		&item.CreatedAt,
		&item.UpdatedAt,
		&deletedAt,
	); err != nil {
		return requests.SavedRequest{}, err
	}
	if collectionID.Valid {
		value := collectionID.String
		item.CollectionID = &value
	}
	if ownerUserID.Valid {
		value := ownerUserID.String
		item.OwnerUserID = &value
	}
	if deletedAt.Valid {
		item.DeletedAt = &deletedAt.Time
	}
	return item, nil
}

func nullableUUID(value **string) any {
	if value == nil {
		return nil
	}
	if *value == nil {
		return ""
	}
	return **value
}
