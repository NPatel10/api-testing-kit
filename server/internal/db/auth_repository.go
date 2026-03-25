package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"api-testing-kit/server/internal/auth"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}

func (r *AuthRepository) CreateUser(ctx context.Context, params auth.CreateUserParams) (auth.UserRecord, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO users (
			email,
			password_hash,
			display_name,
			email_verified_at,
			status,
			role,
			locale,
			timezone,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING
			id,
			email,
			password_hash,
			email_verified_at,
			COALESCE(display_name, ''),
			status,
			role,
			locale,
			timezone,
			last_login_at,
			created_at,
			updated_at,
			deleted_at
	`, params.Email, params.PasswordHash, params.DisplayName, params.EmailVerifiedAt, params.Status, params.Role, params.Locale, params.Timezone, params.CreatedAt, params.UpdatedAt)

	record, err := scanUserRecord(row.Scan)
	if err != nil {
		if isUniqueViolation(err) {
			return auth.UserRecord{}, auth.ErrConflict
		}

		return auth.UserRecord{}, err
	}

	return record, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (auth.UserRecord, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT
			id,
			email,
			password_hash,
			email_verified_at,
			COALESCE(display_name, ''),
			status,
			role,
			locale,
			timezone,
			last_login_at,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`, email)

	record, err := scanUserRecord(row.Scan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.UserRecord{}, auth.ErrNotFound
		}

		return auth.UserRecord{}, err
	}

	return record, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, id string) (auth.UserRecord, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT
			id,
			email,
			password_hash,
			email_verified_at,
			COALESCE(display_name, ''),
			status,
			role,
			locale,
			timezone,
			last_login_at,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`, id)

	record, err := scanUserRecord(row.Scan)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.UserRecord{}, auth.ErrNotFound
		}

		return auth.UserRecord{}, err
	}

	return record, nil
}

func (r *AuthRepository) UpdateUserLastLoginAt(ctx context.Context, userID string, lastLoginAt time.Time) error {
	ct, err := r.pool.Exec(ctx, `
		UPDATE users
		SET last_login_at = $2, updated_at = $2
		WHERE id = $1
	`, userID, lastLoginAt)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return auth.ErrNotFound
	}

	return nil
}

func (r *AuthRepository) CreateSession(ctx context.Context, params auth.CreateSessionParams) (auth.SessionRecord, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO sessions (
			user_id,
			session_token_hash,
			status,
			expires_at,
			last_seen_at,
			created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING
			id,
			user_id,
			session_token_hash,
			status,
			expires_at,
			last_seen_at,
			revoked_at,
			created_at
	`, params.UserID, params.TokenHash, params.Status, params.ExpiresAt, params.CreatedAt, params.CreatedAt)

	record, err := scanSessionRecord(row.Scan)
	if err != nil {
		if isUniqueViolation(err) {
			return auth.SessionRecord{}, auth.ErrConflict
		}

		return auth.SessionRecord{}, err
	}

	return record, nil
}

func (r *AuthRepository) GetSessionIdentityByTokenHash(ctx context.Context, tokenHash string) (auth.UserRecord, auth.SessionRecord, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT
			u.id,
			u.email,
			u.password_hash,
			u.email_verified_at,
			COALESCE(u.display_name, ''),
			u.status,
			u.role,
			u.locale,
			u.timezone,
			u.last_login_at,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			s.id,
			s.user_id,
			s.session_token_hash,
			s.status,
			s.expires_at,
			s.last_seen_at,
			s.revoked_at,
			s.created_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.session_token_hash = $1
		  AND s.status = 'active'
		  AND s.revoked_at IS NULL
		  AND s.expires_at > now()
		  AND u.deleted_at IS NULL
	`, tokenHash)

	var user auth.UserRecord
	var session auth.SessionRecord
	var emailVerifiedAt sql.NullTime
	var lastLoginAt sql.NullTime
	var deletedAt sql.NullTime
	var lastSeenAt sql.NullTime
	var revokedAt sql.NullTime
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&emailVerifiedAt,
		&user.DisplayName,
		&user.Status,
		&user.Role,
		&user.Locale,
		&user.Timezone,
		&lastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.Status,
		&session.ExpiresAt,
		&lastSeenAt,
		&revokedAt,
		&session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.UserRecord{}, auth.SessionRecord{}, auth.ErrNotFound
		}

		return auth.UserRecord{}, auth.SessionRecord{}, err
	}

	if emailVerifiedAt.Valid {
		user.EmailVerifiedAt = &emailVerifiedAt.Time
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	if lastSeenAt.Valid {
		session.LastSeenAt = &lastSeenAt.Time
	}

	if revokedAt.Valid {
		session.RevokedAt = &revokedAt.Time
	}

	return user, session, nil
}

func (r *AuthRepository) RevokeSession(ctx context.Context, sessionID string, revokedAt time.Time) error {
	ct, err := r.pool.Exec(ctx, `
		UPDATE sessions
		SET status = 'revoked',
			revoked_at = $2
		WHERE id = $1
	`, sessionID, revokedAt)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return auth.ErrNotFound
	}

	return nil
}

func scanUserRecord(scan func(dest ...any) error) (auth.UserRecord, error) {
	var record auth.UserRecord
	var emailVerifiedAt sql.NullTime
	var lastLoginAt sql.NullTime
	var deletedAt sql.NullTime

	err := scan(
		&record.ID,
		&record.Email,
		&record.PasswordHash,
		&emailVerifiedAt,
		&record.DisplayName,
		&record.Status,
		&record.Role,
		&record.Locale,
		&record.Timezone,
		&lastLoginAt,
		&record.CreatedAt,
		&record.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return auth.UserRecord{}, err
	}

	if emailVerifiedAt.Valid {
		record.EmailVerifiedAt = &emailVerifiedAt.Time
	}

	if lastLoginAt.Valid {
		record.LastLoginAt = &lastLoginAt.Time
	}

	if deletedAt.Valid {
		record.DeletedAt = &deletedAt.Time
	}

	return record, nil
}

func scanSessionRecord(scan func(dest ...any) error) (auth.SessionRecord, error) {
	var record auth.SessionRecord
	var lastSeenAt sql.NullTime
	var revokedAt sql.NullTime

	err := scan(
		&record.ID,
		&record.UserID,
		&record.TokenHash,
		&record.Status,
		&record.ExpiresAt,
		&lastSeenAt,
		&revokedAt,
		&record.CreatedAt,
	)
	if err != nil {
		return auth.SessionRecord{}, err
	}

	if lastSeenAt.Valid {
		record.LastSeenAt = &lastSeenAt.Time
	}

	if revokedAt.Valid {
		record.RevokedAt = &revokedAt.Time
	}

	return record, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505"
}
