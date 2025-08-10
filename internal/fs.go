package internal

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

const defaultMod = 0666
const fsAppend = os.O_APPEND | os.O_RDWR

func FSAppendFlag() int {
	return fsAppend
}

func IOToFile(from io.Reader, to string, bar *pb.ProgressBar, flags ...int) error {
	return ioToFile(from, NormalizePath(to), defaultMod, bar, flags...)
}

func flagIsEqual(a, b int) bool {
	if a == b {
		return true
	}
	return a&b != 0
}

func ioToFile(from io.Reader, to string, mod os.FileMode, bar *pb.ProgressBar, flags ...int) error {
	if mod == 0 {
		mod = defaultMod
	}
	flag := os.O_CREATE | os.O_RDWR | os.O_TRUNC
	if len(flags) > 0 {
		flag = flags[0]
	}
	if st, _ := StatPath(filepath.Dir(to)); st == nil && flagIsEqual(flag, os.O_CREATE) {
		_ = os.MkdirAll(filepath.Dir(to), os.ModeDir|defaultMod)
	}
	toFile, err := os.OpenFile(
		to, flag, mod)
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
	if from.IsDir() {
		if from == nil {
			_ = os.Mkdir(to, os.ModeDir|defaultMod)
		}
		return filepath.Walk(from.Name(),
			func(path string, i fs.FileInfo, err error) error {
				if path == to || i.Name() == "" {
					return nil
				}

				rel, _ := filepath.Rel(from.Name(), path)
				dstFull := filepath.Join(to, rel)
				if len(excludes) > 0 && rel == excludes[0] {
					return nil
				}
				if stat, _ := StatPath(dstFull); stat != nil && stat.IsDir() {
					return nil
				}
				if i.IsDir() {
					return os.Mkdir(dstFull, os.ModeDir|defaultMod)
				} else {
					stat, _ := StatPath(path)
					return copyToFile(stat, dstFull)
				}
			})
	}
	fileSrc, _ := os.Open(from.Name())
	defer CloseFileAndFinishBar(fileSrc, nil)
	bar := pb.Default.Start64(from.Size())
	err := ioToFile(fileSrc, to, defaultMod, bar)
	if err != nil {
		return err
	}
	return nil
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
	return os.MkdirAll(NormalizePath(path), defaultMod)
}

func Remove(path string) error {
	statSrc, _ := StatPath(path)
	if !statSrc.IsDir() {
		return os.Remove(statSrc.Name())
	} else {
		return os.RemoveAll(statSrc.Name())
	}
}

func StatPath(path string) (*statInfo, error) {
	result, err := os.Stat(NormalizePath(path))
	if err != nil {
		return nil, err
	}
	return toStatInfo(result, NormalizePath(path)), nil
}

func IfFolderNotExists(path string) bool {
	stat, _ := StatPath(path)
	if stat == nil {
		return true
	}
	return !stat.IsDir()
}

type statInfo struct {
	os.FileInfo
	name string
}

func (i statInfo) Name() string {
	if len(i.name) == 0 {
		return i.FileInfo.Name()
	}
	return i.name
}

func toStatInfo(stat os.FileInfo, path string) *statInfo {
	result := statInfo{
		FileInfo: stat,
		name:     path,
	}
	return &result
}

func (i statInfo) Open(flag ...int) (*os.File, error) {
	if len(flag) == 0 {
		return os.Open(i.Name())
	}
	return os.OpenFile(i.Name(), flag[0], i.Mode().Perm())
}

func Open(src string) (*os.File, error) {
	stat, err := StatPath(src)
	if err != nil {
		return nil, err
	}
	return stat.Open()
}
