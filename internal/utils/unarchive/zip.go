package unarchive

import (
	"archive/zip"
	"os"
)

func ZipReader(src *os.File) (*zip.Reader, error) {
	stat, err := src.Stat()
	if err != nil {
		return nil, err
	}
	return zip.NewReader(src, stat.Size())
}
