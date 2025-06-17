package cmds

import (
	"corntron/internal/utils"
	"errors"
)

const MkdirCmdID = "md"
const RemoveDirCmdID = "rd"

func init() {
	AppendCmd(CmdName(MkdirCmdID), MkdirCmd)
	AppendCmd(CmdName(RemoveDirCmdID), RemoveDirCmd)
}

func MkdirCmd(args []string) error {
	if len(args) < 1 {
		return errors.New("i-md correct usage was: i-md dir [options]")
	}
	return utils.Mkdir(args[0])
}

func RemoveDirCmd(args []string) error {
	if len(args) < 1 {
		return errors.New("i-rd correct usage was: i-rd dir [options]")
	}
	return utils.Remove(args[0])
}
