package gitutil

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var submoduleSectionRE = regexp.MustCompile(`^\[submodule "(.+)"\]$`)

// ParseGitmodules reads a .gitmodules file and returns a map from each
// submodule's path to its name (the quoted section header).
func ParseGitmodules(r io.Reader) (map[string]string, error) {
	names := make(map[string]string)
	var current string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if m := submoduleSectionRE.FindStringSubmatch(line); m != nil {
			current = m[1]
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(key) != "path" {
			continue
		}
		path := strings.TrimSpace(value)
		if current == "" {
			return nil, fmt.Errorf("path %q defined outside a submodule section", path)
		}
		names[path] = current
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading .gitmodules: %w", err)
	}
	return names, nil
}
