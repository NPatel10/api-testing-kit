package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"api-testing-kit/server/internal/auth"
	"api-testing-kit/server/internal/entitlements"
)

const sessionCookieName = "api_testing_kit_session"

type AuthHandler struct {
	service *auth.Service
}

type authEnvelope struct {
	User         userResponse        `json:"user"`
	Session      sessionResponse     `json:"session"`
	Entitlements entitlementResponse `json:"entitlements"`
}

type userResponse struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	EmailVerifiedAt *time.Time `json:"emailVerifiedAt,omitempty"`
	DisplayName     string     `json:"displayName,omitempty"`
	Status          string     `json:"status"`
	Role            string     `json:"role,omitempty"`
	Locale          string     `json:"locale,omitempty"`
	Timezone        string     `json:"timezone,omitempty"`
	LastLoginAt     *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type sessionResponse struct {
	ID        string    `json:"id"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type entitlementResponse struct {
	Plan         planResponse         `json:"plan"`
	Capabilities []capabilityResponse `json:"capabilities"`
}

type planResponse struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Source string `json:"source"`
}

type capabilityResponse struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Scope       string `json:"scope"`
	Limit       *int   `json:"limit,omitempty"`
	LimitLabel  string `json:"limitLabel,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

type authRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName,omitempty"`
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/signup", h.handleSignup)
	mux.HandleFunc("POST /api/v1/auth/login", h.handleLogin)
	mux.HandleFunc("POST /api/v1/auth/logout", h.handleLogout)
	mux.HandleFunc("GET /api/v1/me", h.handleMe)
}

func (h *AuthHandler) handleSignup(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "auth_unavailable", "authentication is temporarily unavailable")
		return
	}

	var input authRequest
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		return
	}

	result, err := h.service.Signup(r.Context(), auth.SignupInput{
		Email:       input.Email,
		Password:    input.Password,
		DisplayName: input.DisplayName,
	})
	if err != nil {
		writeAuthError(w, err)
		return
	}

	setSessionCookie(w, r, result.Token, result.Session.ExpiresAt)
	writeJSON(w, http.StatusCreated, authEnvelope{
		User:         toUserResponse(result.User),
		Session:      toSessionResponse(result.Session),
		Entitlements: toEntitlementResponse(entitlements.Resolve(result.User)),
	})
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "auth_unavailable", "authentication is temporarily unavailable")
		return
	}

	var input authRequest
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		return
	}

	result, err := h.service.Login(r.Context(), auth.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		writeAuthError(w, err)
		return
	}

	setSessionCookie(w, r, result.Token, result.Session.ExpiresAt)
	writeJSON(w, http.StatusOK, authEnvelope{
		User:         toUserResponse(result.User),
		Session:      toSessionResponse(result.Session),
		Entitlements: toEntitlementResponse(entitlements.Resolve(result.User)),
	})
}

func (h *AuthHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		clearSessionCookie(w, r)
		writeJSON(w, http.StatusOK, map[string]string{
			"status": "logged_out",
		})
		return
	}

	if cookie, err := r.Cookie(sessionCookieName); err == nil && cookie.Value != "" {
		_ = h.service.Logout(r.Context(), cookie.Value)
	}

	clearSessionCookie(w, r)
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "logged_out",
	})
}

func (h *AuthHandler) handleMe(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "auth_unavailable", "authentication is temporarily unavailable")
		return
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}

	user, session, err := h.service.CurrentUser(r.Context(), cookie.Value)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user":         toUserResponse(user),
		"session":      toSessionResponse(session),
		"entitlements": toEntitlementResponse(entitlements.Resolve(user)),
	})
}

func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, "invalid_request", err.Error())
	case errors.Is(err, auth.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "email or password is incorrect")
	case errors.Is(err, auth.ErrConflict):
		writeError(w, http.StatusConflict, "email_exists", "an account with this email already exists")
	case errors.Is(err, auth.ErrNotFound):
		writeError(w, http.StatusUnauthorized, "unauthorized", "missing or expired session")
	case errors.Is(err, auth.ErrUnavailable):
		writeError(w, http.StatusServiceUnavailable, "auth_unavailable", "authentication is temporarily unavailable")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "unexpected authentication failure")
	}
}

func toUserResponse(user auth.UserRecord) userResponse {
	return userResponse{
		ID:              user.ID,
		Email:           user.Email,
		EmailVerifiedAt: user.EmailVerifiedAt,
		DisplayName:     user.DisplayName,
		Status:          user.Status,
		Role:            user.Role,
		Locale:          user.Locale,
		Timezone:        user.Timezone,
		LastLoginAt:     user.LastLoginAt,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

func toSessionResponse(session auth.SessionRecord) sessionResponse {
	return sessionResponse{
		ID:        session.ID,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
	}
}

func toEntitlementResponse(state entitlements.State) entitlementResponse {
	capabilities := make([]capabilityResponse, 0, len(state.Capabilities))
	for _, capability := range state.Capabilities {
		capabilities = append(capabilities, capabilityResponse{
			Key:         string(capability.Key),
			Label:       capability.Label,
			Description: capability.Description,
			Enabled:     capability.Enabled,
			Scope:       capability.Scope,
			Limit:       capability.Limit,
			LimitLabel:  capability.LimitLabel,
			Reason:      capability.Reason,
		})
	}

	return entitlementResponse{
		Plan: planResponse{
			Code:   state.Plan.Code,
			Name:   state.Plan.Name,
			Source: string(state.Plan.Source),
		},
		Capabilities: capabilities,
	}
}

func setSessionCookie(w http.ResponseWriter, r *http.Request, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   isSecureRequest(r),
	})
}

func clearSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   isSecureRequest(r),
	})
}

func isSecureRequest(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
