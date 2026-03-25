package db

import (
	"context"
	"database/sql"
	"errors"

	"api-testing-kit/server/internal/collections"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CollectionRepository struct {
	pool *pgxpool.Pool
}

func NewCollectionRepository(pool *pgxpool.Pool) *CollectionRepository {
	return &CollectionRepository{pool: pool}
}

func (r *CollectionRepository) ListByOwner(ctx context.Context, ownerUserID string) ([]collections.Collection, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			id,
			owner_user_id,
			name,
			slug,
			COALESCE(description, ''),
			visibility,
			COALESCE(color, ''),
			sort_order,
			shared_token,
			metadata,
			created_at,
			updated_at,
			deleted_at
		FROM collections
		WHERE owner_user_id = $1
		  AND deleted_at IS NULL
		ORDER BY sort_order ASC, created_at ASC
	`, ownerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]collections.Collection, 0)
	for rows.Next() {
		item, err := scanCollection(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *CollectionRepository) Create(ctx context.Context, params collections.CreateParams) (collections.Collection, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO collections (
			owner_user_id,
			name,
			slug,
			description,
			visibility,
			color,
			sort_order,
			metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			id,
			owner_user_id,
			name,
			slug,
			COALESCE(description, ''),
			visibility,
			COALESCE(color, ''),
			sort_order,
			shared_token,
			metadata,
			created_at,
			updated_at,
			deleted_at
	`, params.OwnerUserID, params.Name, params.Slug, params.Description, params.Visibility, params.Color, params.SortOrder, params.Metadata)

	item, err := scanCollection(row.Scan)
	if err != nil {
		return collections.Collection{}, err
	}

	return item, nil
}

func (r *CollectionRepository) Update(ctx context.Context, params collections.UpdateParams) (collections.Collection, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE collections
		SET
			name = COALESCE($3, name),
			slug = CASE
				WHEN $4::text IS NULL THEN slug
				ELSE $4::text
			END,
			description = COALESCE($5, description),
			visibility = COALESCE($6, visibility),
			color = COALESCE($7, color),
			sort_order = COALESCE($8, sort_order),
			metadata = COALESCE($9, metadata)
		WHERE id = $1
		  AND owner_user_id = $2
		  AND deleted_at IS NULL
		RETURNING
			id,
			owner_user_id,
			name,
			slug,
			COALESCE(description, ''),
			visibility,
			COALESCE(color, ''),
			sort_order,
			shared_token,
			metadata,
			created_at,
			updated_at,
			deleted_at
	`, params.ID, params.OwnerUserID, params.Name, nullableSlug(params.Slug), params.Description, params.Visibility, params.Color, params.SortOrder, params.Metadata)

	item, err := scanCollection(row.Scan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return collections.Collection{}, collections.ErrNotFound
		}
		return collections.Collection{}, err
	}

	return item, nil
}

func (r *CollectionRepository) Delete(ctx context.Context, id string, ownerUserID string) error {
	commandTag, err := r.pool.Exec(ctx, `
		UPDATE collections
		SET deleted_at = now()
		WHERE id = $1
		  AND owner_user_id = $2
		  AND deleted_at IS NULL
	`, id, ownerUserID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return collections.ErrNotFound
	}

	return nil
}

func scanCollection(scan func(dest ...any) error) (collections.Collection, error) {
	var item collections.Collection
	var ownerUserID sql.NullString
	var slug sql.NullString
	var sharedToken sql.NullString
	var deletedAt sql.NullTime

	if err := scan(
		&item.ID,
		&ownerUserID,
		&item.Name,
		&slug,
		&item.Description,
		&item.Visibility,
		&item.Color,
		&item.SortOrder,
		&sharedToken,
		&item.Metadata,
		&item.CreatedAt,
		&item.UpdatedAt,
		&deletedAt,
	); err != nil {
		return collections.Collection{}, err
	}

	if ownerUserID.Valid {
		value := ownerUserID.String
		item.OwnerUserID = &value
	}
	if slug.Valid {
		value := slug.String
		item.Slug = &value
	}
	if sharedToken.Valid {
		value := sharedToken.String
		item.SharedToken = &value
	}
	if deletedAt.Valid {
		item.DeletedAt = &deletedAt.Time
	}

	return item, nil
}

func nullableSlug(slug **string) any {
	if slug == nil {
		return nil
	}
	if *slug == nil {
		return ""
	}
	return **slug
}
