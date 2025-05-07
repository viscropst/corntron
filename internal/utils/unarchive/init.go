package unarchive

import (
	"io/fs"
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
