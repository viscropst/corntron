package cmds

import (
	"corntron/internal"
	"corntron/internal/log"
	"strings"
)

const CopyCmdID = "cp"
const MoveCmdID = "mv"
const RemoveCmdID = "rf"

func init() {
	AppendCmd(CmdName(CopyCmdID), CpCmd)
	AppendCmd(CmdName(MoveCmdID), MvCmd)
	AppendCmd(CmdName(RemoveCmdID), RemoveFileCmd)
}

func CpCmd(args []string) error {
	if len(args) < 2 {
		return internal.Error(
			"i-cp correct usage was: i-cp src dst [options]")
	}
	src := internal.NormalizePath(args[0])
	statSrc, _ := internal.StatPath(args[0])
	if statSrc == nil {
		return internal.Error("i-cp: src is not exists")
	}

	dst := internal.NormalizePath(args[1])
	statDst, _ := internal.StatPath(dst)
	if statDst != nil &&
		statDst.Mode().Type() != statSrc.Mode().Type() {
		return internal.Error("i-cp: dst type mismatch with src")
	}

	internal.LogCLI(log.InfoLevel).Println("i-cp", ":", "Copying", src, "to", dst)
	exclude := make([]string, 0)
	if len(args) > 2 && strings.HasPrefix(args[2], "-ex:") {
		exclude = append(exclude, strings.TrimPrefix(args[2], "-ex:"))
	}
	return internal.CopyToFile(statSrc, dst, exclude...)
}

func MvCmd(args []string) error {
	if len(args) < 2 {
		return internal.Error(
			"i-mv correct usage was: i-mv src dst [options]")
	}
	src := internal.NormalizePath(args[0])
	statSrc, _ := internal.StatPath(src)
	if statSrc == nil {
		return internal.Error("i-mv: src is not exists")
	}

	dst := internal.NormalizePath(args[1])
	statDst, _ := internal.StatPath(dst)
	if statDst != nil &&
		statDst.Mode().Type() != statSrc.Mode().Type() {
		return internal.Error("i-mv: dst type mismatch with src")
	}

	internal.LogCLI(log.InfoLevel).Println("i-mv", ":", "Moving", src, "to", dst)
	if err := internal.CopyToFile(statSrc, dst); err != nil {
		return err
	}
	return internal.Remove(src)
}

func RemoveFileCmd(args []string) error {
	if len(args) < 1 {
		return internal.Error("i-rf correct usage was: i-rf dir [options]")
	}
	file := internal.NormalizePath(args[0])
	internal.LogCLI(log.InfoLevel).Println("i-rf", ":", "Removing File", file)
	return internal.Remove(file)
}
