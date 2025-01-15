/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package diff

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	type TestCase struct {
		name     string
		a, b     []string
		expected []DiffPart
	}

	testCases := []TestCase{
		{
			"hello_world",
			[]string{"Hello", "world!"},
			[]string{"world!"},
			[]DiffPart{
				{DiffRemoved, "Hello"},
				{DiffIdentical, "world!"},
			},
		},
		{
			"longer",
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
	}

	for _, testCase := range testCases {
		result := Diff(testCase.a, testCase.b)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, Got %v", testCase.expected, result)
			t.Fail()
		}
	}
}
