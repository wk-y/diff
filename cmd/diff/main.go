// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Diffs two files line-by-line, and shows the diff roughly in unified format
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wk-y/diff/cmd/diff/internal/directorydiff"
	"github.com/wk-y/diff/cmd/diff/internal/filediff"
)

var recursive bool

func init() {
	flag.BoolVar(&recursive, "r", false, "Recurse")
}

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	if recursive {
		dirDiff := directorydiff.DiffDirectories(flag.Arg(0), flag.Arg(1))
		for _, msg := range dirDiff {
			fmt.Printf("%#v\n", msg)
		}
	} else {
		fdiff, err := filediff.DiffFiles(flag.Arg(0), flag.Arg(1))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to calculate diff: %v", err)
			os.Exit(1)
		}
		fmt.Print(fdiff)
	}
}
