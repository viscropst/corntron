package cmds

const InternalCmdPrefix = "i-"

func CmdName(name string) string {
	return InternalCmdPrefix + name
}

type Command func(args []string) error

var Commnads map[string]Command

func AppendCmd(name string, cmd Command) {
	if Commnads == nil {
		Commnads = make(map[string]Command)
	}
	Commnads[name] = cmd
}
