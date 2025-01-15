package strutils

import (
	"reflect"
	"testing"
)

func TestSplitLines(t *testing.T) {
	type testCase struct {
		text     string
		expected []string
	}

	tests := []testCase{
		{
			text:     "Hello\nworld",
			expected: []string{"Hello\n", "world"},
		},
		{
			text:     "Hello\nworld\n",
			expected: []string{"Hello\n", "world\n"},
		},
		{
			text:     "",
			expected: []string{},
		},
	}
	for _, test := range tests {
		result := SplitLines(test.text)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Expected %v, got %v", result, test.expected)
		}
	}
}
