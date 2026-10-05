package resolver

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Angel-Tornero/git-subpin/internal/gitutil"
	"github.com/Angel-Tornero/git-subpin/internal/policy"
)

// TestResolveIntegration drives Resolve against a real repository that has a
// genuine submodule pointer conflict, built with the git CLI.
func TestResolveIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test builds a real git repo; skipped in -short mode")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found in PATH")
	}

	// Upstream submodule with a base commit and two divergent children, so the
	// super-project pointer genuinely conflicts rather than fast-forwarding.
	sub := t.TempDir()
	git(t, sub, "init", "-q", "-b", "main")
	writeFile(t, sub, "VERSION", "A")
	git(t, sub, "add", "VERSION")
	git(t, sub, "commit", "-q", "-m", "A")
	shaA := git(t, sub, "rev-parse", "HEAD")

	git(t, sub, "checkout", "-q", "-b", "b", shaA)
	writeFile(t, sub, "VERSION", "B")
	git(t, sub, "commit", "-q", "-am", "B")
	shaB := git(t, sub, "rev-parse", "HEAD")

	git(t, sub, "checkout", "-q", "-b", "c", shaA)
	writeFile(t, sub, "VERSION", "C")
	git(t, sub, "commit", "-q", "-am", "C")
	shaC := git(t, sub, "rev-parse", "HEAD")

	// Super-project embedding the submodule.
	super := t.TempDir()
	subDir := filepath.Join(super, "third_party/keycloak")
	git(t, super, "init", "-q", "-b", "main")
	// -c propagates to the clone subprocess git spawns for the submodule, which
	// otherwise refuses the local file transport.
	git(t, super, "-c", "protocol.file.allow=always", "submodule", "add", "-q", sub, "third_party/keycloak")

	git(t, subDir, "checkout", "-q", shaA)
	git(t, super, "add", "third_party/keycloak")
	git(t, super, "commit", "-q", "-m", "base at A")

	// target (main) moves the pointer to B.
	git(t, subDir, "checkout", "-q", shaB)
	git(t, super, "add", "third_party/keycloak")
	git(t, super, "commit", "-q", "-m", "target moves to B")

	// source branches off base and moves the pointer to C.
	git(t, super, "checkout", "-q", "-b", "source", "main~1")
	git(t, subDir, "checkout", "-q", shaC)
	git(t, super, "add", "third_party/keycloak")
	git(t, super, "commit", "-q", "-m", "source moves to C")

	// Merging source into target conflicts on the submodule gitlink.
	git(t, super, "checkout", "-q", "main")
	if out, err := runGit(super, "merge", "--no-edit", "source"); err == nil {
		t.Fatalf("expected the merge to conflict, but it succeeded:\n%s", out)
	}

	ctx := context.Background()
	repo := gitutil.New(super)
	conflicts, err := repo.Conflicts(ctx)
	if err != nil {
		t.Fatalf("Conflicts(): %v", err)
	}
	if len(conflicts) != 1 {
		t.Fatalf("got %d conflicts, want 1: %+v", len(conflicts), conflicts)
	}

	rules := []policy.Rule{{From: "*", To: "*", Submodule: "*", Strategy: policy.KeepSource}}
	summary, err := Resolve(ctx, repo, rules, "source", "main")
	if err != nil {
		t.Fatalf("Resolve(): %v", err)
	}
	if len(summary.Resolved) != 1 {
		t.Fatalf("resolved %d submodules, want 1", len(summary.Resolved))
	}

	// keep-source keeps the source (theirs) pointer, staged at stage 0.
	if got := git(t, super, "rev-parse", ":third_party/keycloak"); got != shaC {
		t.Errorf("resolved pointer = %s; want source sha %s", got, shaC)
	}
	if remaining, err := repo.Conflicts(ctx); err != nil {
		t.Fatalf("Conflicts() after resolve: %v", err)
	} else if len(remaining) != 0 {
		t.Errorf("got %d conflicts after resolve, want 0", len(remaining))
	}
}

func gitEnv() []string {
	return append(os.Environ(),
		"GIT_CONFIG_GLOBAL="+os.DevNull,
		"GIT_CONFIG_SYSTEM="+os.DevNull,
		"GIT_AUTHOR_NAME=test",
		"GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test",
		"GIT_COMMITTER_EMAIL=test@example.com",
	)
}

func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = gitEnv()
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := runGit(dir, args...)
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(out)
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
