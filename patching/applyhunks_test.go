// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package patching

import (
	"reflect"
	"testing"

	"github.com/wk-y/diff"
	"github.com/wk-y/diff/internal/strutils"
	"github.com/wk-y/diff/patching/internal/testdata/test1"
)

func TestApplyHunks(t *testing.T) {
	a := strutils.SplitLines(test1.A)
	b := strutils.SplitLines(test1.B)
	d := diff.Diff(a, b)
	hunks := HunkDiff(d)

	const expectedHunks = 48
	if len(hunks) != expectedHunks { // arbitrary number
		t.Fatalf("Wrong number of hunks in test case (test error) %v!=%v", len(hunks), expectedHunks)
	}

	// Test re-alignment by offsetting the hunks
	for i := range hunks {
		hunks[i].aStart += (i*7)%10 - 5 // Arbitrary formula
	}

	reconstructedB, err := ApplyHunks(a, hunks)
	if err != nil {
		t.Errorf("Failed to apply hunks: %v", err)
	}

	if !reflect.DeepEqual(reconstructedB, b) {
		t.Error("Reconstructed file doesn't match!")
		bDiff := diff.Diff(b, reconstructedB)
		t.Error(DiffString(bDiff))
	}
}

// Test that a mismatched hunk fails to apply
func TestApplyHunksMismatch(t *testing.T) {
	hunk := Hunk{
		aStart: 1,
		bStart: 1,
		parts: []diff.DiffPart{
			{Action: diff.DiffIdentical, Value: "NONEXISTENT"},
			{Action: diff.DiffAdded, Value: "Some text"},
			{Action: diff.DiffIdentical, Value: "NONEXISTENT"},
		},
	}
	a := []string{"The", "quick", "brown", "fox"}
	_, err := ApplyHunks(a, []Hunk{hunk})

	if err == nil {
		t.Error("Hunk application was supposed to fail!")
	}
}
