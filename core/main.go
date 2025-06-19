package core

import (
	"corntron/internal"
	"errors"
	"strings"
)

type Core struct {
	*ValueScope
	Environ     map[string]string
	ProfileDir  string
	Config      MainConfig
	CornsEnv    map[string]CornsEnvConfig
	RuntimesEnv []RtEnvConfig
}

func (c *Core) ExecCmd(command string, isWaiting bool, args ...string) error {
	scope := c.ComposeRtEnv()
	cmd := Command{
		Exec:          command,
		Args:          args,
		GlobalWaiting: isWaiting,
	}

	return c.execCmd(&cmd, scope)
}

func (c *Core) checkByPATH(command *Command, scope *ValueScope) (string, error) {
	pthVal := scope.Env["PATH"]
	pthVal = strings.Replace(pthVal, PathPlaceHolder, c.Environ["PATH"], 1)
	scope.Env["PATH"] = pthVal
	return internal.GetExecPath(command.Exec, pthVal)
}

func (c *Core) execCmd(command *Command, scope *ValueScope) error {
	if exec, err := c.checkByPATH(command, scope); err == nil {
		command.Exec = exec
	}
	err := command.SetEnv(scope.Env).ExecWithAttr(scope.Vars)
	if err != nil {
		newErr := errors.New("error while executing")
		return errors.Join(newErr, err)
	}
	return nil
}

func (c *Core) Prepare() {
	if c.Environ != nil {
		return
	}
	c.fillEnviron()
}
