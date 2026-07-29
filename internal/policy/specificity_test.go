package policy

import (
	"testing"
)

func TestSpecificityOf(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    Specificity
	}{
		{name: "bare wildcard is full wildcard", pattern: "*", want: FullWildcard},
		{name: "suffix wildcard is partial", pattern: "test.*", want: PartialWildcard},
		{name: "prefix wildcard is partial", pattern: "*.test", want: PartialWildcard},
		{name: "no wildcard chars is literal", pattern: "prod", want: Literal},
		{name: "literal with slash is still literal", pattern: "third_party/keycloak", want: Literal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := specificityOf(tt.pattern)
			if got != tt.want {
				t.Errorf("specificityOf(%q) = %v; want %v", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestRuleSpecificity(t *testing.T) {
	tests := []struct {
		name string
		rule Rule
		want RuleSpecificity
	}{
		{
			name: "mixed specificity per field",
			rule: Rule{From: "*", To: "test.*", Submodule: "third_party/keycloak"},
			want: RuleSpecificity{From: FullWildcard, To: PartialWildcard, Submodule: Literal},
		},
		{
			name: "all literal",
			rule: Rule{From: "prod", To: "test", Submodule: "third_party/grpc"},
			want: RuleSpecificity{From: Literal, To: Literal, Submodule: Literal},
		},
		{
			name: "all full wildcard",
			rule: Rule{From: "*", To: "*", Submodule: "*"},
			want: RuleSpecificity{From: FullWildcard, To: FullWildcard, Submodule: FullWildcard},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.Specificity()
			if got != tt.want {
				t.Errorf("Specificity() = %+v; want %+v", got, tt.want)
			}
		})
	}
}

func TestDominates(t *testing.T) {
	tests := []struct {
		name string
		a    RuleSpecificity
		b    RuleSpecificity
		want bool
	}{
		{
			name: "literal in all fields dominates full wildcard in all fields",
			a:    RuleSpecificity{From: Literal, To: Literal, Submodule: Literal},
			b:    RuleSpecificity{From: FullWildcard, To: FullWildcard, Submodule: FullWildcard},
			want: true,
		},
		{
			name: "full wildcard never dominates literal",
			a:    RuleSpecificity{From: FullWildcard, To: FullWildcard, Submodule: FullWildcard},
			b:    RuleSpecificity{From: Literal, To: Literal, Submodule: Literal},
			want: false,
		},
		{
			name: "identical specificity does not dominate itself",
			a:    RuleSpecificity{From: PartialWildcard, To: PartialWildcard, Submodule: PartialWildcard},
			b:    RuleSpecificity{From: PartialWildcard, To: PartialWildcard, Submodule: PartialWildcard},
			want: false,
		},
		{
			name: "better in one field, tied in the rest, dominates",
			a:    RuleSpecificity{From: Literal, To: PartialWildcard, Submodule: PartialWildcard},
			b:    RuleSpecificity{From: PartialWildcard, To: PartialWildcard, Submodule: PartialWildcard},
			want: true,
		},
		{
			name: "worse in one field despite tying the rest does not dominate",
			a:    RuleSpecificity{From: PartialWildcard, To: PartialWildcard, Submodule: PartialWildcard},
			b:    RuleSpecificity{From: Literal, To: PartialWildcard, Submodule: PartialWildcard},
			want: false,
		},
		{
			name: "winning one field and losing another is not dominance (from wildcard, to literal vs from literal, to wildcard)",
			a:    RuleSpecificity{From: FullWildcard, To: Literal, Submodule: Literal},
			b:    RuleSpecificity{From: Literal, To: FullWildcard, Submodule: Literal},
			want: false,
		},
		{
			name: "reverse of the mixed-field case is also not dominance",
			a:    RuleSpecificity{From: Literal, To: FullWildcard, Submodule: Literal},
			b:    RuleSpecificity{From: FullWildcard, To: Literal, Submodule: Literal},
			want: false,
		},
		{
			name: "higher sum does not imply dominance",
			a:    RuleSpecificity{From: Literal, To: FullWildcard, Submodule: FullWildcard},
			b:    RuleSpecificity{From: PartialWildcard, To: PartialWildcard, Submodule: PartialWildcard},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dominates(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("dominates(%+v, %+v) = %v; want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
