package resolver

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Angel-Tornero/git-subpin/internal/gitutil"
	"github.com/Angel-Tornero/git-subpin/internal/policy"
)

type setCall struct {
	path, sha string
}

// fakeRepo is an in-memory Repo so the resolver can be tested without git.
type fakeRepo struct {
	conflicts []gitutil.ConflictStatus
	names     map[string]string
	setErr    error
	sets      []setCall
}

func (f *fakeRepo) Conflicts(context.Context) ([]gitutil.ConflictStatus, error) {
	return f.conflicts, nil
}

func (f *fakeRepo) SubmoduleNames(context.Context) (map[string]string, error) {
	return f.names, nil
}

func (f *fakeRepo) SetPointer(_ context.Context, path, sha string) error {
	if f.setErr != nil {
		return f.setErr
	}
	f.sets = append(f.sets, setCall{path, sha})
	return nil
}

func TestResolve(t *testing.T) {
	keycloakConflict := gitutil.ConflictStatus{
		Path:       "third_party/keycloak",
		StatusCode: "UU",
		OursSHA:    "oursSHA",
		TheirsSHA:  "theirsSHA",
	}
	names := map[string]string{"third_party/keycloak": "keycloak"}

	tests := []struct {
		name     string
		repo     *fakeRepo
		rules    []policy.Rule
		wantSets []setCall
		wantErr  string
	}{
		{
			name: "keep-source stages the theirs pointer",
			repo: &fakeRepo{conflicts: []gitutil.ConflictStatus{keycloakConflict}, names: names},
			rules: []policy.Rule{
				{From: "*", To: "*", Submodule: "*", Strategy: policy.KeepSource},
			},
			wantSets: []setCall{{"third_party/keycloak", "theirsSHA"}},
		},
		{
			name: "no conflicts is a no-op",
			repo: &fakeRepo{names: names},
			rules: []policy.Rule{
				{From: "*", To: "*", Submodule: "*", Strategy: policy.KeepSource},
			},
			wantSets: nil,
		},
		{
			name: "keep-target is unsupported in v1",
			repo: &fakeRepo{conflicts: []gitutil.ConflictStatus{keycloakConflict}, names: names},
			rules: []policy.Rule{
				{From: "*", To: "*", Submodule: "*", Strategy: policy.KeepTarget},
			},
			wantErr: "not supported in v1",
		},
		{
			name:    "no matching rule is an error",
			repo:    &fakeRepo{conflicts: []gitutil.ConflictStatus{keycloakConflict}, names: names},
			rules:   []policy.Rule{{From: "a", To: "b", Submodule: "other", Strategy: policy.KeepSource}},
			wantErr: "no rule matches",
		},
		{
			name: "ambiguous rules are an error",
			repo: &fakeRepo{conflicts: []gitutil.ConflictStatus{keycloakConflict}, names: names},
			rules: []policy.Rule{
				{From: "dev", To: "*", Submodule: "keycloak", Strategy: policy.KeepSource},
				{From: "*", To: "prod", Submodule: "keycloak", Strategy: policy.KeepTarget},
			},
			wantErr: "ambiguous",
		},
		{
			name:    "conflicted path missing from .gitmodules is an error",
			repo:    &fakeRepo{conflicts: []gitutil.ConflictStatus{keycloakConflict}, names: map[string]string{}},
			rules:   []policy.Rule{{From: "*", To: "*", Submodule: "*", Strategy: policy.KeepSource}},
			wantErr: "no submodule registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary, err := Resolve(context.Background(), tt.repo, tt.rules, "dev", "prod")

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("Resolve() expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Resolve() error = %q; want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Resolve() unexpected error: %v", err)
			}
			if len(summary.Resolved) != len(tt.wantSets) {
				t.Errorf("Resolve() resolved %d submodules; want %d", len(summary.Resolved), len(tt.wantSets))
			}
			if !equalSetCalls(tt.repo.sets, tt.wantSets) {
				t.Errorf("SetPointer calls = %v; want %v", tt.repo.sets, tt.wantSets)
			}
		})
	}
}

func TestResolveSetPointerError(t *testing.T) {
	repo := &fakeRepo{
		conflicts: []gitutil.ConflictStatus{{Path: "third_party/keycloak", TheirsSHA: "theirsSHA"}},
		names:     map[string]string{"third_party/keycloak": "keycloak"},
		setErr:    errors.New("boom"),
	}
	rules := []policy.Rule{{From: "*", To: "*", Submodule: "*", Strategy: policy.KeepSource}}

	_, err := Resolve(context.Background(), repo, rules, "dev", "prod")
	if err == nil || !strings.Contains(err.Error(), "resolving submodule") {
		t.Fatalf("Resolve() error = %v; want it to wrap the SetPointer failure", err)
	}
}

func equalSetCalls(a, b []setCall) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
