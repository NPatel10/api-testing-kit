package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnavailable        = errors.New("auth repository is unavailable")
	ErrNotFound           = errors.New("auth record not found")
	ErrInvalidInput       = errors.New("invalid auth input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrConflict           = errors.New("auth record already exists")
)

const (
	DefaultSessionDuration = 30 * 24 * time.Hour
	DefaultUserStatus      = "active"
	DefaultUserRole        = "user"
	DefaultUserLocale      = "en"
	DefaultUserTimezone    = "Asia/Calcutta"
)

type UserRecord struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	PasswordHash    string     `json:"-"`
	EmailVerifiedAt *time.Time `json:"emailVerifiedAt,omitempty"`
	DisplayName     string     `json:"displayName,omitempty"`
	Status          string     `json:"status"`
	Role            string     `json:"role,omitempty"`
	Locale          string     `json:"locale,omitempty"`
	Timezone        string     `json:"timezone,omitempty"`
	LastLoginAt     *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty"`
}

type SessionRecord struct {
	ID         string     `json:"id"`
	UserID     string     `json:"userId"`
	TokenHash  string     `json:"-"`
	Status     string     `json:"status"`
	ExpiresAt  time.Time  `json:"expiresAt"`
	LastSeenAt *time.Time `json:"lastSeenAt,omitempty"`
	RevokedAt  *time.Time `json:"revokedAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
}

type CreateUserParams struct {
	Email           string
	PasswordHash    string
	DisplayName     string
	EmailVerifiedAt *time.Time
	Status          string
	Role            string
	Locale          string
	Timezone        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CreateSessionParams struct {
	UserID    string
	TokenHash string
	Status    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type Repository interface {
	CreateUser(ctx context.Context, params CreateUserParams) (UserRecord, error)
	GetUserByEmail(ctx context.Context, email string) (UserRecord, error)
	GetUserByID(ctx context.Context, id string) (UserRecord, error)
	UpdateUserLastLoginAt(ctx context.Context, userID string, lastLoginAt time.Time) error
	CreateSession(ctx context.Context, params CreateSessionParams) (SessionRecord, error)
	GetSessionIdentityByTokenHash(ctx context.Context, tokenHash string) (UserRecord, SessionRecord, error)
	RevokeSession(ctx context.Context, sessionID string, revokedAt time.Time) error
}

type SignupInput struct {
	Email       string
	Password    string
	DisplayName string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	User    UserRecord
	Session SessionRecord
	Token   string
}

type Service struct {
	repo            Repository
	now             func() time.Time
	tokenGenerator  func() (string, error)
	passwordCost    int
	sessionDuration time.Duration
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:            repo,
		now:             time.Now,
		tokenGenerator:  generateSessionToken,
		passwordCost:    bcrypt.DefaultCost,
		sessionDuration: DefaultSessionDuration,
	}
}

func (s *Service) Signup(ctx context.Context, input SignupInput) (AuthResult, error) {
	if err := validateCredentials(input.Email, input.Password); err != nil {
		return AuthResult{}, err
	}

	if s.repo == nil {
		return AuthResult{}, ErrUnavailable
	}

	email := normalizeEmail(input.Email)
	displayName := strings.TrimSpace(input.DisplayName)
	if existing, err := s.repo.GetUserByEmail(ctx, email); err == nil && existing.ID != "" {
		return AuthResult{}, ErrConflict
	} else if err != nil && !errors.Is(err, ErrNotFound) {
		return AuthResult{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), s.passwordCost)
	if err != nil {
		return AuthResult{}, err
	}

	now := s.now()
	verifiedAt := now
	user, err := s.repo.CreateUser(ctx, CreateUserParams{
		Email:           email,
		PasswordHash:    string(passwordHash),
		DisplayName:     displayName,
		EmailVerifiedAt: &verifiedAt,
		Status:          DefaultUserStatus,
		Role:            DefaultUserRole,
		Locale:          DefaultUserLocale,
		Timezone:        DefaultUserTimezone,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return AuthResult{}, err
	}

	return s.createSession(ctx, user)
}

func (s *Service) Login(ctx context.Context, input LoginInput) (AuthResult, error) {
	if err := validateCredentials(input.Email, input.Password); err != nil {
		return AuthResult{}, err
	}

	if s.repo == nil {
		return AuthResult{}, ErrUnavailable
	}

	user, err := s.repo.GetUserByEmail(ctx, normalizeEmail(input.Email))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return AuthResult{}, ErrInvalidCredentials
		}

		return AuthResult{}, err
	}

	if user.PasswordHash == "" {
		return AuthResult{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	return s.createSession(ctx, user)
}

func (s *Service) Logout(ctx context.Context, token string) error {
	if s.repo == nil {
		return ErrUnavailable
	}

	tokenHash := hashSessionToken(token)
	_, session, err := s.repo.GetSessionIdentityByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil
		}

		return err
	}

	return s.repo.RevokeSession(ctx, session.ID, s.now())
}

func (s *Service) CurrentUser(ctx context.Context, token string) (UserRecord, SessionRecord, error) {
	if s.repo == nil {
		return UserRecord{}, SessionRecord{}, ErrUnavailable
	}

	return s.repo.GetSessionIdentityByTokenHash(ctx, hashSessionToken(token))
}

func (s *Service) createSession(ctx context.Context, user UserRecord) (AuthResult, error) {
	token, err := s.tokenGenerator()
	if err != nil {
		return AuthResult{}, err
	}

	now := s.now()
	session, err := s.repo.CreateSession(ctx, CreateSessionParams{
		UserID:    user.ID,
		TokenHash: hashSessionToken(token),
		Status:    "active",
		ExpiresAt: now.Add(s.sessionDuration),
		CreatedAt: now,
	})
	if err != nil {
		return AuthResult{}, err
	}

	if err := s.repo.UpdateUserLastLoginAt(ctx, user.ID, now); err != nil {
		return AuthResult{}, err
	}

	session.TokenHash = ""
	return AuthResult{
		User:    user,
		Session: session,
		Token:   token,
	}, nil
}

func validateCredentials(email, password string) error {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return ErrInvalidInput
	}

	if !strings.Contains(email, "@") {
		return fmt.Errorf("%w: email must contain @", ErrInvalidInput)
	}

	if len(password) < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters", ErrInvalidInput)
	}

	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func generateSessionToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func hashSessionToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
