package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"api-testing-kit/server/internal/auth"
	"api-testing-kit/server/internal/billing"
	"api-testing-kit/server/internal/collections"
	"api-testing-kit/server/internal/db"
	"api-testing-kit/server/internal/guest"
	"api-testing-kit/server/internal/history"
	"api-testing-kit/server/internal/ratelimit"
	"api-testing-kit/server/internal/requests"
	"api-testing-kit/server/internal/runner"
	"api-testing-kit/server/internal/safety"
	"api-testing-kit/server/internal/templates"
)

type healthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

type RouterDeps struct {
	Store *db.Store
	Auth  *auth.Service
}

func NewRouter(deps RouterDeps) http.Handler {
	mux := http.NewServeMux()

	registerCoreRoutes(mux)
	registerTemplateRoutes(mux)
	registerGuestRoutes(mux, deps)
	registerAuthRoutes(mux, deps)
	registerBillingRoutes(mux)
	registerAdminRoutes(mux, deps)
	registerCollectionRoutes(mux, deps)
	registerSavedRequestRoutes(mux, deps)
	registerRunRoutes(mux, deps)

	return mux
}

func registerCoreRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "API Testing Kit server scaffold is running",
		})
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{
			Status:    "ok",
			Service:   "api-testing-kit-server",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{
			Status:    "ok",
			Service:   "api-testing-kit-server",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	})
}

func registerTemplateRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/templates", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"templates": templates.List(),
		})
	})

	mux.HandleFunc("GET /api/v1/templates/{slug}", func(w http.ResponseWriter, r *http.Request) {
		template, ok := templates.Get(r.PathValue("slug"))
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]any{
				"error": map[string]string{
					"code":    "template_not_found",
					"message": "template not found",
				},
			})
			return
		}

		writeJSON(w, http.StatusOK, template)
	})
}

func registerAuthRoutes(mux *http.ServeMux, deps RouterDeps) {
	service := deps.Auth
	if service == nil && deps.Store != nil && deps.Store.Auth != nil {
		service = auth.NewService(deps.Store.Auth)
	}

	NewAuthHandler(service).Register(mux)
}

func registerGuestRoutes(mux *http.ServeMux, deps RouterDeps) {
	var usageRecorder guest.UsageRecorder
	var abuseRecorder guest.AbuseRecorder
	if deps.Store != nil {
		usageRecorder = deps.Store.Usage
		abuseRecorder = deps.Store.Abuse
	}

	service := guest.NewService(nil, nil, usageRecorder, abuseRecorder, safety.DefaultOptions())
	NewGuestRunsHandler(service).Register(mux)
}

func registerBillingRoutes(mux *http.ServeMux) {
	NewBillingHandler(billing.NewStubService()).Register(mux)
}

func registerAdminRoutes(mux *http.ServeMux, deps RouterDeps) {
	authService := deps.Auth
	if authService == nil && deps.Store != nil && deps.Store.Auth != nil {
		authService = auth.NewService(deps.Store.Auth)
	}

	var abuseRepo adminAbuseRepository
	var blockedTargetsRepo adminBlockedTargetRepository
	if deps.Store != nil {
		abuseRepo = deps.Store.Abuse
		blockedTargetsRepo = deps.Store.BlockedTargets
	}

	NewAdminHandler(abuseRepo, blockedTargetsRepo, authService).Register(mux)
}

func registerCollectionRoutes(mux *http.ServeMux, deps RouterDeps) {
	var service *collections.Service
	var requestsService *requests.Service
	authService := deps.Auth
	if authService == nil && deps.Store != nil && deps.Store.Auth != nil {
		authService = auth.NewService(deps.Store.Auth)
	}
	if deps.Store != nil && deps.Store.Collections != nil {
		service = collections.NewService(deps.Store.Collections)
	}
	if deps.Store != nil && deps.Store.SavedRequests != nil {
		requestsService = requests.NewService(deps.Store.SavedRequests)
	}

	NewCollectionsHandler(service, requestsService, authService).Register(mux)
}

func registerSavedRequestRoutes(mux *http.ServeMux, deps RouterDeps) {
	var requestsService *requests.Service
	var historyService *history.Service
	authService := deps.Auth
	if authService == nil && deps.Store != nil && deps.Store.Auth != nil {
		authService = auth.NewService(deps.Store.Auth)
	}
	if deps.Store != nil && deps.Store.SavedRequests != nil {
		requestsService = requests.NewService(deps.Store.SavedRequests)
	}
	if deps.Store != nil && deps.Store.History != nil {
		historyService = history.NewService(deps.Store.History)
	}
	NewRequestsHandler(requestsService, historyService, authService).Register(mux)
}

func registerRunRoutes(mux *http.ServeMux, deps RouterDeps) {
	var historyService *history.Service
	var usageRecorder runner.UsageRecorder
	var abuseRecorder runner.AbuseRecorder
	authService := deps.Auth
	if authService == nil && deps.Store != nil && deps.Store.Auth != nil {
		authService = auth.NewService(deps.Store.Auth)
	}
	if deps.Store != nil && deps.Store.History != nil {
		historyService = history.NewService(deps.Store.History)
	}
	if deps.Store != nil {
		usageRecorder = deps.Store.Usage
		abuseRecorder = deps.Store.Abuse
	}
	runnerService := runner.NewService(
		nil,
		historyService,
		safety.DefaultOptions(),
		runner.WithLimiter(ratelimit.NewLimiter(ratelimit.AuthenticatedConfig())),
		runner.WithUsageRecorder(usageRecorder),
		runner.WithAbuseRecorder(abuseRecorder),
	)
	NewRunsHandler(runnerService, authService).Register(mux)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
