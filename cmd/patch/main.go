// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/wk-y/diff/internal/exitcodes"
	"github.com/wk-y/diff/internal/strutils"
	"github.com/wk-y/diff/patching"
)

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		return
	}

	originalFileName := flag.Arg(0)
	patchFileName := flag.Arg(1)

	originalBytes, err := os.ReadFile(originalFileName)
	if err != nil {
		fmt.Printf("Failed to read original file: %v\n", err)
		os.Exit(exitcodes.IoError)
	}

	patchBytes, err := os.ReadFile(patchFileName)
	if err != nil {
		fmt.Printf("Failed to read patch file: %v\n", err)
		os.Exit(exitcodes.IoError)
	}

	// Skip the first two lines of the diff
	// TODO: Actually parse the diff header
	patchLines := strutils.SplitLines(string(patchBytes))
	if len(patchLines) < 2 {
		fmt.Println("Patch file too short")
		os.Exit(1)
	}

	patchString := strings.Join(patchLines[2:], "")

	if err != nil {
		fmt.Printf("Failed to read patch file: %v\n", err)
		os.Exit(exitcodes.IoError)
	}

	a := strutils.SplitLines(string(originalBytes))

	hunks, err := patching.ParseHunks(patchString)
	if err != nil {
		fmt.Printf("Failed to parse patch: %v\n", err)
		os.Exit(1)
	}

	b, err := patching.ApplyHunks(a, hunks)
	if err != nil {
		fmt.Printf("Failed to apply patch: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(originalFileName, []byte(strings.Join(b, "")), 0o664)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}
}
