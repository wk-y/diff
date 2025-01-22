package patching

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/wk-y/diff"
)

// Apply hunks to a sequence of parts.
func ApplyHunks(a []string, hunks []Hunk) ([]string, error) {
	// To avoid modifying the passed hunks, a new array is made with the
	// corrected hunks
	h := make([]Hunk, len(hunks))

	previousAdjustment := 0 // Used to change starting adjustment of hunks based on previous adjustment
	for i, hunk := range hunks {
		// Recompute line counts
		hunk.aLines = 0
		hunk.bLines = 0
		for _, part := range hunk.parts {
			switch part.Action {
			case diff.DiffIdentical:
				hunk.aLines++
				hunk.bLines++
			case diff.DiffAdded:
				hunk.bLines++
			case diff.DiffRemoved:
				hunk.aLines++
			}
		}

		// Find where the hunk matches
		expectedALines := make([]string, 0, hunk.aLines)
		for _, part := range hunk.parts {
			if part.Action == diff.DiffIdentical || part.Action == diff.DiffRemoved {
				expectedALines = append(expectedALines, part.Value)
			}
		}

		// Adjust the start of the hunk
		n := len(a)
		hunkPositionFound := false
		// adjustedStart is reduced by 1, so that lines can be treated as 0 indexed
		adjustedStart := hunk.aStart + previousAdjustment - 1
		for ; adjustedStart+hunk.aLines <= n; adjustedStart++ {
			if (adjustedStart+hunk.aLines <= n) &&
				(adjustedStart+hunk.aLines >= 0) &&
				reflect.DeepEqual(a[adjustedStart:adjustedStart+hunk.aLines], expectedALines) {
				hunkPositionFound = true
				break
			}
		}
		if !hunkPositionFound {
			for adjustedStart = hunk.aStart + previousAdjustment - 1; adjustedStart >= 0; adjustedStart-- {
				if (adjustedStart+hunk.aLines <= n) &&
					(adjustedStart+hunk.aLines >= 0) &&
					reflect.DeepEqual(a[adjustedStart:adjustedStart+hunk.aLines], expectedALines) {
					hunkPositionFound = true
					break
				}
			}
		}
		if !hunkPositionFound {
			return nil, fmt.Errorf("could not find location of hunk %v", i)
		}

		adjustedStart++ // switch to one indexed
		previousAdjustment = adjustedStart - hunk.aStart
		hunk.aStart = adjustedStart

		h[i] = hunk
	}

	// TODO: Reject overlapping hunks?

	// Ensure the hunks are in order from first to last.
	slices.SortFunc(h, func(a, b Hunk) int {
		return a.aStart - b.aStart
	})

	// Apply the hunks as we go
	b := []string{}
	for aln := 1; aln <= len(a); aln++ {
		for len(h) > 0 && h[0].aStart <= aln {
			for _, part := range h[0].parts {
				switch part.Action {
				case diff.DiffIdentical:
					aln++
					b = append(b, part.Value)
				case diff.DiffAdded:
					b = append(b, part.Value)
				case diff.DiffRemoved:
					aln++
				}
			}
			h = h[1:]
		}
		if aln <= len(a) {
			b = append(b, a[aln-1])
		}
	}

	return b, nil
}
