package internal

import (
	"corntron/internal/unarchive"
	"io"
	"os"
	"path/filepath"
)

func UnTarFromFile(src *os.File, dst string, include ...string) error {
	return UnTarFromFileWithBaseDir(src, dst, "", include...)
}

func UnTarFromFileWithBaseDir(src *os.File, dst string, baseDir string, include ...string) error {
	reader, err := unarchive.TarReader(src)
	if err != nil {
		return err
	}
	for {
		h, err := reader.Next()
		if err == io.EOF {
			break
		} else if err != nil && err != io.EOF {
			return err
		}
		if unarchive.IsInInclude(h, include...) {
			fileName := unarchive.FileNameInArchive(h.Name, baseDir)
			dstFile := filepath.Join(dst, fileName)
			err = copyFromFile(reader, dstFile, h.FileInfo())
			if err != nil && err != io.EOF {
				return err
			}
		}
	}
	return nil
}
