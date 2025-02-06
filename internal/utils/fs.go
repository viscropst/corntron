package utils

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const defaultMod = 0666

func IOToFile(from io.Reader, to string) error {
	return ioToFile(from, NormalizePath(to))
}

func ioToFile(from io.Reader, to string) error {
	toFile, err := os.OpenFile(
		to, os.O_CREATE|os.O_RDWR, defaultMod)
	if err != nil {
		return err
	}
	_, err = io.Copy(toFile, from)
	_ = toFile.Close()
	return err
}

func CopyToFile(from os.FileInfo, to string, excludes ...string) error {
	return copyToFile(from, NormalizePath(to), excludes...)
}

func copyToFile(from os.FileInfo, to string, excludes ...string) error {
	if !from.IsDir() {
		fileSrc, _ := os.Open(from.Name())
		err := ioToFile(fileSrc, to)
		_ = fileSrc.Close()
		return err
	} else {
		if from == nil {
			_ = os.Mkdir(to, os.ModeDir|defaultMod)
		}
		return filepath.Walk(from.Name(),
			func(path string, i fs.FileInfo, err error) error {
				if path == to || i.Name() == "" {
					return nil
				}

				rel, _ := filepath.Rel(to, path)
				dstFull := filepath.Join(to, rel)
				if len(excludes) > 0 && rel == excludes[0] {
					return nil
				}
				if i.IsDir() {
					return os.Mkdir(dstFull, os.ModeDir|defaultMod)
				} else {
					return copyToFile(i, dstFull)
				}
			})
	}
}
