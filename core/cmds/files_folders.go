package cmds

const MkdirCmdID = "md"
const RemoveDirCmdID = "rd"

func init() {
	AppendCmd(CmdName(MkdirCmdID), MkdirCmd)
	AppendCmd(CmdName(RemoveDirCmdID), RemoveDirCmd)
}

func MkdirCmd(args []string) error {
	if len(args) < 1 {
		return cmdErrors(CmdName(MkdirCmdID), " use -h or -help to get help")
	}
	if isNeedHelpAtFirst(args) {
		LogInfo(CmdName(MkdirCmdID), ": Makes a Folder")
		LogInfo("usage was: " + MkdirCmdID + " dir")
		LogInfo("dir is required for folder name or path")
		return nil
	}
	_, err := StatFilePath(args[0])
	if err == nil {
		LogInfo("the", args[0], "is already created skipping.")
		return nil
	}
	return MkDir(args[0])
}

func RemoveDirCmd(args []string) error {
	if len(args) < 1 {
		return cmdErrors(CmdName(RemoveDirCmdID), " use -h or -help to get help")
	}
	if isNeedHelpAtFirst(args) {
		LogInfo(CmdName(RemoveDirCmdID), ": Removes a Folder")
		LogInfo("usage was: " + RemoveDirCmdID + " dir")
		LogInfo("dir is required for folder name or path")
		return nil
	}
	stat, _ := StatFilePath(args[0])
	if stat == nil {
		return cmdError(CmdName(RemoveDirCmdID) + ": folder is not exists")
	}
	if !stat.IsDir() {
		return cmdErrors(CmdName(RemoveDirCmdID) + ": cannot remove a file")
	}
	return RemoveFileAndFolders(args[0])
}
