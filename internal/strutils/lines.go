// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package strutils

import (
	"io"
	"strings"
)

// This follows a UNIX-style approach. Lines end with a \n, and the \n is
// preserved. This means that joining the strings with "" will reproduce the
// original string. If the original string did not end in a \n, the last line
// in the split won't either.
func SplitLines(a string) []string {
	lines, _ := ReadLines(strings.NewReader(a))
	return lines
}

func ReadLines(r io.Reader) ([]string, error) {
	lines := []string{}
	stringBuffer := []byte{}
	var readBuffer [4096]byte // Buffer to store read bytes
	for {
		len, err := r.Read(readBuffer[:])
		for _, c := range readBuffer[:len] {
			stringBuffer = append(stringBuffer, c)
			if c == '\n' {
				lines = append(lines, string(stringBuffer))
				stringBuffer = stringBuffer[:0]
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return lines, err
		}
	}

	if len(stringBuffer) > 0 {
		lines = append(lines, string(stringBuffer))
	}

	return lines, nil
}
