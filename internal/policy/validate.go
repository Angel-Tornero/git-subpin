package policy

import "fmt"

// Validate checks that the policy version is supported and that every rule has
// the required from/to/submodule patterns and a known strategy.
func (p *Policy) Validate() error {
	if p.Version != "1" {
		return fmt.Errorf("unsupported policy version %q", p.Version)
	}
	var invalidRules []int
	for i, rule := range p.Rules {
		if rule.From == "" || rule.To == "" || rule.Submodule == "" || !rule.Strategy.Valid() {
			invalidRules = append(invalidRules, i+1)
		}
	}
	if len(invalidRules) > 0 {
		return fmt.Errorf("invalid rules (need from/to/submodule and a known strategy): %v", invalidRules)
	}
	return nil
}
