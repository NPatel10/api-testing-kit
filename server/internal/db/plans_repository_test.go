package db

import (
	"strings"
	"testing"
)

func TestBuildPlanQueries(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name  string
		query string
		want  []string
	}{
		{
			name:  "list active",
			query: buildPlanListQuery(),
			want:  []string{"FROM plans", "WHERE is_active = TRUE", "ORDER BY sort_order ASC, code ASC"},
		},
		{
			name:  "by code",
			query: buildPlanByCodeQuery(),
			want:  []string{"FROM plans", "WHERE code = $1", "AND is_active = TRUE"},
		},
		{
			name:  "by id",
			query: buildPlanByIDQuery(),
			want:  []string{"FROM plans", "WHERE id = $1", "AND is_active = TRUE"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, want := range tc.want {
				if !strings.Contains(tc.query, want) {
					t.Fatalf("expected query to contain %q, got %q", want, tc.query)
				}
			}
		})
	}
}
