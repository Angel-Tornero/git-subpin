package policy

import "fmt"

func (p *Policy) Validate() error {
	if p.Version != "1" {
		return fmt.Errorf("unsupported yaml version %q", p.Version)
	}
	var invalidRules []int
	for i, rule := range p.Rules {
		if rule.From == "" || rule.To == "" || rule.Submodule == "" {
			invalidRules = append(invalidRules, i+1)
		}
	}
	if len(invalidRules) > 0 {
		return fmt.Errorf("missing required fields (from/to/submodule) in rules: %v", invalidRules)
	}
	return nil
}
