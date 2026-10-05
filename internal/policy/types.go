// Package policy parses and evaluates the declarative rules that decide how
// submodule pointer conflicts are resolved during a merge between branches.
package policy

// Rule resolves a submodule pointer conflict for any merge whose source, target
// and submodule name match the From, To and Submodule patterns.
type Rule struct {
	From, To, Submodule string
	Strategy            Strategy
	Description         string
}

// Policy is a versioned, ordered set of rules.
type Policy struct {
	Version string
	Rules   []Rule
}
