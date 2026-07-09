package unarchive

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
)

func TarReader(src fs.File) (*tar.Reader, error) {
	st, err := src.Stat()
	if err != nil {
		return nil, err
	}
	if isTgz(st) {
		return tgzReader(src)
	}
	if isXz(st) {
		return xzReader(src)
	}
	if isBz2(st) {
		return bz2Reader(src)
	}
	if isZst(st) {
		return zstReader(src)
	}
	return tarFromReader(src)
}

func tarFromReader(src io.Reader) (*tar.Reader, error) {
	return tar.NewReader(src), nil
}

func isTgz(src fs.FileInfo) bool {
	ext := filepath.Ext(src.Name())
	return ext == ".gz" ||
		ext == ".tgz"
}

func tgzReader(src fs.File) (*tar.Reader, error) {
	gzReader, err := gzip.NewReader(src)
	if err != nil {
		return nil, err
	}
	return tarFromReader(gzReader)
}

func isXz(src fs.FileInfo) bool {
	ext := filepath.Ext(src.Name())
	return ext == ".xz"
}

func xzReader(src fs.File) (*tar.Reader, error) {
	xzReader, err := xz.NewReader(src)
	if err != nil {
		return nil, err
	}
	return tarFromReader(xzReader)
}

func isBz2(src fs.FileInfo) bool {
	ext := filepath.Ext(src.Name())
	return ext == ".bz2"
}

func bz2Reader(src fs.File) (*tar.Reader, error) {
	bz2Reader := bzip2.NewReader(src)
	return tarFromReader(bz2Reader)
}

func isZst(src fs.FileInfo) bool {
	ext := filepath.Ext(src.Name())
	return ext == ".zst" || ext == ".zstd" || ext == ".zstandard"
}

func zstReader(src fs.File) (*tar.Reader, error) {
	dec, err := zstd.NewReader(src)
	if err != nil {
		return nil, err
	}
	return tarFromReader(dec)
}
