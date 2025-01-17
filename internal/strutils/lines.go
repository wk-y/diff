// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package strutils

// This follows a UNIX-style approach. Lines end with a \n, and the \n is
// preserved. This means that joining the strings with "" will reproduce the
// original string. If the original string did not end in a \n, the last line
// in the split won't either.
func SplitLines(a string) []string {
	split := make([]string, 0)
	buffer := make([]rune, 0)
	for _, r := range a {
		buffer = append(buffer, r)
		if r == '\n' {
			split = append(split, string(buffer))
			buffer = []rune{}
		}
	}
	if len(buffer) > 0 {
		split = append(split, string(buffer))
	}
	return split
}
