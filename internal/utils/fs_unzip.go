package utils

import (
	"cryphtron/internal/utils/unarchive"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

func UnZipFromFile(src *os.File, dst string, include ...string) error {
	return UnZipFromFileWithBaseDir(src, dst, "", include...)
}
func UnZipFromFileWithBaseDir(src *os.File, dst string, baseDir string, include ...string) error {
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
		fileName := unarchive.FileNameInArchive(file.Name, baseDir)
		dstFile := filepath.Join(dst, fileName)
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
