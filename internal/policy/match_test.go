package policy

import (
	"slices"
	"testing"
)

func TestMatchingRules(t *testing.T) {
	tests := []struct {
		name      string
		rules     []Rule
		from      string
		to        string
		submodule string
		want      []Rule
	}{
		{
			name: "exact literal match",
			rules: []Rule{
				{From: "prod", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
			},
			from:      "prod",
			to:        "test",
			submodule: "third_party/keycloak",
			want: []Rule{
				{From: "prod", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
			},
		},
		{
			name: "full wildcard matches anything",
			rules: []Rule{
				{From: "*", To: "*", Submodule: "*", Strategy: KeepTarget},
			},
			from:      "prod",
			to:        "test",
			submodule: "third_party/grpc",
			want: []Rule{
				{From: "*", To: "*", Submodule: "*", Strategy: KeepTarget},
			},
		},
		{
			name: "partial wildcard matches with the literal prefix",
			rules: []Rule{
				{From: "prod", To: "test.*", Submodule: "*", Strategy: KeepSource},
			},
			from:      "prod",
			to:        "test.client-a",
			submodule: "third_party/grpc",
			want: []Rule{
				{From: "prod", To: "test.*", Submodule: "*", Strategy: KeepSource},
			},
		},
		{
			name: "partial wildcard does not match without the literal prefix",
			rules: []Rule{
				{From: "prod", To: "test.*", Submodule: "*", Strategy: KeepSource},
			},
			from:      "prod",
			to:        "test",
			submodule: "third_party/grpc",
			want:      nil,
		},
		{
			name: "mismatch in a single field excludes the rule entirely",
			rules: []Rule{
				{From: "prod", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
			},
			from:      "prod",
			to:        "test",
			submodule: "third_party/grpc",
			want:      nil,
		},
		{
			name: "filters candidates out of a mixed set of rules",
			rules: []Rule{
				{From: "*", To: "*", Submodule: "*", Strategy: KeepTarget},
				{From: "prod", To: "test", Submodule: "third_party/keycloak", Strategy: KeepSource},
				{From: "staging", To: "test", Submodule: "*", Strategy: Manual},
			},
			from:      "prod",
			to:        "test",
			submodule: "third_party/keycloak",
			want: []Rule{
				{From: "*", To: "*", Submodule: "*", Strategy: KeepTarget},
				{From: "prod", To: "test", Submodule: "third_party/keycloak", Strategy: KeepSource},
			},
		},
		{
			name:      "no rules at all returns no candidates",
			rules:     nil,
			from:      "prod",
			to:        "test",
			submodule: "third_party/keycloak",
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchingRules(tt.rules, tt.from, tt.to, tt.submodule)
			if !slices.Equal(got, tt.want) {
				t.Errorf("MatchingRules(...) = %+v; want %+v", got, tt.want)
			}
		})
	}
}
