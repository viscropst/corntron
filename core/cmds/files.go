package cmds

import (
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
		return cmdError(
			"i-cp correct usage was: i-cp src dst [options]")
	}
	src := NormalizeFilePath(args[0])
	statSrc, _ := StatFilePath(args[0])
	if statSrc == nil {
		return cmdError("i-cp: src is not exists")
	}

	dst := NormalizeFilePath(args[1])
	statDst, _ := StatFilePath(dst)
	if statDst != nil &&
		statDst.Mode().Type() != statSrc.Mode().Type() {
		return cmdError("i-cp: dst type mismatch with src")
	}

	LogInfo("i-cp", ":", "Copying", src, "to", dst)
	exclude := make([]string, 0)
	if len(args) > 2 && strings.HasPrefix(args[2], "-ex:") {
		exclude = append(exclude, strings.TrimPrefix(args[2], "-ex:"))
	}
	return CopyFile(statSrc, dst, exclude...)
}

func MvCmd(args []string) error {
	if len(args) < 2 {
		return cmdError(
			"i-mv correct usage was: i-mv src dst [options]")
	}
	src := NormalizeFilePath(args[0])
	statSrc, _ := StatFilePath(src)
	if statSrc == nil {
		return cmdError("i-mv: src is not exists")
	}

	dst := NormalizeFilePath(args[1])
	statDst, _ := StatFilePath(dst)
	if statDst != nil &&
		statDst.Mode().Type() != statSrc.Mode().Type() {
		return cmdError("i-mv: dst type mismatch with src")
	}

	LogInfo("i-mv", ":", "Moving", src, "to", dst)
	if err := CopyFile(statSrc, dst); err != nil {
		return err
	}
	return RemoveFileAndFolders(src)
}

func RemoveFileCmd(args []string) error {
	if len(args) < 1 {
		return cmdError("i-rf correct usage was: i-rf dir [options]")
	}
	file := NormalizeFilePath(args[0])
	LogInfo("i-rf", ":", "Removing File", file)
	return RemoveFileAndFolders(file)
}
