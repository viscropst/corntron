package unarchive

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type inArchiveFileInfo interface {
	FileInfo() fs.FileInfo
}

func IsInInclude(src inArchiveFileInfo, includes ...string) bool {
	if len(includes) == 0 {
		return true
	}
	if src == nil {
		return false
	}
	for _, v := range includes {
		fileInfo := src.FileInfo()
		if strings.HasPrefix(fileInfo.Name(), v) {
			return true
		}
	}
	return false
}

func FileNameInArchive(fileName, baseDir string) string {
	result := filepath.Join(fileName)
	hasPrefix := strings.HasPrefix(result, baseDir)
	if len(baseDir) > 0 && hasPrefix {
		result = strings.TrimPrefix(fileName, baseDir)
		result = filepath.Clean(result)
	}
	return result
}
