package cmds

const MkdirCmdID = "md"
const RemoveDirCmdID = "rd"

func init() {
	AppendCmd(CmdName(MkdirCmdID), MkdirCmd)
	AppendCmd(CmdName(RemoveDirCmdID), RemoveDirCmd)
}

func MkdirCmd(args []string) error {
	if len(args) < 1 {
		return cmdError("i-md correct usage was: i-md dir [options]")
	}
	return MkDir(args[0])
}

func RemoveDirCmd(args []string) error {
	if len(args) < 1 {
		return cmdError("i-rd correct usage was: i-rd dir [options]")
	}
	return RemoveFileAndFolders(args[0])
}
