package actions

import (
	"cryphtron"
	"cryphtron/cmd/cptron"
	"cryphtron/core"
	"errors"
	"os"
)

type runCorn struct {
	cptron.BaseAction
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

func (c *runCorn) ParseArg(info cptron.FlagInfo) error {
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

func (c *runCorn) InsertFlags(flag *cptron.CmdFlag) error {
	c.globalWaiting = flag.IsWaiting
	return nil
}

func (c *runCorn) Exec(core *cryphtron.Core) error {
	err := core.ProcessRtBootstrap(true)
	if err != nil {
		newErr := errors.New("error while bootstrapping")
		cptron.CliLog().Println(errors.Join(newErr, err))
	}

	err = core.ProcessRtMirror(true)
	if err != nil {
		newErr := errors.New("error while processing mirror")
		cptron.CliLog().Println(errors.Join(newErr, err))
	}

	err = core.ProcessRtConfig(true)
	if err != nil {
		newErr := errors.New("error while processing config")
		cptron.CliLog().Println(errors.Join(newErr, err))
	}
	return core.ExecCorn(c.appName, c.globalWaiting, c.args...)
}
