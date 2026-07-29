package policy

import (
	"regexp"
	"strings"
)

func wildcardToRegexp(pattern string) string {
	literals := strings.Split(pattern, "*")
	for i, literal := range literals {
		literals[i] = regexp.QuoteMeta(literal)
	}
	return "^" + strings.Join(literals, ".*") + "$"
}

func matchPattern(pattern string, value string) bool {
	result, _ := regexp.MatchString(wildcardToRegexp(pattern), value)
	return result
}

func MatchingRules(rules []Rule, from, to, submodule string) []Rule {
	var matchingRules []Rule
	for _, rule := range rules {
		if matchPattern(rule.From, from) && matchPattern(rule.To, to) && matchPattern(rule.Submodule, submodule) {
			matchingRules = append(matchingRules, rule)
		}
	}
	return matchingRules
}
