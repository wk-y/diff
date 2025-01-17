// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package patching

import "github.com/wk-y/diff"

type Hunk struct {
	aStart, bStart int             // Starting line number
	aLines, bLines int             // Number of lines covered by the hunk
	parts          []diff.DiffPart // The lines in the diff
}

func HunkDiff(d []diff.DiffPart) []Hunk {
	// Find line numbers
	aln := make([]int, len(d))
	bln := make([]int, len(d))

	ai := 0
	bi := 0
	for i, part := range d {
		switch part.Action {
		case diff.DiffIdentical:
			ai++
			bi++
		case diff.DiffAdded:
			bi++
		case diff.DiffRemoved:
			ai++
		}
		aln[i] = ai
		bln[i] = bi
	}

	const contextLines = 3

	hunks := make([]Hunk, 0)
	for i := 0; i < len(d); i++ {
		if d[i].Action != diff.DiffIdentical {
			dStart := max(i-contextLines, 0)
			newHunk := Hunk{
				aStart: aln[dStart],
				bStart: bln[dStart],
			}

			distancePastEdit := 0
			j := i
			for j < len(d)-1 && distancePastEdit <= contextLines*2 {
				j++
				distancePastEdit++
				if d[j].Action != diff.DiffIdentical {
					distancePastEdit = 0
				}
			}
			i = j

			dEnd := j
			if distancePastEdit > contextLines {
				dEnd -= distancePastEdit - contextLines
			}

			newHunk.aLines = aln[dEnd] + 1 - newHunk.aStart
			newHunk.bLines = bln[dEnd] + 1 - newHunk.bStart
			newHunk.parts = d[dStart : dEnd+1]

			hunks = append(hunks, newHunk)
		}
	}

	return hunks
}
