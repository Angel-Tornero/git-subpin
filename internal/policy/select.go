package policy

import (
	"fmt"
	"strings"
)

type AmbiguousRulesError struct {
	Rules []Rule
}

func (e *AmbiguousRulesError) Error() string {
	descriptions := make([]string, len(e.Rules))
	for i, r := range e.Rules {
		descriptions[i] = fmt.Sprintf("%s -> %s (%s): %s", r.From, r.To, r.Submodule, r.Strategy)
	}
	return fmt.Sprintf("ambiguous rules, no clear winner: %s", strings.Join(descriptions, "; "))
}

func SelectRule(candidates []Rule) (*Rule, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidate rules to select from")
	}
	if len(candidates) == 1 {
		return &candidates[0], nil
	}
	var maximal []Rule
	for i, candidate := range candidates {
		dominated := false
		for j, other := range candidates {
			if i == j {
				continue
			}
			if dominates(other.Specificity(), candidate.Specificity()) {
				dominated = true
				break
			}
		}
		if !dominated {
			maximal = append(maximal, candidate)
		}
	}
	if len(maximal) == 1 {
		return &maximal[0], nil
	}
	sameStrategy := true
	for _, r := range maximal[1:] {
		if r.Strategy != maximal[0].Strategy {
			sameStrategy = false
			break
		}
	}
	if sameStrategy {
		return &maximal[0], nil
	}

	return nil, &AmbiguousRulesError{Rules: maximal}
}
