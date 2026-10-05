package policy

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Strategy is how a submodule pointer conflict should be resolved.
type Strategy string

const (
	// KeepSource keeps the pointer from the branch being merged in.
	KeepSource Strategy = "keep-source"
	// KeepTarget keeps the pointer already on the branch being merged into.
	KeepTarget Strategy = "keep-target"
	// Manual leaves the conflict for a human to resolve.
	Manual Strategy = "manual"
)

// Valid reports whether s is a known strategy.
func (s Strategy) Valid() bool {
	return s == KeepSource || s == KeepTarget || s == Manual
}

// UnmarshalYAML decodes a strategy value and rejects unknown ones, so an
// unsupported strategy such as fast-forward fails when the policy is parsed.
func (s *Strategy) UnmarshalYAML(value *yaml.Node) error {
	var raw string
	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("decoding strategy: %w", err)
	}

	candidate := Strategy(raw)
	if !candidate.Valid() {
		return fmt.Errorf("invalid strategy %q: must be one of keep-source, keep-target, manual", raw)
	}

	*s = candidate
	return nil
}
