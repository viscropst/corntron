package utils

import (
	"cryphtron/internal/utils/unarchive"
	"io"
	"os"
	"path/filepath"
)

func UnTarFromFile(src *os.File, dst string, include ...string) error {
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
			dstFile := filepath.Join(dst, NormalizePath(h.Name))
			err = copyFromFile(reader, dstFile, h.FileInfo())
			if err != nil && err != io.EOF {
				return err
			}
		}
	}
	return nil
}
