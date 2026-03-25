package contracts_test

import (
	"net/http"
	"testing"

	"api-testing-kit/server/internal/httpapi"
	"api-testing-kit/server/tests/integration/helpers"
)

func TestTemplatesContractIncludesGuestSafeExamples(t *testing.T) {
	t.Parallel()

	server := helpers.NewServer(t, httpapi.NewRouter())

	response, err := server.Client().Get(server.URL + "/api/v1/templates")
	if err != nil {
		t.Fatalf("GET /api/v1/templates: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET /api/v1/templates returned %d", response.StatusCode)
	}

	var listPayload struct {
		Templates []struct {
			Slug string `json:"slug"`
		} `json:"templates"`
	}

	helpers.DecodeJSONResponse(t, response, &listPayload)

	if len(listPayload.Templates) == 0 {
		t.Fatalf("expected at least one template")
	}

	response, err = server.Client().Get(server.URL + "/api/v1/templates/weather-demo")
	if err != nil {
		t.Fatalf("GET /api/v1/templates/weather-demo: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET /api/v1/templates/weather-demo returned %d", response.StatusCode)
	}

	var detailPayload struct {
		Slug string `json:"slug"`
	}

	helpers.DecodeJSONResponse(t, response, &detailPayload)

	if detailPayload.Slug != "weather-demo" {
		t.Fatalf("detail slug = %q", detailPayload.Slug)
	}
}
