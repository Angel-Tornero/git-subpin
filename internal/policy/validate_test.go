package policy

import (
	"strings"
	"testing"
)

func TestPolicyValidate(t *testing.T) {
	tests := []struct {
		name      string
		policy    Policy
		wantError bool
		wantMsg   string
	}{
		{
			name: "valid policy passes",
			policy: Policy{
				Version: "1",
				Rules: []Rule{
					{From: "prod", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
				},
			},
		},
		{
			name: "empty rules list is valid",
			policy: Policy{
				Version: "1",
				Rules:   nil,
			},
		},
		{
			name: "unsupported version fails",
			policy: Policy{
				Version: "2",
			},
			wantError: true,
		},
		{
			name: "missing version fails",
			policy: Policy{
				Version: "",
			},
			wantError: true,
		},
		{
			name: "rule missing from fails",
			policy: Policy{
				Version: "1",
				Rules: []Rule{
					{From: "", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
				},
			},
			wantError: true,
		},
		{
			name: "rule missing to fails",
			policy: Policy{
				Version: "1",
				Rules: []Rule{
					{From: "prod", To: "", Submodule: "third_party/keycloak", Strategy: KeepTarget},
				},
			},
			wantError: true,
		},
		{
			name: "rule missing submodule fails",
			policy: Policy{
				Version: "1",
				Rules: []Rule{
					{From: "prod", To: "test", Submodule: "", Strategy: KeepTarget},
				},
			},
			wantError: true,
		},
		{
			name: "reports positions of multiple invalid rules",
			policy: Policy{
				Version: "1",
				Rules: []Rule{
					{From: "prod", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
					{From: "", To: "test", Submodule: "third_party/keycloak", Strategy: KeepTarget},
					{From: "prod", To: "test", Submodule: "third_party/grpc", Strategy: KeepSource},
					{From: "prod", To: "", Submodule: "third_party/grpc", Strategy: KeepSource},
				},
			},
			wantError: true,
			wantMsg:   "[2 4]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policy.Validate()

			if tt.wantError {
				if err == nil {
					t.Fatalf("Validate() expected error, got nil")
				}
				if tt.wantMsg != "" && !strings.Contains(err.Error(), tt.wantMsg) {
					t.Errorf("Validate() error = %q; want it to contain %q", err.Error(), tt.wantMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}
