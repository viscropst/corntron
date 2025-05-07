package utils

import (
	"archive/zip"
	"cryphtron/internal/utils/unarchive"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

type inArchiveFileInfo interface {
	FileInfo() fs.FileInfo
}

func isInInclude(src inArchiveFileInfo, includes ...string) bool {
	if len(includes) == 0 {
		return true
	}
	for _, v := range includes {
		fileInfo := src.FileInfo()
		if !strings.HasPrefix(fileInfo.Name(), v) {
			return true
		}
	}
	return false
}

func filterZipFiles(src []*zip.File, includes ...string) []*zip.File {
	if len(includes) == 0 {
		return src
	}
	result := make([]*zip.File, 0)
	for _, v := range src {
		if !isInInclude(v, includes...) {
			result = append(result, v)
		}
	}
	return result
}

func UnZipFromFile(src *os.File, dst string, include ...string) error {
	reader, err := unarchive.ZipReader(src)
	if err != nil {
		return err
	}
	toUnarchive := filterZipFiles(reader.File, include...)
	pbTotal := pb.StartNew(len(toUnarchive))
	for _, file := range toUnarchive {
		tmp, err := file.Open()
		if err != nil {
			return err
		}
		dstFile := filepath.Join(dst, file.Name)
		pbTotal.Increment()
		filePb := pb.Default.Start64(file.FileInfo().Size())
		err = ioToFile(tmp, dstFile, filePb)
		if err != nil {
			return err
		}
	}
	pbTotal.Finish()
	return nil
}
