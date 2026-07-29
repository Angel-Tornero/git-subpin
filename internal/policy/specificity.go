package policy

import (
	"strings"
)

type Specificity int

const (
	FullWildcard Specificity = iota
	PartialWildcard
	Literal
)

type RuleSpecificity struct {
	From, To, Submodule Specificity
}

func specificityOf(pattern string) Specificity {
	switch {
	case pattern == "*":
		return FullWildcard
	case strings.Contains(pattern, "*"):
		return PartialWildcard
	default:
		return Literal
	}
}

func (r Rule) Specificity() RuleSpecificity {
	return RuleSpecificity{
		From:      specificityOf(r.From),
		To:        specificityOf(r.To),
		Submodule: specificityOf(r.Submodule),
	}
}

func dominates(a, b RuleSpecificity) bool {
	return a.From >= b.From && a.To >= b.To && a.Submodule >= b.Submodule &&
		(a.From > b.From || a.To > b.To || a.Submodule > b.Submodule)
}
