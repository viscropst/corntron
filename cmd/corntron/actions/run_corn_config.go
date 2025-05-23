package actions

import (
	"corntron"
	cmdcontron "corntron/cmd/corntron"
	"corntron/core"
	"errors"
	"flag"
	"os"
	"path/filepath"
)

type runCornConfig struct {
	cmdcontron.BaseAction
	fileName        string
	args            []string
	globalWaiting   bool
	flagSet         *flag.FlagSet
	configDirAsBase bool
}

func init() {
	cfg := runCornConfig{}
	cfg.flagSet = flag.NewFlagSet("", flag.ContinueOnError)
	cfg.flagSet.BoolVar(&cfg.configDirAsBase,
		"dir-as-base", false, "use the current corn file's dir as config base")
	appendAction(&cfg)
}

func (c *runCornConfig) ActionName() string {
	return "run-" + core.CornsIdentifier + "-config"
}

func (c *runCornConfig) ParseArg(info cmdcontron.FlagInfo) error {
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

func (c *runCornConfig) BeforeLoad(flags *cmdcontron.CmdFlag) (string, []string) {
	c.globalWaiting = !flags.NoWaiting
	if c.configDirAsBase {
		base := c.fileName
		if filepath.IsLocal(c.fileName) {
			wd, _ := os.Getwd()
			base = filepath.Join(wd, base)
		}
		flags.ConfigBase = filepath.Dir(base)
	}
	return c.BaseAction.BeforeLoad(flags)
}

func (c *runCornConfig) BeforeCore(coreConfig *core.MainConfig) error {
	coreConfig.WithCorn = true
	return nil
}

func (c *runCornConfig) Exec(inCore *corntron.Core) error {
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

	tmpEnv := core.BaseEnv(inCore.Config)
	tmpEnv.Top = inCore.ValueScope

	config, err := core.LoadCornConfig(c.fileName, tmpEnv)
	if err != nil {
		newErr := errors.New("error while loading corn config:" + err.Error())
		return newErr
	}

	return inCore.ExecCornConfig(config, c.globalWaiting, c.args...)
}
