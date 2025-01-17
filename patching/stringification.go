// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package patching

import (
	"fmt"
	"strings"

	"github.com/wk-y/diff"
)

// DiffString formats an array of DiffParts into a unified diff.
// DiffString will produce strange results if d is not from LineDiff.
func DiffString(d []diff.DiffPart) string {
	hunks := HunkDiff(d)
	hunkStrings := make([]string, len(hunks))
	for i, hunk := range hunks {
		hunkStrings[i] = hunk.String()
	}
	return strings.Join(hunkStrings, "")
}

func (h Hunk) String() string {
	diffLines := make([]string, 0)
	header := fmt.Sprintf("@@ -%v,%v +%v,%v @@\n", h.aStart, h.aLines, h.bStart, h.bLines)
	diffLines = append(diffLines, header)
	for _, part := range h.parts {
		switch part.Action {
		case diff.DiffAdded:
			diffLines = append(diffLines, fmt.Sprint("+", part.Value))
		case diff.DiffRemoved:
			diffLines = append(diffLines, fmt.Sprint("-", part.Value))
		case diff.DiffIdentical:
			diffLines = append(diffLines, fmt.Sprint(" ", part.Value))
		}

		// Since it is assumed that LineDiff is used, it is assumed that the
		// only place a newline can be missing is at the end of the file.
		if !strings.HasSuffix(part.Value, "\n") {
			diffLines = append(diffLines, "\n\\ No newline at end of file\n")
		}
	}
	return strings.Join(diffLines, "")
}
