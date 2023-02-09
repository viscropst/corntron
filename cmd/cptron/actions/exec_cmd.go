package actions

import (
	"cryphtron"
	"cryphtron/cmd/cptron"
	ct_core "cryphtron/core"
	"cryphtron/internal"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/skerkour/rz"
)

type execCmd struct {
	cptron.BaseAction
	Execute     string
	ExecArgs    []string
	withWaiting bool
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
		cptron.CliLog(rz.WarnLevel).Println("warn:no command to exec,will use default shell or cmd")
	}

	if path, err := exec.LookPath(c.Execute); err != nil {
		errBuilder := strings.Builder{}
		errBuilder.WriteString("exec argument invalid: usage ")
		errBuilder.WriteString(info.CmdName + " ")
		errBuilder.WriteString(info.Name + " ")
		errBuilder.WriteString("<command>")
		return errors.New(errBuilder.String())
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

func (c *execCmd) InsertFlags(flag *cptron.CmdFlag) error {
	c.withWaiting = !flag.NoWaiting
	return nil
}

func (c *execCmd) Exec(core *cryphtron.Core) error {
	var err error
	scope := core.ComposeRtEnv()

	err = core.ProcessRtBootstrap()
	if err != nil {
		newErr := errors.New("error while bootstrapping:")
		return errors.Join(newErr, err)
	}

	err = core.ProcessRtMirror()
	if err != nil {
		newErr := errors.New("error while processing mirror:")
		return errors.Join(newErr, err)
	}

	cmd := ct_core.Command{
		Exec:        c.Execute,
		Args:        c.ExecArgs,
		WithWaiting: c.withWaiting,
	}

	pthVal := scope.Env["PATH"]
	pthVal = strings.Replace(pthVal, internal.PathPlaceHolder, core.Environ["PATH"], 1)
	scope.Env["PATH"] = pthVal

	err = cmd.SetEnv(scope.Env).Execute(scope.Vars)
	if err != nil {
		newErr := errors.New("error while executing:")
		return errors.Join(newErr, err)
	}

	return nil
}
