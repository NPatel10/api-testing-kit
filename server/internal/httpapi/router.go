package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"api-testing-kit/server/internal/auth"
	"api-testing-kit/server/internal/collections"
	"api-testing-kit/server/internal/db"
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
	registerAuthRoutes(mux, deps)
	registerCollectionRoutes(mux, deps)

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

func registerCollectionRoutes(mux *http.ServeMux, deps RouterDeps) {
	var service *collections.Service
	authService := deps.Auth
	if authService == nil && deps.Store != nil && deps.Store.Auth != nil {
		authService = auth.NewService(deps.Store.Auth)
	}
	if deps.Store != nil && deps.Store.Collections != nil {
		service = collections.NewService(deps.Store.Collections)
	}

	NewCollectionsHandler(service, authService).Register(mux)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
