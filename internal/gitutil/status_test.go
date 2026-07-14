package gitutil

import (
	"slices"
	"testing"
)

func TestParseConflictStatus(t *testing.T) {
	tests := []struct {
		name      string
		statusStr string
		want      []ConflictStatus
		wantError bool
	}{
		{
			name:      "valid UU conflict",
			statusStr: `u UU S... 160000 160000 160000 160000 baseSHA oursSHA theirsSHA third_party/keycloak`,
			want: []ConflictStatus{{
				Path:       "third_party/keycloak",
				StatusCode: "UU",
				BaseSHA:    "baseSHA",
				OursSHA:    "oursSHA",
				TheirsSHA:  "theirsSHA",
			}},
			wantError: false,
		},
		{
			name: "ignores non-conflict lines",
			statusStr: `u UU S... 160000 160000 160000 160000 baseSHA oursSHA theirsSHA third_party/keycloak
1 .M N... 100644 100644 100644 100644 abc123 abc123 abc123 some/other/file.go`,
			want: []ConflictStatus{{
				Path:       "third_party/keycloak",
				StatusCode: "UU",
				BaseSHA:    "baseSHA",
				OursSHA:    "oursSHA",
				TheirsSHA:  "theirsSHA",
			}},
			wantError: false,
		},
		{
			name:      "too few fields returns error",
			statusStr: `u UU S...`,
			want:      nil,
			wantError: true,
		},
		{
			name:      "wrong mode at base",
			statusStr: `u UU S... 1337 160000 160000 160000 baseSHA oursSHA theirsSHA third_party/keycloak`,
			want:      nil,
			wantError: false,
		},
		{
			name:      "wrong mode at ours",
			statusStr: `u UU S... 160000 1337 160000 160000 baseSHA oursSHA theirsSHA third_party/keycloak`,
			want:      nil,
			wantError: false,
		},
		{
			name:      "wrong mode at theirs",
			statusStr: `u UU S... 160000 160000 1337 160000 baseSHA oursSHA theirsSHA third_party/keycloak`,
			want:      nil,
			wantError: false,
		},
		{
			name:      "wrong mode at worktree",
			statusStr: `u UU S... 160000 160000 160000 1337 baseSHA oursSHA theirsSHA third_party/keycloak`,
			want:      nil,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConflictStatuses(tt.statusStr)
			if tt.wantError {
				if err == nil {
					t.Errorf("ParseConflictStatuses(%q) expected error, got nil", tt.statusStr)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseConflictStatuses(%q) unexpected error: %v", tt.statusStr, err)
				return
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("ParseConflictStatuses(%q) = %v; want %v", tt.statusStr, got, tt.want)
			}
		})
	}
}
