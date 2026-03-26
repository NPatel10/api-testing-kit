package db

import (
	"strings"
	"testing"
)

func TestBuildPlanEntitlementsQuery(t *testing.T) {
	t.Parallel()

	query := buildPlanEntitlementsByPlanIDQuery()
	for _, want := range []string{"FROM plan_entitlements", "WHERE plan_id = $1", "ORDER BY key ASC"} {
		if !strings.Contains(query, want) {
			t.Fatalf("expected query to contain %q, got %q", want, query)
		}
	}
}
