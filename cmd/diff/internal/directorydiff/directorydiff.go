// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package directorydiff

import (
	"os"
	"path"
	"sort"

	"github.com/wk-y/diff"
	"github.com/wk-y/diff/cmd/diff/internal/filediff"
)

type DiffMessage interface {
	Path() string
}

type diffMessage struct {
	path string
}

func (d diffMessage) Path() string {
	return d.path
}

type DiffMessageDeleted struct {
	diffMessage
}

type DiffMessageIdentical struct {
	diffMessage
}

type DiffMessageModified struct {
	diffMessage
	filediff.FileDiff
}

type DiffMessageAdded struct {
	diffMessage
}

type DiffMessageError struct {
	diffMessage
	Error error
}

type DiffMessageDifferentTypes struct {
	diffMessage
	AType, BType string
}

// DiffDirectories will write messages to ch.
func DiffDirectories(aPath, bPath string, callback func(DiffMessage)) {
	diffDirectories(aPath, bPath, "/", callback)
}

func diffDirectories(aPath, bPath, commonPath string, callback func(DiffMessage)) {
	aDir := path.Join(aPath, commonPath)
	bDir := path.Join(bPath, commonPath)

	aEntries, err := os.ReadDir(aDir)
	if err != nil {
		callback(DiffMessageError{
			diffMessage: diffMessage{
				path: path.Join(aPath),
			},
			Error: err,
		})
		return
	}
	sort.Slice(aEntries, func(i, j int) bool {
		return aEntries[i].Name() < aEntries[j].Name()
	})

	bEntries, err := os.ReadDir(bDir)
	if err != nil {
		callback(DiffMessageError{
			diffMessage: diffMessage{
				path: path.Join(bPath),
			},
			Error: err,
		})
		return
	}
	sort.Slice(bEntries, func(i, j int) bool {
		return bEntries[i].Name() < bEntries[j].Name()
	})

	d := diff.DiffAlgorithm(len(aEntries), len(bEntries), func(i, j int) bool {
		return aEntries[i].Name() == bEntries[j].Name()
	})

	var i, j int
	for _, action := range d {
		switch action {
		case diff.DiffIdentical:
			aIsDir := aEntries[i].IsDir()
			bIsDir := bEntries[j].IsDir()
			if aIsDir != bIsDir {
				callback(DiffMessageDifferentTypes{
					diffMessage: diffMessage{path: path.Join(commonPath, aEntries[i].Name())},
					AType:       fileType(aEntries[i]),
					BType:       fileType(bEntries[j]),
				})
			} else if aIsDir {
				diffDirectories(aPath, bPath, path.Join(commonPath, aEntries[i].Name()), callback)
			} else {
				callback(diffFiles(aPath, bPath, path.Join(commonPath, aEntries[i].Name())))
			}
			i++
			j++
		case diff.DiffRemoved:
			callback(DiffMessageDeleted{diffMessage: diffMessage{
				path: path.Join(commonPath, aEntries[i].Name()),
			}})
			i++
		case diff.DiffAdded:
			callback(DiffMessageAdded{diffMessage: diffMessage{
				path: path.Join(commonPath, bEntries[j].Name()),
			}})
			j++
		}
	}
}

func fileType(e os.DirEntry) string {
	if e.IsDir() {
		return "a directory"
	} else {
		return "a regular file"
	}
}

func sortDir(d []os.DirEntry) {
	sort.Slice(d, func(i, j int) bool {
		return d[i].Name() < d[j].Name()
	})
}

func recursiveListDir(root string) ([]string, error) {
	result := []string{}
	var search func(string) error
	search = func(dirname string) error {
		entries, err := os.ReadDir(path.Join(root, dirname))

		if err != nil {
			return err
		}

		sortDir(entries)

		for _, entry := range entries {
			entryPath := path.Join(dirname, entry.Name())
			if entry.IsDir() {
				if err := search(entryPath); err != nil {
					return err
				}
			} else {
				result = append(result, entryPath)
			}
		}
		return nil
	}
	err := search("/")
	return result, err
}

func diffFiles(aPath, bPath, relPath string) DiffMessage {
	aName := path.Join(aPath, relPath)
	bName := path.Join(bPath, relPath)

	fdiff, err := filediff.DiffFiles(aName, bName)
	if err != nil {
		return DiffMessageError{
			diffMessage: diffMessage{
				path: relPath,
			},
			Error: err,
		}
	}

	for _, part := range fdiff.Diff {
		if part.Action != diff.DiffIdentical {
			return DiffMessageModified{diffMessage: diffMessage{
				path: relPath,
			},
				FileDiff: fdiff,
			}
		}
	}

	return DiffMessageIdentical{diffMessage: diffMessage{
		path: relPath,
	},
	}
}
