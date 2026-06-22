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
	if isNeedHelpAtFirst(args) {
		LogInfo(CmdName(CopyCmdID), ": Copies a File or Folder")
		LogInfo("usage was: " + CmdName(CopyCmdID) + " src dst")
		LogInfo("src is required for source name or path of file and folder")
		LogInfo("dst is required for source name or path of file and folder")
		return nil
	}
	if len(args) < 2 {
		return cmdError(
			CmdName(CopyCmdID) + " use -h or -help to get help")
	}
	src := NormalizeFilePath(args[0])
	statSrc, _ := StatFilePath(args[0])
	if statSrc == nil {
		return cmdError(CmdName(CopyCmdID) + ": src is not exists")
	}

	dst := NormalizeFilePath(args[1])
	statDst, _ := StatFilePath(dst)
	if statDst != nil &&
		statDst.Mode().Type() != statSrc.Mode().Type() {
		return cmdError(CmdName(CopyCmdID) + ": dst type mismatch with src")
	}

	LogInfo(CmdName(CopyCmdID), ":", "Copying", src, "to", dst)
	exclude := make([]string, 0)
	if len(args) > 2 && strings.HasPrefix(args[2], "-ex:") {
		exclude = append(exclude, strings.TrimPrefix(args[2], "-ex:"))
	}
	return CopyFile(statSrc, dst, exclude...)
}

func MvCmd(args []string) error {
	if isNeedHelpAtFirst(args) {
		LogInfo(CmdName(MoveCmdID), ": Moves a File or Folder")
		LogInfo("usage was: " + CmdName(MoveCmdID) + " src dst")
		LogInfo("src is required for source name or path of file and folder")
		LogInfo("dst is required for source name or path of file and folder")
		return nil
	}
	if len(args) < 2 {
		return cmdError(
			CmdName(MoveCmdID) + " use -h or -help to get help")
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
		return cmdError(CmdName(MoveCmdID) + ": dst type mismatch with src")
	}

	LogInfo(CmdName(MoveCmdID), ":", "Moving", src, "to", dst)
	if err := CopyFile(statSrc, dst); err != nil {
		return err
	}
	return RemoveFileAndFolders(src)
}

func RemoveFileCmd(args []string) error {
	if len(args) < 1 {
		return cmdError(CmdName(RemoveCmdID) + " use -h or -help to get help")
	}
	if isNeedHelpAtFirst(args) {
		LogInfo(CmdName(RemoveCmdID), ": Removes a File")
		LogInfo("usage was: " + CmdName(RemoveCmdID) + " file")
		LogInfo("file is required for file name or path")
		return nil
	}
	file := NormalizeFilePath(args[0])
	stat, _ := StatFilePath(file)
	if stat == nil {
		return cmdError(CmdName(RemoveCmdID) + ": file is not exists")
	}
	if stat.IsDir() {
		return cmdErrors(CmdName(RemoveCmdID) + ": cannot remove a folder")
	}
	LogInfo(CmdName(RemoveCmdID), ":", "Removing File", file)
	return RemoveFileAndFolders(file)
}
