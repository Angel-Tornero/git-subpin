package policy

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Strategy string

const (
	KeepSource  Strategy = "keep-source"
	KeepTarget  Strategy = "keep-target"
	FastForward Strategy = "fast-forward"
	Manual      Strategy = "manual"
)

func (s Strategy) Valid() bool {
	return s == KeepSource || s == KeepTarget || s == FastForward || s == Manual
}

func (s *Strategy) UnmarshalYAML(value *yaml.Node) error {
	var raw string
	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("decoding strategy: %w", err)
	}

	candidate := Strategy(raw)
	if !candidate.Valid() {
		return fmt.Errorf("invalid strategy %q: must be one of keep-source, keep-target, fast-forward, manual", raw)
	}

	*s = candidate
	return nil
}
