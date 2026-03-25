package helpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return server
}

func DecodeJSONResponse(t *testing.T, response *http.Response, dst any) {
	t.Helper()
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	if err := decoder.Decode(dst); err != nil {
		t.Fatalf("decode JSON response: %v", err)
	}
}
