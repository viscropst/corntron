package utils

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

const defaultMod = 0666

func IOToFile(from io.Reader, to string, bar *pb.ProgressBar) error {
	return ioToFile(from, NormalizePath(to), defaultMod, bar)
}

func ioToFile(from io.Reader, to string, mod os.FileMode, bar *pb.ProgressBar) error {
	if mod == 0 {
		mod = defaultMod
	}
	if st, _ := os.Stat(filepath.Dir(to)); st == nil {
		_ = os.MkdirAll(filepath.Dir(to), os.ModeDir|defaultMod)
	}
	toFile, err := os.OpenFile(
		to, os.O_CREATE|os.O_RDWR|os.O_TRUNC, mod)
	if err != nil {
		return err
	}
	defer CloseFileAndFinishBar(toFile, bar)
	if bar != nil {
		barReader := bar.NewProxyReader(from)
		_, err = io.Copy(toFile, barReader)
	} else {
		_, err = io.Copy(toFile, from)
	}
	return err
}

func CopyToFile(from os.FileInfo, to string, excludes ...string) error {
	return copyToFile(from, NormalizePath(to), excludes...)
}

func copyToFile(from os.FileInfo, to string, excludes ...string) error {
	if !from.IsDir() {
		fileSrc, _ := os.Open(from.Name())
		bar := pb.Default.Start64(from.Size())
		return ioToFile(fileSrc, to, defaultMod, bar)
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

func CloseFileAndFinishBar(file io.Closer, bar *pb.ProgressBar) {
	if file != nil {
		_ = file.Close()
	}
	if bar != nil {
		bar.Finish()
	}
}

func copyFromFile(file io.Reader, to string, fileInfo os.FileInfo) error {
	if fileInfo == nil {
		fileInfo = file.(fs.FileInfo)
	}
	if err := os.MkdirAll(filepath.Dir(to), defaultMod); err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return os.MkdirAll(filepath.Dir(to), defaultMod)
	}
	filePb := pb.Default.Start64(fileInfo.Size())
	return ioToFile(file, to, fileInfo.Mode(), filePb)
}

func Mkdir(path string) error {
	return os.MkdirAll(path, defaultMod)
}

func RemoveDir(path string) error {
	return os.RemoveAll(path)
}
