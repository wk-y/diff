// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package diff

import (
	"reflect"
	"strings"
	"testing"
)

func TestDiff(t *testing.T) {
	type TestCase struct {
		a, b     []string
		expected []DiffPart
	}

	testCases := []TestCase{
		{
			[]string{"Hello", "world!"},
			[]string{"world!"},
			[]DiffPart{
				{DiffRemoved, "Hello"},
				{DiffIdentical, "world!"},
			},
		},
		{
			[]string{"A", "B", "C", "C" /**/},
			[]string{"A" /**/, "C", "C", "B"},
			[]DiffPart{
				{DiffIdentical, "A"},
				{DiffRemoved, "B"},
				{DiffIdentical, "C"},
				{DiffIdentical, "C"},
				{DiffAdded, "B"},
			},
		},
		{
			strings.Split("The quick brown fox jumps over the lazy dog", " "),
			strings.Split("A over quick red fox jumps the lazy dog", " "),
			[]DiffPart{
				{DiffRemoved, "The"},
				{DiffAdded, "A"},
				{DiffAdded, "over"},
				{DiffIdentical, "quick"},
				{DiffRemoved, "brown"},
				{DiffAdded, "red"},
				{DiffIdentical, "fox"},
				{DiffIdentical, "jumps"},
				{DiffRemoved, "over"},
				{DiffIdentical, "the"},
				{DiffIdentical, "lazy"},
				{DiffIdentical, "dog"},
			},
		},
	}

	for _, testCase := range testCases {
		result := Diff(testCase.a, testCase.b)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, Got %v", testCase.expected, result)
		}

		// Verify that the original files can be reconstructed from result
		a, b := extractOriginals(result)
		if !reflect.DeepEqual(a, testCase.a) {
			t.Errorf("Original a was %v, got %v", testCase.a, a)
		}
		if !reflect.DeepEqual(b, testCase.b) {
			t.Errorf("Original b was %v, got %v", testCase.b, b)
		}
	}
}

func TestDiffActionString(t *testing.T) {
	testCases := []struct {
		action         DiffAction
		expectedString string
	}{
		{DiffAdded, "Added"},
		{DiffRemoved, "Removed"},
		{DiffIdentical, "Identical"},
	}
	for _, testCase := range testCases {
		if s := testCase.action.String(); s != testCase.expectedString {
			t.Errorf("Expected %v, got %v", testCase.expectedString, s)
		}
	}
}

// Given a diff, extract the original two files.
func extractOriginals(d []DiffPart) (a, b []string) {
	for _, part := range d {
		if part.Action != DiffAdded {
			a = append(a, part.Value)
		}
		if part.Action != DiffRemoved {
			b = append(b, part.Value)
		}
	}
	return
}
