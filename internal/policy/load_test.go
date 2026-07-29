package policy

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name       string
		policyYAML io.Reader
		want       *Policy
		wantError  bool
	}{
		{
			name: "valid single rule",
			policyYAML: strings.NewReader(`
version: 1
rules:
  - from: prod
    to: test
    submodule: "*"
    strategy: keep-target
    description: "test always wins by default"
`),
			want: &Policy{
				Version: "1",
				Rules: []Rule{
					{From: "prod", To: "test", Submodule: "*", Strategy: KeepTarget, Description: "test always wins by default"},
				},
			},
		},
		{
			name: "multiple rules preserve order",
			policyYAML: strings.NewReader(`
version: 1
rules:
  - from: prod
    to: test
    submodule: "*"
    strategy: keep-target
  - from: prod
    to: test
    submodule: third_party/grpc
    strategy: keep-source
`),
			want: &Policy{
				Version: "1",
				Rules: []Rule{
					{From: "prod", To: "test", Submodule: "*", Strategy: KeepTarget},
					{From: "prod", To: "test", Submodule: "third_party/grpc", Strategy: KeepSource},
				},
			},
		},
		{
			name: "description is optional",
			policyYAML: strings.NewReader(`
version: 1
rules:
  - from: prod
    to: test
    submodule: third_party/keycloak
    strategy: manual
`),
			want: &Policy{
				Version: "1",
				Rules: []Rule{
					{From: "prod", To: "test", Submodule: "third_party/keycloak", Strategy: Manual},
				},
			},
		},
		{
			name: "invalid strategy value fails",
			policyYAML: strings.NewReader(`
version: 1
rules:
  - from: prod
    to: test
    submodule: third_party/keycloak
    strategy: keep-sourcee
`),
			wantError: true,
		},
		{
			name: "unknown field fails (typo)",
			policyYAML: strings.NewReader(`
version: 1
rules:
  - from: prod
    to: test
    submodule: third_party/grpc
    strategy: keep-target
    soruce: prod
`),
			wantError: true,
		},
		{
			name:       "malformed yaml syntax fails",
			policyYAML: strings.NewReader("version: 1\nrules:\n  - from: prod\n\tto: test\n"),
			wantError:  true,
		},
		{
			name: "version as quoted string still parses correctly",
			policyYAML: strings.NewReader(`
version: "1"
rules: []
`),
			want: &Policy{
				Version: "1",
				Rules:   []Rule{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.policyYAML)
			if tt.wantError {
				if err == nil {
					t.Errorf("Load() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Load() unexpected error: %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v; want %v", got, tt.want)
			}
		})
	}
}
