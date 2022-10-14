package internal

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Command func(args []string) error

var Commands = map[string]Command{
	"i-cp": cpCmd,
	"i-mv": mvCmd,
}

const defaultMod = 0666

func cpCmd(args []string) error {
	if len(args) < 2 {
		return errors.New(
			"i-cp correct usage was: i-cp src dst [options]")
	}
	src := filepath.Join(args[0])
	statSrc, _ := os.Stat(src)
	if statSrc == nil {
		return errors.New("i-cp: src is not exists")
	}

	dst := filepath.Join(args[1])
	statDst, _ := os.Stat(dst)
	if statDst != nil &&
		statDst.Mode().Type() != statSrc.Mode().Type() {
		return errors.New("i-cp: dst type mismatch with src")
	}

	exclude := make([]string, 0)
	if len(args) > 2 && strings.HasPrefix(args[2], "-ex:") {
		exclude = append(exclude, strings.TrimPrefix(args[2], "-ex:"))
	}

	if !statSrc.IsDir() {
		fileSrc, _ := os.Open(src)
		fileDst, _ := os.OpenFile(
			dst, os.O_CREATE|os.O_RDWR, statSrc.Mode())
		_, err := io.Copy(fileDst, fileSrc)
		_ = fileDst.Close()
		_ = fileSrc.Close()
		return err
	} else {
		if statDst == nil {
			_ = os.Mkdir(dst, os.ModeDir|defaultMod)
		}
		return filepath.Walk(src,
			func(path string, i fs.FileInfo, err error) error {
				if path == src || i.Name() == "" {
					return nil
				}

				rel, _ := filepath.Rel(src, path)
				dstFull := filepath.Join(dst, rel)
				if len(exclude) > 0 && rel == exclude[0] {
					return nil
				}
				if i.IsDir() {
					return os.Mkdir(dstFull, os.ModeDir|defaultMod)
				} else {
					return cpCmd([]string{path, dstFull})
				}
			})
	}
}

func mvCmd(args []string) error {
	if len(args) < 2 {
		return errors.New(
			"i-mv correct usage was: i-mv src dst [options]")
	}
	src := filepath.Join(args[0])
	statSrc, _ := os.Stat(src)
	if statSrc == nil {
		return errors.New("i-mv: src is not exists")
	}

	dst := filepath.Join(args[1])
	statDst, _ := os.Stat(dst)
	if statDst != nil &&
		statDst.Mode().Type() != statSrc.Mode().Type() {
		return errors.New("i-mv: dst type mismatch with src")
	}

	if !statSrc.IsDir() {
		fileSrc, _ := os.Open(src)
		fileDst, _ := os.OpenFile(
			dst, os.O_CREATE|os.O_RDWR, statSrc.Mode())
		_, err := fileDst.ReadFrom(fileSrc)
		_ = fileDst.Close()
		_ = fileSrc.Close()
		_ = os.Remove(src)
		return err
	} else {
		if statDst == nil {
			_ = os.Mkdir(dst, os.ModeDir|defaultMod)
		}
		filepath.Walk(src,
			func(path string, i fs.FileInfo, err error) error {
				if path == src || i.Name() == "" {
					return nil
				}

				rel, _ := filepath.Rel(src, path)
				dstFull := filepath.Join(dst, rel)
				if i.IsDir() {
					return os.Mkdir(dstFull, os.ModeDir|defaultMod)
				} else {
					return mvCmd([]string{path, dstFull})
				}
			})
		_ = os.RemoveAll(src)
	}

	return nil
}
