// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Diffs two files line-by-line, and shows the diff roughly in unified format
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/wk-y/diff"
	"github.com/wk-y/diff/internal/exitcodes"
	"github.com/wk-y/diff/patching"
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

// formatHeaderInfo formats information about a file into a string matching the
// format used in diff's header lines.
// Name is passed separately, as info.Name does not provide the full path.
func formatHeaderInfo(name string, info os.FileInfo) string {
	const headerDateFormat = "2006-01-02 15:04:05.000000000 -0700"
	if strings.ContainsRune(name, ' ') {
		name = fmt.Sprintf("\"%v\"", strings.ReplaceAll(name, "\"", "\\\""))
	}
	return fmt.Sprintf("%v\t%v\n", name, info.ModTime().Format(headerDateFormat))
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %v file1 file2\n", os.Args[0])
		os.Exit(exitcodes.UsageError)
	}

	a, err := getFileData(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read %v: %v", os.Args[1], err)
		os.Exit(exitcodes.IoError)
	}

	b, err := getFileData(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read %v: %v", os.Args[2], err)
		os.Exit(exitcodes.IoError)
	}

	diffed := diff.LineDiff(string(a.Contents), string(b.Contents))
	fmt.Printf("--- %v", formatHeaderInfo(a.Name, a.Info))
	fmt.Printf("+++ %v", formatHeaderInfo(b.Name, b.Info))
	fmt.Print(patching.DiffString(diffed))
}
