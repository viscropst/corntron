package core

import (
	"corntron/internal"
	"errors"
)

type Core struct {
	*ValueScope
	Environ     map[string]string
	EnvironPath PathList
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
	pthVal := c.EnvironPath
	pthVal = pthVal.AppendList(scope.EnvPath)
	return internal.GetExecPath(command.Exec, pthVal.String())
}

func (c *Core) execCmd(command *Command, scope *ValueScope) error {
	if exec, err := c.checkByPATH(command, scope); err == nil {
		command.Exec = exec
	}
	command.withAttr = true
	command.EnvPath = c.EnvironPath.AppendList(scope.EnvPath)
	err := command.SetEnv(scope.Env).Execute(scope.Vars)
	if err != nil {
		newErr := errors.New("error while executing:" + err.Error())
		return newErr
	}
	return nil
}

func (c *Core) Prepare() {
	c.EnvironPath = EnvironPathList()
	if len(c.Env) > 0 {
		return
	}
	c.Env = internal.FillEnviron(c.ProfileDir)
	if len(c.Environ) > 0 {
		return
	}
	c.Environ = internal.GetEnvironMap()
}
