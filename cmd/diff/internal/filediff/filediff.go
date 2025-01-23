package filediff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/wk-y/diff"
	"github.com/wk-y/diff/patching"
)

type FileDiff struct {
	originalName, modifiedName string
	originalInfo, modifiedInfo os.FileInfo
	diff                       []diff.DiffPart
}

func DiffFiles(aName, bName string) (FileDiff, error) {
	result := FileDiff{
		originalName: aName,
		modifiedName: bName,
	}
	aFile, err := os.Open(result.originalName)
	if err != nil {
		return result, err
	}
	defer aFile.Close()

	result.originalInfo, err = aFile.Stat()
	if err != nil {
		return result, err
	}

	aBytes, err := io.ReadAll(aFile)
	if err != nil {
		return result, err
	}

	bFile, err := os.Open(result.modifiedName)
	if err != nil {
		return result, err
	}
	defer bFile.Close()

	result.modifiedInfo, err = bFile.Stat()
	if err != nil {
		return result, err
	}

	bBytes, err := io.ReadAll(bFile)
	if err != nil {
		return result, err
	}

	result.diff = diff.LineDiff(string(aBytes), string(bBytes))

	return result, nil
}

func formatHeaderInfo(name string, info os.FileInfo) string {
	const headerDateFormat = "2006-01-02 15:04:05.000000000 -0700"
	if strings.ContainsRune(name, ' ') {
		name = fmt.Sprintf("\"%v\"", strings.ReplaceAll(name, "\"", "\\\""))
	}
	return fmt.Sprintf("%v\t%v\n", name, info.ModTime().Format(headerDateFormat))
}

func (f FileDiff) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("--- %v", formatHeaderInfo(f.originalName, f.originalInfo)))
	builder.WriteString(fmt.Sprintf("+++ %v", formatHeaderInfo(f.modifiedName, f.modifiedInfo)))
	builder.WriteString(patching.DiffString(f.diff))
	return builder.String()
}
