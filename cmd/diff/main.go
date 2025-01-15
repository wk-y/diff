/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

// Diffs two files line-by-line, and shows the diff roughly in unified format
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/wk-y/diff/diff"
)

// Utility function to read f until EOF as a string
func slurp(f *os.File) (string, error) {
	builder := strings.Builder{}
	buffer := make([]byte, 4096) // Arbitrary size
	for {
		n, err := f.Read(buffer)
		builder.Write(buffer[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return builder.String(), err
		}
	}
	return builder.String(), nil
}

type fileData struct {
	Name     string
	Contents string
	Info     os.FileInfo
}

// Utility function to get all the data needed for a diff
func getFileData(name string) (fileData, error) {
	result := fileData{Name: name}
	f, err := os.Open(name)
	if err != nil {
		return result, err
	}

	defer f.Close()

	result.Contents, err = slurp(f)
	if err != nil {
		return result, err
	}

	result.Info, err = f.Stat()
	return result, err
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %v file1 file2\n", os.Args[0])
		return
	}

	a, err := getFileData(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to read %v: %v", os.Args[1], err)
	}

	b, err := getFileData(os.Args[2])
	if err != nil {
		fmt.Printf("Failed to read %v: %v", os.Args[2], err)
	}

	diffed := diff.LineDiff(string(a.Contents), string(b.Contents))
	const headerDateFormat = "2006-01-02 15:04:05.999999999 -0700"
	fmt.Printf("--- %v\t%v\n", a.Name, a.Info.ModTime().Format(headerDateFormat))
	fmt.Printf("+++ %v\t%v\n", b.Name, b.Info.ModTime().Format(headerDateFormat))
	fmt.Print(diff.DiffString(diffed))
}
