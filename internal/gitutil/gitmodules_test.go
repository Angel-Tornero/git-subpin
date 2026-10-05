package gitutil

import (
	"maps"
	"strings"
	"testing"
)

func TestParseGitmodules(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		want      map[string]string
		wantError bool
	}{
		{
			name: "single submodule",
			content: `[submodule "keycloak"]
	path = third_party/keycloak
	url = https://example.com/keycloak.git`,
			want: map[string]string{"third_party/keycloak": "keycloak"},
		},
		{
			name: "multiple submodules",
			content: `[submodule "keycloak"]
	path = third_party/keycloak
	url = https://example.com/keycloak.git
[submodule "grpc"]
	path = third_party/grpc
	url = https://example.com/grpc.git`,
			want: map[string]string{
				"third_party/keycloak": "keycloak",
				"third_party/grpc":     "grpc",
			},
		},
		{
			name:    "name differs from path",
			content: "[submodule \"auth\"]\n\tpath = vendor/keycloak\n",
			want:    map[string]string{"vendor/keycloak": "auth"},
		},
		{
			name:    "empty file",
			content: "",
			want:    map[string]string{},
		},
		{
			name:      "path outside a section is an error",
			content:   "path = third_party/keycloak\n",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGitmodules(strings.NewReader(tt.content))
			if tt.wantError {
				if err == nil {
					t.Fatalf("ParseGitmodules() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseGitmodules() unexpected error: %v", err)
			}
			if !maps.Equal(tt.want, got) {
				t.Errorf("ParseGitmodules() = %v; want %v", got, tt.want)
			}
		})
	}
}
