package contracts_test

import (
	"net/http"
	"testing"

	"api-testing-kit/server/internal/httpapi"
	"api-testing-kit/server/tests/integration/helpers"
)

func TestHealthEndpointsRespondWithStableContract(t *testing.T) {
	t.Parallel()

	server := helpers.NewServer(t, httpapi.NewRouter())

	for _, path := range []string{"/healthz", "/api/v1/health"} {
		response, err := server.Client().Get(server.URL + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("GET %s returned %d", path, response.StatusCode)
		}

		var payload struct {
			Status  string `json:"status"`
			Service string `json:"service"`
		}

		helpers.DecodeJSONResponse(t, response, &payload)

		if payload.Status != "ok" {
			t.Fatalf("GET %s status = %q", path, payload.Status)
		}

		if payload.Service != "api-testing-kit-server" {
			t.Fatalf("GET %s service = %q", path, payload.Service)
		}
	}
}
