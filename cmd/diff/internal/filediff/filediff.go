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
	OriginalName, ModifiedName string
	OriginalInfo, ModifiedInfo os.FileInfo
	Diff                       []diff.DiffPart
}

func DiffFiles(aName, bName string) (FileDiff, error) {
	result := FileDiff{
		OriginalName: aName,
		ModifiedName: bName,
	}
	aFile, err := os.Open(result.OriginalName)
	if err != nil {
		return result, err
	}
	defer aFile.Close()

	result.OriginalInfo, err = aFile.Stat()
	if err != nil {
		return result, err
	}

	aBytes, err := io.ReadAll(aFile)
	if err != nil {
		return result, err
	}

	bFile, err := os.Open(result.ModifiedName)
	if err != nil {
		return result, err
	}
	defer bFile.Close()

	result.ModifiedInfo, err = bFile.Stat()
	if err != nil {
		return result, err
	}

	bBytes, err := io.ReadAll(bFile)
	if err != nil {
		return result, err
	}

	result.Diff = diff.LineDiff(string(aBytes), string(bBytes))

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
	builder.WriteString(fmt.Sprintf("--- %v", formatHeaderInfo(f.OriginalName, f.OriginalInfo)))
	builder.WriteString(fmt.Sprintf("+++ %v", formatHeaderInfo(f.ModifiedName, f.ModifiedInfo)))
	builder.WriteString(patching.DiffString(f.Diff))
	return builder.String()
}
