package directorydiff

import (
	"os"
	"path"
	"sort"

	"github.com/wk-y/diff"
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
	diff []diff.DiffPart
}

type DiffMessageAdded struct {
	diffMessage
}

type DiffMessageError struct {
	diffMessage
	Error error
}

// DiffDirectories will write messages to ch.
func DiffDirectories(aPath, bPath string) (result []DiffMessage) {
	aEntries, err := recursiveListDir(aPath)
	if err != nil {
		result = append(result, DiffMessageError{
			diffMessage: diffMessage{
				path: path.Join(aPath),
			},
			Error: err,
		})
		return
	}
	sort.Slice(aEntries, func(i, j int) bool {
		return aEntries[i] < aEntries[j]
	})

	bEntries, err := recursiveListDir(bPath)
	if err != nil {
		result = append(result, DiffMessageError{
			diffMessage: diffMessage{
				path: path.Join(bPath),
			},
			Error: err,
		})
		return
	}
	sort.Slice(bEntries, func(i, j int) bool {
		return bEntries[i] < bEntries[j]
	})

	d := diff.Diff(aEntries, bEntries)
	for _, part := range d {
		switch part.Action {
		case diff.DiffAdded:
			result = append(result, DiffMessageAdded{diffMessage: diffMessage{
				path: part.Value,
			}})
		case diff.DiffRemoved:
			result = append(result, DiffMessageDeleted{diffMessage: diffMessage{
				path: part.Value,
			}})
		case diff.DiffIdentical:
			// TODO: Implement DiffModified and identical detection
			result = append(result, DiffMessageModified{diffMessage: diffMessage{
				path: part.Value,
			}})
		}
	}
	return
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
