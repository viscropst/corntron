package cryphtron

import (
	"cryphtron/core"
	"cryphtron/internal"
	"errors"
	"os"
	"os/exec"
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

func (c *Core) checkByPATH(command *core.Command, scope *internal.ValueScope) (string, error) {
	pthVal := scope.Env["PATH"]
	pthVal = strings.Replace(pthVal, internal.PathPlaceHolder, c.Environ["PATH"], 1)
	scope.Env["PATH"] = pthVal

	os.Setenv("PATH", scope.Env["PATH"])
	defer os.Unsetenv("PATH")
	if path, err := exec.LookPath(command.Exec); err != nil {
		errBuilder := strings.Builder{}
		errBuilder.WriteString("exec argument invalid: the command could not found")
		return command.Exec, errors.New(errBuilder.String())
	} else {
		return path, nil
	}
}

func (c *Core) execCmd(command *core.Command, scope *internal.ValueScope) error {
	if exec, err := c.checkByPATH(command, scope); err == nil {
		command.Exec = exec
	}
	err := command.SetEnv(scope.Env).Execute(scope.Vars)
	if err != nil {
		newErr := errors.New("error while executing")
		return errors.Join(newErr, err)
	}
	return nil
}
