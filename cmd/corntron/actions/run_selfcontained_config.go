package actions

import (
	"bytes"
	"corntron"
	cmdcontron "corntron/cmd/corntron"
	"errors"
	"flag"
)

type runCornConfigDirect struct {
	cmdcontron.BaseAction
	fileName        string
	args            []string
	globalWaiting   bool
	flagSet         *flag.FlagSet
	configDirAsBase bool
}

func init() {
	cfg := runCornConfigDirect{}
	cfg.flagSet = flag.NewFlagSet("", flag.ContinueOnError)
	appendAction(&cfg)
}

func (c *runCornConfigDirect) ActionName() string {
	return "run-selfcontained-config"
}

func (c *runCornConfigDirect) ParseArg(info cmdcontron.FlagInfo) error {
	err := c.flagSet.Parse(info.Args[info.Index+1:])
	if err != nil {
		return err
	}
	argCmdIdx := info.Index + c.flagSet.NFlag() + 1
	if (info.TotalLen + 1) < (argCmdIdx + (*c.flagSet).NFlag()) {
		argCmdIdx -= 1
	}
	isValidArgNum := len(info.Args) <= argCmdIdx
	if isValidArgNum {
		return nil
	}
	if isValidArgNum && len(info.Args[argCmdIdx]) == 0 {
		return nil
	}

	c.fileName = info.Args[argCmdIdx]
	if len(info.Args) > argCmdIdx+1 {
		c.args = info.Args[argCmdIdx+1:]
	}
	return nil
}

func (c *runCornConfigDirect) BeforeLoad(flags *cmdcontron.CmdFlag) (string, []string) {
	c.globalWaiting = !flags.NoWaiting
	return c.BaseAction.BeforeLoad(flags)
}

func (c *runCornConfigDirect) BeforeCore(coreConfig *corntron.MainConfig) error {
	coreConfig.WithCorn = true
	return nil
}

func (c *runCornConfigDirect) Exec(inCore *corntron.Core) error {
	err := inCore.ProcessRtBootstrap(true)
	if err != nil {
		newErr := errors.New("error while bootstrapping:" + err.Error())
		cmdcontron.CliLog().Println(newErr)
	}

	err = inCore.ProcessRtMirror(true)
	if err != nil {
		newErr := errors.New("error while processing mirror:" + err.Error())
		cmdcontron.CliLog().Println(newErr)
	}

	err = inCore.ProcessRtConfig(true)
	if err != nil {
		newErr := errors.New("error while processing config:" + err.Error())
		cmdcontron.CliLog().Println(newErr)
	}

	configBytes := make([]byte, 8192)
	config, err := corntron.LoadCornConfigReader(inCore, c.fileName, bytes.NewBuffer(configBytes))
	if err != nil {
		newErr := errors.New("error while loading corn config:" + err.Error())
		return newErr
	}

	return inCore.ExecCornConfig(*config, c.globalWaiting, c.args...)
}
