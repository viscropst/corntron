package actions

import (
	"cryphtron"
	"cryphtron/cmd/cptron"
	"cryphtron/core"
	"errors"
	"os"
)

type runCornConfig struct {
	cptron.BaseAction
	fileName      string
	args          []string
	globalWaiting bool
}

func init() {
	appendAction(&runCornConfig{})
}

func (c *runCornConfig) ActionName() string {
	return "run-" + core.CornsIdentifier + "-config"
}

func (c *runCornConfig) ParseArg(info cptron.FlagInfo) error {
	argCmdIdx := info.Index + 1
	isValidArgNum := len(os.Args) <= argCmdIdx
	if isValidArgNum {
		return nil
	}
	if isValidArgNum && len(os.Args[argCmdIdx]) == 0 {
		return nil
	}
	c.fileName = os.Args[argCmdIdx]
	if len(os.Args) > argCmdIdx+1 {
		c.args = os.Args[argCmdIdx+1:]
	}
	return nil
}

func (c *runCornConfig) BeforeCore(coreConfig *core.MainConfig) error {
	coreConfig.WithCorn = true
	return nil
}

func (c *runCornConfig) InsertFlags(flag *cptron.CmdFlag) error {
	c.globalWaiting = !flag.NoWaiting
	return nil
}

func (c *runCornConfig) Exec(inCore *cryphtron.Core) error {
	err := inCore.ProcessRtBootstrap(true)
	if err != nil {
		newErr := errors.New("error while bootstrapping")
		cptron.CliLog().Println(errors.Join(newErr, err))
	}

	err = inCore.ProcessRtMirror(true)
	if err != nil {
		newErr := errors.New("error while processing mirror")
		cptron.CliLog().Println(errors.Join(newErr, err))
	}

	err = inCore.ProcessRtConfig(true)
	if err != nil {
		newErr := errors.New("error while processing config")
		cptron.CliLog().Println(errors.Join(newErr, err))
	}

	tmpEnv := core.BaseEnv(inCore.Config)
	tmpEnv.Top = inCore.ValueScope

	config, err := core.LoadCornConfig(c.fileName, tmpEnv)
	if err != nil {
		newErr := errors.New("error while loading corn config")
		return newErr
	}

	return inCore.ExecCornConfig(config, c.globalWaiting, c.args...)
}
