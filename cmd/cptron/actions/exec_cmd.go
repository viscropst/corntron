package actions

import (
	"cryphtron"
	"cryphtron/cmd/cptron"
	ct_core "cryphtron/core"
	"cryphtron/internal"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type execCmd struct {
	Execute  string
	ExecArgs []string
}

func (c *execCmd) ActionName() string {
	return "exec-cmd"
}

func (c *execCmd) ParseArg(info cptron.FlagInfo) error {
	argCmdIdx := info.Index + 1
	if len(os.Args) > argCmdIdx && len(os.Args[argCmdIdx]) > 0 {
		c.Execute = os.Args[argCmdIdx]
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
		fmt.Println("warn:no command to exec,will use default shell or cmd")
	}

	if path, err := exec.LookPath(c.Execute); err != nil {
		return fmt.Errorf(
			"exec argument invalid:usage %s %s <command>",
			info.CmdName,
			info.Name)
	} else {
		c.Execute = path
	}

	if len(os.Args) > argCmdIdx+1 {
		c.ExecArgs = os.Args[argCmdIdx+1:]
	}

	return nil
}

func (c *execCmd) BeforeCore(coreConfig *ct_core.MainConfig) error {
	return nil
}

func (c *execCmd) Exec(core *cryphtron.Core) error {
	var err error
	scope := core.ComposeRtEnv()

	err = core.ProcessRtMirror()
	if err != nil {
		err = fmt.Errorf("error while processing mirror %s", err.Error())
		return err
	}

	cmd := ct_core.Command{
		Exec: c.Execute,
		Args: c.ExecArgs,
	}

	pthVal := scope.Env["PATH"]
	pthVal = strings.Replace(pthVal, internal.PathPlaceHolder, core.Environ["PATH"], 1)
	scope.Env["PATH"] = pthVal

	err = cmd.SetEnv(scope.Env).Execute(scope.Vars)
	if err != nil {
		err = fmt.Errorf("error while exec %s", err.Error())
		return err
	}

	return nil
}
