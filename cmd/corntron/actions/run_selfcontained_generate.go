package actions

import (
	"corntron"
	cmdcontron "corntron/cmd/corntron"
	"corntron/core"
	"errors"
	"flag"
)

type runGenerateSelfcontained struct {
	cmdcontron.BaseAction
	fileName        string
	args            []string
	globalWaiting   bool
	flagSet         *flag.FlagSet
	configDirAsBase bool
}

func init() {
	cfg := runGenerateSelfcontained{}
	cfg.flagSet = flag.NewFlagSet("", flag.ContinueOnError)
	cfg.flagSet.BoolVar(&cfg.configDirAsBase,
		"dir-as-base", false, "use the current corn file's dir as config base")
	appendAction(&cfg)
}

func (c *runGenerateSelfcontained) ActionName() string {
	return "generate-selfcontained"
}

func (c *runGenerateSelfcontained) ParseArg(info cmdcontron.FlagInfo) error {
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

func (c *runGenerateSelfcontained) BeforeLoad(flags *cmdcontron.CmdFlag) (string, []string) {
	return c.BaseAction.BeforeLoad(flags)
}

func (c *runGenerateSelfcontained) BeforeCore(coreConfig *core.MainConfig) error {
	coreConfig.WithCorn = true
	return nil
}

func (c *runGenerateSelfcontained) Exec(inCore *corntron.Core) error {
	return errors.New("not implemented")
}
