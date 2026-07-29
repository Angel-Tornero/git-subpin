package policy

import (
	"errors"
	"slices"
	"testing"
)

func TestSelectRule(t *testing.T) {
	tests := []struct {
		name               string
		candidates         []Rule
		want               *Rule
		wantError          bool
		wantAmbiguousRules []Rule
	}{
		{
			name:       "no candidates returns an error",
			candidates: nil,
			wantError:  true,
		},
		{
			name: "single candidate wins unconditionally",
			candidates: []Rule{
				{From: "prod", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
			},
			want: &Rule{From: "prod", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
		},
		{
			name: "literal rule dominates a wildcard catch-all",
			candidates: []Rule{
				{From: "*", To: "*", Submodule: "*", Strategy: KeepTarget},
				{From: "prod", To: "test", Submodule: "third_party/keycloak", Strategy: KeepSource},
			},
			want: &Rule{From: "prod", To: "test", Submodule: "third_party/keycloak", Strategy: KeepSource},
		},
		{
			name: "tie with matching strategy resolves without ambiguity",
			candidates: []Rule{
				{From: "*", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
				{From: "prod", To: "*", Submodule: "third_party/keycloak", Strategy: KeepTarget},
			},
			want: &Rule{From: "*", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
		},
		{
			name: "tie with differing strategy is ambiguous",
			candidates: []Rule{
				{From: "*", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
				{From: "prod", To: "*", Submodule: "third_party/keycloak", Strategy: KeepSource},
			},
			wantError: true,
			wantAmbiguousRules: []Rule{
				{From: "*", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
				{From: "prod", To: "*", Submodule: "third_party/keycloak", Strategy: KeepSource},
			},
		},
		{
			name: "a dominated third candidate is excluded from the ambiguity",
			candidates: []Rule{
				{From: "*", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
				{From: "prod", To: "*", Submodule: "third_party/keycloak", Strategy: KeepSource},
				{From: "*", To: "*", Submodule: "*", Strategy: Manual},
			},
			wantError: true,
			wantAmbiguousRules: []Rule{
				{From: "*", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
				{From: "prod", To: "*", Submodule: "third_party/keycloak", Strategy: KeepSource},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectRule(tt.candidates)

			if tt.wantError {
				if err == nil {
					t.Errorf("SelectRule(...) expected error, got nil")
					return
				}
				if tt.wantAmbiguousRules != nil {
					var ambiguousErr *AmbiguousRulesError
					if !errors.As(err, &ambiguousErr) {
						t.Fatalf("SelectRule(...) expected *AmbiguousRulesError, got %v", err)
					}
					if !slices.Equal(ambiguousErr.Rules, tt.wantAmbiguousRules) {
						t.Errorf("AmbiguousRulesError.Rules = %+v; want %+v", ambiguousErr.Rules, tt.wantAmbiguousRules)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("SelectRule(...) unexpected error: %v", err)
				return
			}
			if got == nil || *got != *tt.want {
				t.Errorf("SelectRule(...) = %+v; want %+v", got, tt.want)
			}
		})
	}
}
