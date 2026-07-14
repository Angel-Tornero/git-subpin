package gitutil

import (
	"fmt"
	"strings"
)

const UnmergedEntryType = "u"
const SubmoduleOctalFileMode = "160000"

type ConflictStatus struct {
	Path, StatusCode, BaseSHA, OursSHA, TheirsSHA string
}

func ParseConflictStatuses(statusStr string) ([]ConflictStatus, error) {
	statusLines := strings.Split(statusStr, "\n")
	var conflicts []ConflictStatus
	for _, line := range statusLines {
		fields := strings.Split(line, " ")
		if fields[0] != UnmergedEntryType {
			continue
		}
		if len(fields) < 10 {
			return nil, fmt.Errorf("could not parse status line (expected at least 10 fields, got %d): %q", len(fields), line)
		}
		if fields[3] != SubmoduleOctalFileMode ||
			fields[4] != SubmoduleOctalFileMode ||
			fields[5] != SubmoduleOctalFileMode ||
			fields[6] != SubmoduleOctalFileMode {

			continue
		}
		path := fields[len(fields)-1]
		statusCode := fields[1]
		baseSHA := fields[7]
		oursSHA := fields[8]
		theirsSHA := fields[9]
		conflicts = append(conflicts, ConflictStatus{
			Path:       path,
			StatusCode: statusCode,
			BaseSHA:    baseSHA,
			OursSHA:    oursSHA,
			TheirsSHA:  theirsSHA,
		})
	}
	return conflicts, nil
}
