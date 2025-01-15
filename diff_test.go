/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

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
			t.Fail()
		}
	}
}
