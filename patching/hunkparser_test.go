package patching

import (
	"reflect"
	"strings"
	"testing"

	"github.com/wk-y/diff"
	"github.com/wk-y/diff/patching/internal/testdata/test1"
)

// Tests that hunks can be parsed back from a string.
func TestHunkParsing(t *testing.T) {
	d := diff.LineDiff(test1.A, test1.B)
	hunks := HunkDiff(d)
	hunkStrings := make([]string, len(hunks))
	for i, hunk := range hunks {
		hunkStrings[i] = hunk.String()
	}
	diffString := strings.Join(hunkStrings, "")

	parsedHunks, err := ParseHunks(diffString)
	if err != nil {
		t.Errorf("Parsing hunks failed")
	}

	if !reflect.DeepEqual(parsedHunks, hunks) {
		t.Error("Parsed hunks do not match expected hunks")
	}
}
