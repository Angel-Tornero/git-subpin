package gitutil

import (
	"fmt"
	"strings"
)

// UnmergedEntryType is the porcelain v2 record type for a conflicted entry.
const UnmergedEntryType = "u"

// SubmoduleOctalFileMode is the git file mode for a submodule gitlink.
const SubmoduleOctalFileMode = "160000"

// unmergedFieldCount is the number of fields in a porcelain v2 unmerged entry.
// The final field is the path, which may itself contain spaces.
const unmergedFieldCount = 11

// ConflictStatus is a submodule pointer conflict parsed from git status.
type ConflictStatus struct {
	Path, StatusCode, BaseSHA, OursSHA, TheirsSHA string
}

// ParseConflictStatus extracts submodule pointer conflicts from the output of
// "git status --porcelain=v2". Non-conflict and non-submodule entries are
// ignored; a malformed unmerged entry is an error.
func ParseConflictStatus(statusStr string) ([]ConflictStatus, error) {
	var conflicts []ConflictStatus
	for _, line := range strings.Split(statusStr, "\n") {
		if line == "" {
			continue
		}
		// Split into at most unmergedFieldCount fields so that a path
		// containing spaces survives whole as the final field.
		fields := strings.SplitN(line, " ", unmergedFieldCount)
		if fields[0] != UnmergedEntryType {
			continue
		}
		if len(fields) != unmergedFieldCount {
			return nil, fmt.Errorf("could not parse unmerged status line (expected %d fields, got %d): %q", unmergedFieldCount, len(fields), line)
		}
		if fields[3] != SubmoduleOctalFileMode ||
			fields[4] != SubmoduleOctalFileMode ||
			fields[5] != SubmoduleOctalFileMode ||
			fields[6] != SubmoduleOctalFileMode {

			continue
		}
		conflicts = append(conflicts, ConflictStatus{
			Path:       fields[10],
			StatusCode: fields[1],
			BaseSHA:    fields[7],
			OursSHA:    fields[8],
			TheirsSHA:  fields[9],
		})
	}
	return conflicts, nil
}
