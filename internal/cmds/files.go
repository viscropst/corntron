package cmds

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const defaultMod = 0666

const CopyCmdID = "cp"
const MoveCmdID = "mv"

func init() {
	AppendCmd(CmdName(CopyCmdID), CpCmd)
	AppendCmd(CmdName(MoveCmdID), MvCmd)
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

func CpCmd(args []string) error {
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
	return copyToFile(statSrc, dst, exclude...)
}

func MvCmd(args []string) error {
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

	if err := copyToFile(statSrc, dst); err != nil {
		return err
	}
	if !statSrc.IsDir() {
		_ = os.Remove(src)
	} else {
		_ = os.RemoveAll(src)
	}
	return nil
}
