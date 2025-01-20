// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package patching

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wk-y/diff"
	"github.com/wk-y/diff/internal/strutils"
)

// Parses a diff text back into hunks.
func ParseHunks(diffString string) ([]Hunk, error) {
	hunks := []Hunk{}

	lines := strutils.SplitLines(diffString)

	// Ensure the last line of the diff has a \n.
	if len(lines) > 0 && !strings.HasSuffix(lines[len(lines)-1], "\n") {
		lines[len(lines)-1] += "\n"
	}

	for _, line := range lines {
		if len(line) == 0 {
			return hunks, errors.New("blank line encountered in diff")
		}
		switch line[0] {
		case '@':
			newHunk, err := parseHunkHeader(line)
			if err != nil {
				return hunks, err
			}
			hunks = append(hunks, newHunk)
		case '+':
			if len(hunks) == 0 {
				return hunks, errors.New("addition encountered before any hunks")
			}
			hunks[len(hunks)-1].parts = append(hunks[len(hunks)-1].parts, diff.DiffPart{
				Action: diff.DiffAdded,
				Value:  line[1:],
			})
		case '-':
			if len(hunks) == 0 {
				return hunks, errors.New("removal encountered before any hunks")
			}
			hunks[len(hunks)-1].parts = append(hunks[len(hunks)-1].parts, diff.DiffPart{
				Action: diff.DiffRemoved,
				Value:  line[1:],
			})
		case ' ':
			if len(hunks) == 0 {
				return hunks, errors.New("identical encountered before any hunks")
			}
			hunks[len(hunks)-1].parts = append(hunks[len(hunks)-1].parts, diff.DiffPart{
				Action: diff.DiffIdentical,
				Value:  line[1:],
			})
		case '\\':
			// The \ indicator is somewhat special in that it is used to indicate
			// that the previous line has a missing newline.
			if len(hunks) == 0 {
				return hunks, errors.New("newline omission encountered before any hunks")
			}
			lastHunk := &hunks[len(hunks)-1]
			if len(lastHunk.parts) == 0 {
				return hunks, errors.New("newline omission encountered before any lines")
			}
		}
	}
	return hunks, nil
}

func parseHunkHeader(line string) (Hunk, error) {
	var hunk Hunk
	n, _ := fmt.Sscanf(line, "@@ -%d,%d +%d,%d @@", &hunk.aStart, &hunk.aLines, &hunk.bStart, &hunk.bLines)
	if n != 4 {
		return hunk, fmt.Errorf("Hunk header only had %v of the 4 fields expected", n)
	}
	return hunk, nil
}
