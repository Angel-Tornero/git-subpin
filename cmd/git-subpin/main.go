// Command git-subpin resolves Git submodule pointer conflicts left by a merge
// between environment branches, according to a declarative policy file.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	githubactions "github.com/sethvargo/go-githubactions"

	"github.com/Angel-Tornero/git-subpin/internal/gitutil"
	"github.com/Angel-Tornero/git-subpin/internal/policy"
	"github.com/Angel-Tornero/git-subpin/internal/resolver"
)

func main() {
	if err := run(context.Background(), githubactions.New(), os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "git-subpin:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, action *githubactions.Action, out io.Writer) error {
	from := action.GetInput("from")
	to := action.GetInput("to")
	configPath := action.GetInput("config")
	if from == "" || to == "" {
		return errors.New("inputs 'from' and 'to' are required")
	}
	if configPath == "" {
		configPath = ".subpin.yml"
	}

	f, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("opening policy file: %w", err)
	}
	defer f.Close()

	pol, err := policy.Load(f)
	if err != nil {
		return err
	}
	if err := pol.Validate(); err != nil {
		return fmt.Errorf("invalid policy: %w", err)
	}

	summary, err := resolver.Resolve(ctx, gitutil.New(""), pol.Rules, from, to)
	if err != nil {
		return err
	}

	if len(summary.Resolved) == 0 {
		fmt.Fprintln(out, "no submodule pointer conflicts to resolve")
		return nil
	}
	for _, r := range summary.Resolved {
		fmt.Fprintf(out, "resolved %s (%s) with %s -> %s\n", r.Submodule, r.Path, r.Strategy, r.SHA)
	}
	return nil
}
