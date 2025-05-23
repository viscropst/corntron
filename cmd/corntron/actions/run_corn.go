package actions

import (
	"corntron"
	cmdcontron "corntron/cmd/corntron"
	"corntron/core"
	"errors"
	"os"
)

type runCorn struct {
	cmdcontron.BaseAction
	appName       string
	args          []string
	globalWaiting bool
}

func init() {
	appendAction(&runCorn{})
}

func (c *runCorn) ActionName() string {
	return "run-" + core.CornsIdentifier
}

func (c *runCorn) ParseArg(info cmdcontron.FlagInfo) error {
	argCmdIdx := info.Index + 1
	isValidArgNum := len(os.Args) <= argCmdIdx
	if isValidArgNum {
		return nil
	}
	if isValidArgNum && len(os.Args[argCmdIdx]) == 0 {
		return nil
	}
	c.appName = os.Args[argCmdIdx]
	if len(os.Args) > argCmdIdx+1 {
		c.args = os.Args[argCmdIdx+1:]
	}
	return nil
}

func (c *runCorn) BeforeCore(coreConfig *core.MainConfig) error {
	coreConfig.WithCorn = true
	return nil
}

func (c *runCorn) BeforeLoad(flag *cmdcontron.CmdFlag) (string, []string) {
	c.globalWaiting = !flag.NoWaiting
	return c.BaseAction.BeforeLoad(flag)
}

func (c *runCorn) Exec(core *corntron.Core) error {
	err := core.ProcessRtBootstrap(true)
	if err != nil {
		newErr := errors.New("error while bootstrapping:" + err.Error())
		cmdcontron.CliLog().Println(errors.Join(newErr, err))
	}

	err = core.ProcessRtMirror(true)
	if err != nil {
		newErr := errors.New("error while processing mirror:" + err.Error())
		cmdcontron.CliLog().Println(errors.Join(newErr, err))
	}

	err = core.ProcessRtConfig(true)
	if err != nil {
		newErr := errors.New("error while processing config:" + err.Error())
		cmdcontron.CliLog().Println(newErr)
	}
	return core.ExecCorn(c.appName, c.globalWaiting, c.args...)
}
