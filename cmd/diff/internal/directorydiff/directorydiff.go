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

// DiffDirectories will write messages to ch.
func DiffDirectories(aPath, bPath string, callback func(DiffMessage)) error {
	aEntries, err := recursiveListDir(aPath)
	if err != nil {
		callback(DiffMessageError{
			diffMessage: diffMessage{
				path: path.Join(aPath),
			},
			Error: err,
		})
		return err
	}

	bEntries, err := recursiveListDir(bPath)
	if err != nil {
		callback(DiffMessageError{
			diffMessage: diffMessage{
				path: path.Join(bPath),
			},
			Error: err,
		})
		return err
	}

	d := diff.Diff(aEntries, bEntries)
	for _, part := range d {
		switch part.Action {
		case diff.DiffAdded:
			callback(DiffMessageAdded{diffMessage: diffMessage{
				path: part.Value,
			}})
		case diff.DiffRemoved:
			callback(DiffMessageDeleted{diffMessage: diffMessage{
				path: part.Value,
			}})
		case diff.DiffIdentical:
			callback(diffFiles(aPath, bPath, part.Value))
		}
	}

	return nil
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
