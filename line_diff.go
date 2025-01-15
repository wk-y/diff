/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package diff

import (
	"fmt"
	"strings"

	"github.com/wk-y/diff/internal/strutils"
)

func LineDiff(a, b string) []DiffPart {
	return Diff(strutils.SplitLines(a), strutils.SplitLines(b))
}

// DiffString formats an array of DiffParts into a unified diff.
// DiffString will produce strange results if d is not from LineDiff.
func DiffString(d []DiffPart) string {
	// Find line numbers
	aln := make([]int, len(d))
	bln := make([]int, len(d))

	ai := 0
	bi := 0
	for i, part := range d {
		switch part.Action {
		case DiffIdentical:
			ai++
			bi++
		case DiffAdded:
			bi++
		case DiffRemoved:
			ai++
		}
		aln[i] = ai
		bln[i] = bi
	}

	const contextLines = 3

	type Hunk struct {
		aStart, bStart int // Starting line number
		aLines, bLines int // Number of lines covered by the hunk
		dStart, dEnd   int // Indices of d covered by the hunk
	}

	hunks := make([]Hunk, 0)
	for i := 0; i < len(d); i++ {
		if d[i].Action != DiffIdentical {
			dStart := max(i-contextLines, 0)
			newHunk := Hunk{
				aStart: aln[dStart],
				bStart: bln[dStart],
				dStart: dStart,
			}

			distancePastEdit := 0
			j := i
			for j < len(d)-1 && distancePastEdit <= contextLines*2 {
				j++
				distancePastEdit++
				if d[j].Action != DiffIdentical {
					distancePastEdit = 0
				}
			}
			i = j

			newHunk.dEnd = j
			if distancePastEdit > contextLines {
				newHunk.dEnd -= distancePastEdit - contextLines
			}

			newHunk.aLines = aln[newHunk.dEnd] + 1 - newHunk.aStart
			newHunk.bLines = bln[newHunk.dEnd] + 1 - newHunk.bStart

			hunks = append(hunks, newHunk)
		}
	}

	diffLines := make([]string, 0)
	for _, hunk := range hunks {
		header := fmt.Sprintf("@@ -%v,%v +%v,%v @@\n", hunk.aStart, hunk.aLines, hunk.bStart, hunk.bLines)
		diffLines = append(diffLines, header)
		for i := hunk.dStart; i <= hunk.dEnd; i++ {
			switch d[i].Action {
			case DiffAdded:
				diffLines = append(diffLines, fmt.Sprint("+", d[i].Value))
			case DiffRemoved:
				diffLines = append(diffLines, fmt.Sprint("-", d[i].Value))
			case DiffIdentical:
				diffLines = append(diffLines, fmt.Sprint(" ", d[i].Value))
			}

			// Since it is assumed that LineDiff is used, it is assumed that the
			// only place a newline can be missing is at the end of the file.
			if len(d[i].Value) == 0 || d[i].Value[len(d[i].Value)-1] != '\n' {
				diffLines = append(diffLines, "\n\\ No newline at end of file\n")
			}
		}
	}

	return strings.Join(diffLines, "")
}
