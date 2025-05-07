package utils

import (
	"cryphtron/internal/utils/unarchive"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
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
		}
		if isInInclude(h) {
			pb := pb.Default.Start64(h.FileInfo().Size())
			dstFile := filepath.Join(dst, h.Name)
			err = ioToFile(reader, dstFile, pb)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
