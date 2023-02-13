package cryphtron

import (
	"cryphtron/core"
	"cryphtron/internal"
	"errors"
	"strings"
)

func (c *Core) ExecCmd(command string, isWaiting bool, args ...string) error {
	scope := c.ComposeRtEnv()
	cmd := core.Command{
		Exec:       command,
		Args:       args,
		WithNoWait: !isWaiting,
	}
	return c.execCmd(&cmd, scope)
}

func (c *Core) execCmd(command *core.Command, scope *internal.ValueScope) error {
	pthVal := scope.Env["PATH"]
	pthVal = strings.Replace(pthVal, internal.PathPlaceHolder, c.Environ["PATH"], 1)
	scope.Env["PATH"] = pthVal

	err := command.SetEnv(scope.Env).Execute(scope.Vars)
	if err != nil {
		newErr := errors.New("error while executing")
		return errors.Join(newErr, err)
	}
	return nil
}
