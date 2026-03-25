package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"api-testing-kit/server/internal/auth"
)

func TestAuthSignupLoginMeAndLogout(t *testing.T) {
	t.Parallel()

	fakeRepo := newFakeAuthRepo()
	handler := NewAuthHandler(auth.NewService(fakeRepo))
	mux := http.NewServeMux()
	handler.Register(mux)

	signupBody := `{"email":"User@example.com","password":"password123","displayName":"Ada"}`
	signupReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", strings.NewReader(signupBody))
	signupReq.Header.Set("Content-Type", "application/json")
	signupRR := httptest.NewRecorder()
	mux.ServeHTTP(signupRR, signupReq)

	if signupRR.Code != http.StatusCreated {
		t.Fatalf("expected signup status %d, got %d", http.StatusCreated, signupRR.Code)
	}

	cookie := signupRR.Result().Cookies()
	if len(cookie) != 1 || cookie[0].Name != sessionCookieName {
		t.Fatalf("expected session cookie to be set, got %#v", cookie)
	}

	var signupPayload struct {
		User struct {
			Email string `json:"email"`
		} `json:"user"`
	}
	if err := json.Unmarshal(signupRR.Body.Bytes(), &signupPayload); err != nil {
		t.Fatalf("failed to decode signup payload: %v", err)
	}
	if signupPayload.User.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", signupPayload.User.Email)
	}

	meReq := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	meReq.AddCookie(cookie[0])
	meRR := httptest.NewRecorder()
	mux.ServeHTTP(meRR, meReq)
	if meRR.Code != http.StatusOK {
		t.Fatalf("expected me status %d, got %d", http.StatusOK, meRR.Code)
	}

	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"email":"user@example.com","password":"password123"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRR := httptest.NewRecorder()
	mux.ServeHTTP(loginRR, loginReq)
	if loginRR.Code != http.StatusOK {
		t.Fatalf("expected login status %d, got %d", http.StatusOK, loginRR.Code)
	}

	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	logoutReq.AddCookie(cookie[0])
	logoutRR := httptest.NewRecorder()
	mux.ServeHTTP(logoutRR, logoutReq)
	if logoutRR.Code != http.StatusOK {
		t.Fatalf("expected logout status %d, got %d", http.StatusOK, logoutRR.Code)
	}

	postLogoutMeReq := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	postLogoutMeReq.AddCookie(cookie[0])
	postLogoutMeRR := httptest.NewRecorder()
	mux.ServeHTTP(postLogoutMeRR, postLogoutMeReq)
	if postLogoutMeRR.Code != http.StatusUnauthorized {
		t.Fatalf("expected me to be unauthorized after logout, got %d", postLogoutMeRR.Code)
	}
}

func TestAuthRoutesWithoutServiceAreUnavailable(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	NewAuthHandler(nil).Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", strings.NewReader(`{"email":"test@example.com","password":"password123"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected unavailable status %d, got %d", http.StatusServiceUnavailable, rr.Code)
	}
}

type fakeAuthRepo struct {
	users    map[string]auth.UserRecord
	sessions map[string]auth.SessionRecord
	nextID   int
}

func newFakeAuthRepo() *fakeAuthRepo {
	return &fakeAuthRepo{
		users:    make(map[string]auth.UserRecord),
		sessions: make(map[string]auth.SessionRecord),
	}
}

func (r *fakeAuthRepo) CreateUser(ctx context.Context, params auth.CreateUserParams) (auth.UserRecord, error) {
	if _, exists := r.users[params.Email]; exists {
		return auth.UserRecord{}, auth.ErrConflict
	}

	r.nextID++
	record := auth.UserRecord{
		ID:              userID(r.nextID),
		Email:           params.Email,
		PasswordHash:    params.PasswordHash,
		EmailVerifiedAt: params.EmailVerifiedAt,
		DisplayName:     params.DisplayName,
		Status:          params.Status,
		Role:            params.Role,
		Locale:          params.Locale,
		Timezone:        params.Timezone,
		CreatedAt:       params.CreatedAt,
		UpdatedAt:       params.UpdatedAt,
	}
	r.users[params.Email] = record
	return record, nil
}

func (r *fakeAuthRepo) GetUserByEmail(ctx context.Context, email string) (auth.UserRecord, error) {
	record, ok := r.users[email]
	if !ok {
		return auth.UserRecord{}, auth.ErrNotFound
	}
	return record, nil
}

func (r *fakeAuthRepo) GetUserByID(ctx context.Context, id string) (auth.UserRecord, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return auth.UserRecord{}, auth.ErrNotFound
}

func (r *fakeAuthRepo) UpdateUserLastLoginAt(ctx context.Context, userID string, lastLoginAt time.Time) error {
	for email, user := range r.users {
		if user.ID == userID {
			user.LastLoginAt = &lastLoginAt
			user.UpdatedAt = lastLoginAt
			r.users[email] = user
			return nil
		}
	}
	return auth.ErrNotFound
}

func (r *fakeAuthRepo) CreateSession(ctx context.Context, params auth.CreateSessionParams) (auth.SessionRecord, error) {
	r.nextID++
	record := auth.SessionRecord{
		ID:        sessionID(r.nextID),
		UserID:    params.UserID,
		TokenHash: params.TokenHash,
		Status:    params.Status,
		ExpiresAt: params.ExpiresAt,
		CreatedAt: params.CreatedAt,
	}
	r.sessions[params.TokenHash] = record
	return record, nil
}

func (r *fakeAuthRepo) GetSessionIdentityByTokenHash(ctx context.Context, tokenHash string) (auth.UserRecord, auth.SessionRecord, error) {
	session, ok := r.sessions[tokenHash]
	if !ok || session.Status != "active" || session.RevokedAt != nil || time.Now().After(session.ExpiresAt) {
		return auth.UserRecord{}, auth.SessionRecord{}, auth.ErrNotFound
	}

	user, err := r.GetUserByID(ctx, session.UserID)
	if err != nil {
		return auth.UserRecord{}, auth.SessionRecord{}, err
	}

	return user, session, nil
}

func (r *fakeAuthRepo) RevokeSession(ctx context.Context, sessionID string, revokedAt time.Time) error {
	for tokenHash, session := range r.sessions {
		if session.ID == sessionID {
			session.Status = "revoked"
			session.RevokedAt = &revokedAt
			r.sessions[tokenHash] = session
			return nil
		}
	}
	return auth.ErrNotFound
}

func userID(n int) string {
	return fmt.Sprintf("user-%d", n)
}

func sessionID(n int) string {
	return fmt.Sprintf("session-%d", n)
}
