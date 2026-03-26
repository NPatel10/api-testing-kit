package httpapi

import (
	"context"
	"net/http"
	"strings"
	"time"

	"api-testing-kit/server/internal/abuse"
	"api-testing-kit/server/internal/auth"
)

type adminAbuseRepository interface {
	SummarizeByCategory(ctx context.Context, filter abuse.SummaryFilter) ([]abuse.SummaryRow, error)
	ListRecent(ctx context.Context, filter abuse.RecentFilter) ([]abuse.Event, error)
}

type adminBlockedTargetRepository interface {
	List(ctx context.Context, filter abuse.BlockedTargetFilter) ([]abuse.BlockedTarget, error)
}

type AdminHandler struct {
	auth           *auth.Service
	abuseRepo      adminAbuseRepository
	blockedTargets adminBlockedTargetRepository
}

type adminAbuseResponse struct {
	GeneratedAt    time.Time             `json:"generatedAt"`
	Summary        []abuse.SummaryRow    `json:"summary"`
	Recent         []abuse.Event         `json:"recent"`
	BlockedTargets []abuse.BlockedTarget `json:"blockedTargets"`
}

type adminBlockedTargetsResponse struct {
	GeneratedAt    time.Time             `json:"generatedAt"`
	BlockedTargets []abuse.BlockedTarget `json:"blockedTargets"`
}

func NewAdminHandler(abuseRepo adminAbuseRepository, blockedTargetsRepo adminBlockedTargetRepository, authService *auth.Service) *AdminHandler {
	return &AdminHandler{
		auth:           authService,
		abuseRepo:      abuseRepo,
		blockedTargets: blockedTargetsRepo,
	}
}

func (h *AdminHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/admin/abuse", h.handleAbuse)
	mux.HandleFunc("GET /api/v1/admin/blocked-targets", h.handleBlockedTargets)
}

func (h *AdminHandler) handleAbuse(w http.ResponseWriter, r *http.Request) {
	_, ok := h.requireAdminUser(w, r)
	if !ok {
		return
	}

	if h == nil || h.abuseRepo == nil || h.blockedTargets == nil {
		writeError(w, http.StatusServiceUnavailable, "admin_unavailable", "abuse monitoring is temporarily unavailable")
		return
	}

	summary, err := h.abuseRepo.SummarizeByCategory(r.Context(), abuse.SummaryFilter{Limit: 25})
	if err != nil {
		writeAdminError(w, err)
		return
	}

	recent, err := h.abuseRepo.ListRecent(r.Context(), abuse.RecentFilter{Limit: 25})
	if err != nil {
		writeAdminError(w, err)
		return
	}

	blockedTargets, err := h.blockedTargets.List(r.Context(), abuse.BlockedTargetFilter{
		IncludeExpired: true,
		Limit:          25,
	})
	if err != nil {
		writeAdminError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, adminAbuseResponse{
		GeneratedAt:    time.Now().UTC(),
		Summary:        summary,
		Recent:         recent,
		BlockedTargets: blockedTargets,
	})
}

func (h *AdminHandler) handleBlockedTargets(w http.ResponseWriter, r *http.Request) {
	_, ok := h.requireAdminUser(w, r)
	if !ok {
		return
	}

	if h == nil || h.blockedTargets == nil {
		writeError(w, http.StatusServiceUnavailable, "admin_unavailable", "blocked targets are temporarily unavailable")
		return
	}

	blockedTargets, err := h.blockedTargets.List(r.Context(), abuse.BlockedTargetFilter{
		IncludeExpired: true,
		Limit:          50,
	})
	if err != nil {
		writeAdminError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, adminBlockedTargetsResponse{
		GeneratedAt:    time.Now().UTC(),
		BlockedTargets: blockedTargets,
	})
}

func (h *AdminHandler) requireAdminUser(w http.ResponseWriter, r *http.Request) (auth.UserRecord, bool) {
	if h == nil || h.auth == nil {
		writeError(w, http.StatusServiceUnavailable, "admin_unavailable", "admin monitoring is temporarily unavailable")
		return auth.UserRecord{}, false
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "missing session")
		return auth.UserRecord{}, false
	}

	user, _, err := h.auth.CurrentUser(r.Context(), cookie.Value)
	if err != nil {
		writeAuthError(w, err)
		return auth.UserRecord{}, false
	}

	if !hasAdminAccess(user.Role) {
		writeError(w, http.StatusForbidden, "forbidden", "admin access required")
		return auth.UserRecord{}, false
	}

	return user, true
}

func hasAdminAccess(role string) bool {
	return strings.EqualFold(role, "admin") || strings.EqualFold(role, "owner")
}

func writeAdminError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusInternalServerError, "internal_error", "unexpected admin monitoring failure")
}
