package internal

import (
	"os"
	"path/filepath"
)

const PathSeparator = "/"
const PathListSeparator = ";"
const OSPathSeparator = string(os.PathSeparator)
const OSPathListSeparator = string(os.PathListSeparator)

func NormalizePath(src string) string {
	return filepath.Join(src)
}
