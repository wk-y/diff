// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package filediff

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/wk-y/diff"
	"github.com/wk-y/diff/internal/strutils"
	"github.com/wk-y/diff/patching"
)

type FileDiff struct {
	OriginalInfo, ModifiedInfo fs.FileInfo
	Diff                       []diff.DiffPart
}

func DiffFiles(a, b fs.File) (FileDiff, error) {
	var err error
	var aLines, bLines []string

	result := FileDiff{}

	aLines, result.OriginalInfo, err = diffFilesHelper(a)
	if err != nil {
		return result, err
	}

	bLines, result.ModifiedInfo, err = diffFilesHelper(b)
	if err != nil {
		return result, err
	}

	result.Diff = diff.Diff(aLines, bLines)

	return result, nil
}

func diffFilesHelper(f fs.File) (lines []string, info fs.FileInfo, err error) {
	info, err = f.Stat()
	if err != nil {
		return
	}

	lines, err = strutils.ReadLines(f)
	return
}

func formatHeaderInfo(name string, info os.FileInfo) string {
	const headerDateFormat = "2006-01-02 15:04:05.000000000 -0700"
	if strings.ContainsRune(name, ' ') {
		name = fmt.Sprintf("\"%v\"", strings.ReplaceAll(name, "\"", "\\\""))
	}
	return fmt.Sprintf("%v\t%v\n", name, info.ModTime().Format(headerDateFormat))
}

func (f FileDiff) HeaderString(aPath, bPath string) string {
	return fmt.Sprint(
		fmt.Sprintf("--- %v", formatHeaderInfo(aPath, f.OriginalInfo)),
		fmt.Sprintf("+++ %v", formatHeaderInfo(bPath, f.ModifiedInfo)),
	)
}

func (f FileDiff) String() string {
	return patching.DiffString(f.Diff)
}
