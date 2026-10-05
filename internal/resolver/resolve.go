// Package resolver applies a policy to the submodule pointer conflicts left by
// a merge, staging the resolved pointers.
package resolver

import (
	"context"
	"fmt"

	"github.com/Angel-Tornero/git-subpin/internal/gitutil"
	"github.com/Angel-Tornero/git-subpin/internal/policy"
)

// Repo is the git access the resolver needs. It is satisfied by *gitutil.Git.
type Repo interface {
	// Conflicts returns the submodule pointer conflicts in the working tree.
	Conflicts(ctx context.Context) ([]gitutil.ConflictStatus, error)
	// SubmoduleNames maps each submodule path to its name.
	SubmoduleNames(ctx context.Context) (map[string]string, error)
	// SetPointer stages path's submodule gitlink at sha.
	SetPointer(ctx context.Context, path, sha string) error
}

// Resolution records how a single submodule conflict was resolved.
type Resolution struct {
	Submodule string
	Path      string
	Strategy  policy.Strategy
	SHA       string
}

// Summary is the result of resolving the conflicts in a merge.
type Summary struct {
	Resolved []Resolution
}

// Resolve resolves every submodule pointer conflict in repo according to rules,
// for a merge of source into target. Only keep-source is supported in v1: a
// matched keep-target or manual strategy, an unmatched submodule, or an
// ambiguous rule set stops resolution with an error so a human can intervene.
func Resolve(ctx context.Context, repo Repo, rules []policy.Rule, source, target string) (Summary, error) {
	conflicts, err := repo.Conflicts(ctx)
	if err != nil {
		return Summary{}, fmt.Errorf("listing conflicts: %w", err)
	}
	if len(conflicts) == 0 {
		return Summary{}, nil
	}

	names, err := repo.SubmoduleNames(ctx)
	if err != nil {
		return Summary{}, fmt.Errorf("reading submodule names: %w", err)
	}

	var summary Summary
	for _, c := range conflicts {
		name, ok := names[c.Path]
		if !ok {
			return summary, fmt.Errorf("no submodule registered in .gitmodules for conflicted path %q", c.Path)
		}

		matching := policy.MatchingRules(rules, source, target, name)
		if len(matching) == 0 {
			return summary, fmt.Errorf("no rule matches submodule %q for %s into %s", name, source, target)
		}
		rule, err := policy.SelectRule(matching)
		if err != nil {
			return summary, fmt.Errorf("selecting rule for submodule %q: %w", name, err)
		}
		if rule.Strategy != policy.KeepSource {
			return summary, fmt.Errorf("strategy %q for submodule %q is not supported in v1 (only keep-source)", rule.Strategy, name)
		}

		// keep-source keeps the source (merged-in) side, which git records as
		// the "theirs" stage during a merge of source into target.
		if err := repo.SetPointer(ctx, c.Path, c.TheirsSHA); err != nil {
			return summary, fmt.Errorf("resolving submodule %q: %w", name, err)
		}
		summary.Resolved = append(summary.Resolved, Resolution{
			Submodule: name,
			Path:      c.Path,
			Strategy:  rule.Strategy,
			SHA:       c.TheirsSHA,
		})
	}
	return summary, nil
}
