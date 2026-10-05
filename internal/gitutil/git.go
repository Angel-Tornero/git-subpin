// Package gitutil runs git commands and parses their output for the resolver.
package gitutil

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Git runs git commands inside a repository working directory.
type Git struct {
	dir string
}

// New returns a Git that runs commands in dir. An empty dir uses the current
// working directory.
func New(dir string) *Git {
	return &Git{dir: dir}
}

// run executes git with args and returns its standard output.
func (g *Git) run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = g.dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, msg)
		}
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return stdout.String(), nil
}

// Conflicts returns the submodule pointer conflicts in the working tree.
func (g *Git) Conflicts(ctx context.Context) ([]ConflictStatus, error) {
	out, err := g.run(ctx, "status", "--porcelain=v2")
	if err != nil {
		return nil, err
	}
	return ParseConflictStatus(out)
}

// SubmoduleNames maps each submodule path to its name, read from .gitmodules.
func (g *Git) SubmoduleNames(_ context.Context) (map[string]string, error) {
	f, err := os.Open(filepath.Join(g.dir, ".gitmodules"))
	if err != nil {
		return nil, fmt.Errorf("opening .gitmodules: %w", err)
	}
	defer f.Close()
	return ParseGitmodules(f)
}

// SetPointer stages path's submodule gitlink at sha, resolving its conflict.
func (g *Git) SetPointer(ctx context.Context, path, sha string) error {
	_, err := g.run(ctx, "update-index", "--cacheinfo", fmt.Sprintf("%s,%s,%s", SubmoduleOctalFileMode, sha, path))
	return err
}
