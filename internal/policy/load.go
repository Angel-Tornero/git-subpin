package policy

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
)

func Load(policyYAML io.Reader) (*Policy, error) {
	decoder := yaml.NewDecoder(policyYAML)
	decoder.KnownFields(true)

	var policy Policy
	if err := decoder.Decode(&policy); err != nil {
		return nil, fmt.Errorf("parsing policy: %w", err)
	}

	return &policy, nil
}
