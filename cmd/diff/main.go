// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Diffs two files line-by-line, and shows the diff roughly in unified format
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

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

	a := flag.Arg(0)
	b := flag.Arg(1)
	if recursive {
		callback := func(msg directorydiff.DiffMessage) {
			switch msg := msg.(type) {
			case directorydiff.DiffMessageAdded:
				parent, file := path.Split(path.Join(b, msg.Path()))
				fmt.Printf("Only in %v: %v\n", strings.TrimSuffix(parent, "/"), file)
			case directorydiff.DiffMessageDeleted:
				parent, file := path.Split(path.Join(a, msg.Path()))
				fmt.Printf("Only in %v: %v\n", strings.TrimSuffix(parent, "/"), file)
			case directorydiff.DiffMessageModified:
				// Is the file a binary? Uses the strategy of checking for null byte
				// https://www.gnu.org/software/diffutils/manual/html_node/Binary.html
				isBinary := false
				for _, line := range msg.FileDiff.Diff {
					for _, c := range []byte(line.Value) {
						if c == 0 {
							isBinary = true
							break
						}
					}
					if isBinary {
						break
					}
				}
				if isBinary {
					fmt.Printf("Binary files %v and %v differ\n", path.Join(a, msg.Path()), path.Join(b, msg.Path()))
				} else {
					fmt.Print(msg.FileDiff)
				}
			case directorydiff.DiffMessageDifferentTypes:
				fmt.Printf("File %v is %v while file %v is a %v\n", path.Join(a, msg.Path()), msg.AType, path.Join(b, msg.Path()), msg.BType)
			case directorydiff.DiffMessageError:
				fmt.Fprintf(os.Stderr, "%v\n", msg)
			}
		}
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get CWD: %v", err)
			os.Exit(1)
		}
		directorydiff.DiffDirectories(
			os.DirFS(path.Join(wd, a)),
			os.DirFS(path.Join(wd, b)),
			callback,
		)
	} else {
		fdiff, err := diffSingle(a, b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to calculate diff: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(fdiff.HeaderString(a, b))
		fmt.Print(fdiff)
	}
}

func diffSingle(a, b string) (filediff.FileDiff, error) {
	aFile, err := os.Open(a)
	if err != nil {
		return filediff.FileDiff{}, err
	}
	defer aFile.Close()

	bFile, err := os.Open(b)
	if err != nil {
		return filediff.FileDiff{}, err
	}
	defer bFile.Close()

	return filediff.DiffFiles(aFile, bFile)
}
