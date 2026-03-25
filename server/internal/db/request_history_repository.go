package db

import (
	"context"
	"database/sql"

	"api-testing-kit/server/internal/history"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestHistoryRepository struct {
	pool *pgxpool.Pool
}

func NewRequestHistoryRepository(pool *pgxpool.Pool) *RequestHistoryRepository {
	return &RequestHistoryRepository{pool: pool}
}

func (r *RequestHistoryRepository) ListByUser(ctx context.Context, userID string, limit int32) ([]history.RunRecord, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, collection_id, saved_request_id, source, status, method, url, final_url, COALESCE(target_host, ''), request_headers, request_query_params, request_auth, request_body, response_status, response_headers, COALESCE(response_body_preview, ''), response_size_bytes, response_time_ms, COALESCE(response_content_type, ''), redirect_count, COALESCE(blocked_reason, ''), COALESCE(error_code, ''), COALESCE(error_message, ''), started_at, completed_at, created_at, metadata
		FROM request_runs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]history.RunRecord, 0)
	for rows.Next() {
		item, err := scanRunRecord(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *RequestHistoryRepository) Create(ctx context.Context, params history.CreateParams) (history.RunRecord, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO request_runs (
			user_id, collection_id, saved_request_id, source, status, method, url, final_url, target_host, request_headers, request_query_params, request_auth, request_body, response_status, response_headers, response_body_preview, response_size_bytes, response_time_ms, response_content_type, redirect_count, blocked_reason, error_code, error_message, started_at, completed_at, metadata
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26)
		RETURNING id, user_id, collection_id, saved_request_id, source, status, method, url, final_url, COALESCE(target_host, ''), request_headers, request_query_params, request_auth, request_body, response_status, response_headers, COALESCE(response_body_preview, ''), response_size_bytes, response_time_ms, COALESCE(response_content_type, ''), redirect_count, COALESCE(blocked_reason, ''), COALESCE(error_code, ''), COALESCE(error_message, ''), started_at, completed_at, created_at, metadata
	`,
		params.UserID, params.CollectionID, params.SavedRequestID, params.Source, params.Status, params.Method, params.URL, params.FinalURL, params.TargetHost, params.RequestHeaders, params.RequestQueryParams, params.RequestAuth, params.RequestBody, params.ResponseStatus, params.ResponseHeaders, params.ResponseBodyPreview, params.ResponseSizeBytes, params.ResponseTimeMS, params.ResponseContentType, params.RedirectCount, params.BlockedReason, params.ErrorCode, params.ErrorMessage, params.StartedAt, params.CompletedAt, params.Metadata,
	)
	return scanRunRecord(row.Scan)
}

func scanRunRecord(scan func(dest ...any) error) (history.RunRecord, error) {
	var item history.RunRecord
	var userID sql.NullString
	var collectionID sql.NullString
	var savedRequestID sql.NullString
	var finalURL sql.NullString
	var responseStatus sql.NullInt32
	var responseSizeBytes sql.NullInt32
	var responseTimeMS sql.NullInt32
	var startedAt sql.NullTime
	var completedAt sql.NullTime

	if err := scan(
		&item.ID,
		&userID,
		&collectionID,
		&savedRequestID,
		&item.Source,
		&item.Status,
		&item.Method,
		&item.URL,
		&finalURL,
		&item.TargetHost,
		&item.RequestHeaders,
		&item.RequestQueryParams,
		&item.RequestAuth,
		&item.RequestBody,
		&responseStatus,
		&item.ResponseHeaders,
		&item.ResponseBodyPreview,
		&responseSizeBytes,
		&responseTimeMS,
		&item.ResponseContentType,
		&item.RedirectCount,
		&item.BlockedReason,
		&item.ErrorCode,
		&item.ErrorMessage,
		&startedAt,
		&completedAt,
		&item.CreatedAt,
		&item.Metadata,
	); err != nil {
		return history.RunRecord{}, err
	}
	if userID.Valid {
		value := userID.String
		item.UserID = &value
	}
	if collectionID.Valid {
		value := collectionID.String
		item.CollectionID = &value
	}
	if savedRequestID.Valid {
		value := savedRequestID.String
		item.SavedRequestID = &value
	}
	if finalURL.Valid {
		value := finalURL.String
		item.FinalURL = &value
	}
	if responseStatus.Valid {
		value := int(responseStatus.Int32)
		item.ResponseStatus = &value
	}
	if responseSizeBytes.Valid {
		value := int(responseSizeBytes.Int32)
		item.ResponseSizeBytes = &value
	}
	if responseTimeMS.Valid {
		value := int(responseTimeMS.Int32)
		item.ResponseTimeMS = &value
	}
	if startedAt.Valid {
		item.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		item.CompletedAt = &completedAt.Time
	}
	return item, nil
}
