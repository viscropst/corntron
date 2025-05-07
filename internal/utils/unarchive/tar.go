package unarchive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/fs"
	"path/filepath"
)

func TarReader(src fs.File) (*tar.Reader, error) {
	st, err := src.Stat()
	if err != nil {
		return nil, err
	}
	if isTgz(st) {
		return tgzReader(src)
	}
	return tarFromReader(src)
}

func tarFromReader(src io.ReadCloser) (*tar.Reader, error) {
	return tar.NewReader(src), nil
}

func isTgz(src fs.FileInfo) bool {
	ext := filepath.Ext(src.Name())
	return ext == "gz" ||
		ext == "tgz"
}

func tgzReader(src fs.File) (*tar.Reader, error) {
	gzReader, err := gzip.NewReader(src)
	if err != nil {
		return nil, err
	}
	return tarFromReader(gzReader)
}
