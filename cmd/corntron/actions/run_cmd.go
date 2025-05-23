package actions

import (
	"corntron"
	cmdcontron "corntron/cmd/corntron"
	ct_core "corntron/core"
	"corntron/internal/utils/log"
	"errors"
	"os"
	"runtime"
	"strings"
)

type runCmd struct {
	cmdcontron.BaseAction
	Execute     string
	ExecArgs    []string
	withWaiting bool
}

func init() {
	appendAction(&runCmd{})
}

func (c *runCmd) ActionName() string {
	return "run-cmd"
}

func (c *runCmd) ParseArg(info cmdcontron.FlagInfo) error {
	argCmdIdx := info.Index + 1
	if len(os.Args) > argCmdIdx && len(info.Args[argCmdIdx]) > 0 {
		c.Execute = info.Args[argCmdIdx]
	} else {
		switch runtime.GOOS {
		case "windows":
			c.Execute = os.Getenv("COMSPEC")
		case "linux", "macos", "openbsd", "freebsd", "ios", "android":
			c.Execute = os.Getenv("SHELL")
			if len(c.Execute) == 0 {
				c.Execute = "/bin/sh"
			}
		default:
		}
		cmdcontron.CliLog(log.WarnLevel).Println("no command to exec,will use default shell or cmd")
	}

	if len(c.Execute) == 0 {
		errBuilder := strings.Builder{}
		errBuilder.WriteString("exec argument invalid: usage ")
		errBuilder.WriteString(info.CmdName + " ")
		errBuilder.WriteString(info.Name + " ")
		errBuilder.WriteString("<command>")
		return errors.New(errBuilder.String())
	}

	if len(os.Args) > argCmdIdx+1 {
		c.ExecArgs = info.Args[argCmdIdx+1:]
	}

	return nil
}

func (c *runCmd) BeforeCore(coreConfig *ct_core.MainConfig) error {
	return nil
}

func (c *runCmd) BeforeLoad(flag *cmdcontron.CmdFlag) (string, []string) {
	c.withWaiting = !flag.NoWaiting
	return c.BaseAction.BeforeLoad(flag)
}

func (c *runCmd) Exec(core *corntron.Core) error {
	var err error

	err = core.ProcessRtBootstrap(false)
	if err != nil {
		newErr := errors.New("error while bootstrapping:" + err.Error())
		return newErr
	}

	err = core.ProcessRtMirror(false)
	if err != nil {
		newErr := errors.New("error while processing mirror:" + err.Error())
		return newErr
	}

	err = core.ProcessRtMirror(true)
	if err != nil {
		newErr := errors.New("error while processing config:" + err.Error())
		cmdcontron.CliLog().Println(newErr)
	}

	return core.ExecCmd(c.Execute, c.withWaiting, c.ExecArgs...)
}
