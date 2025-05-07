package utils

import (
	"cryphtron/internal/utils/unarchive"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

func UnZipFromFile(src *os.File, dst string, include ...string) error {
	reader, err := unarchive.ZipReader(src)
	if err != nil {
		return err
	}
	toUnarchive := unarchive.FilterZipFiles(reader.File, include...)
	pbTotal := pb.StartNew(len(toUnarchive))
	for _, file := range toUnarchive {
		tmp, err := file.Open()
		if err != nil {
			return err
		}
		dstFile := filepath.Join(dst, NormalizePath(file.Name))
		err = copyFromFile(tmp, dstFile, file.FileInfo())
		if err != nil {
			return err
		}
		pbTotal.Increment()
		CloseFileAndFinishBar(tmp, nil)
	}
	pbTotal.Finish()
	return nil
}
