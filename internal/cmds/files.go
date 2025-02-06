package cmds

import (
	"cryphtron/internal/utils"
	"errors"
	"os"
	"strings"
)

const CopyCmdID = "cp"
const MoveCmdID = "mv"

func init() {
	AppendCmd(CmdName(CopyCmdID), CpCmd)
	AppendCmd(CmdName(MoveCmdID), MvCmd)
}

func CpCmd(args []string) error {
	if len(args) < 2 {
		return errors.New(
			"i-cp correct usage was: i-cp src dst [options]")
	}
	src := utils.NormalizePath(args[0])
	statSrc, _ := os.Stat(src)
	if statSrc == nil {
		return errors.New("i-cp: src is not exists")
	}

	dst := utils.NormalizePath(args[1])
	statDst, _ := os.Stat(dst)
	if statDst != nil &&
		statDst.Mode().Type() != statSrc.Mode().Type() {
		return errors.New("i-cp: dst type mismatch with src")
	}

	exclude := make([]string, 0)
	if len(args) > 2 && strings.HasPrefix(args[2], "-ex:") {
		exclude = append(exclude, strings.TrimPrefix(args[2], "-ex:"))
	}
	return utils.CopyToFile(statSrc, dst, exclude...)
}

func MvCmd(args []string) error {
	if len(args) < 2 {
		return errors.New(
			"i-mv correct usage was: i-mv src dst [options]")
	}
	src := utils.NormalizePath(args[0])
	statSrc, _ := os.Stat(src)
	if statSrc == nil {
		return errors.New("i-mv: src is not exists")
	}

	dst := utils.NormalizePath(args[1])
	statDst, _ := os.Stat(dst)
	if statDst != nil &&
		statDst.Mode().Type() != statSrc.Mode().Type() {
		return errors.New("i-mv: dst type mismatch with src")
	}

	if err := utils.CopyToFile(statSrc, dst); err != nil {
		return err
	}
	if !statSrc.IsDir() {
		_ = os.Remove(src)
	} else {
		_ = os.RemoveAll(src)
	}
	return nil
}
