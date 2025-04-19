package cptron

import (
	"cryphtron"
	"cryphtron/core"
	"os"
)

type CmdAction interface {
	ActionName() string
	ParseArg(info FlagInfo) error
	BeforeLoad(flags *CmdFlag) (string, []string)
	BeforeCore(coreConfig *core.MainConfig) error
	Exec(core *cryphtron.Core) error
}

type BaseAction struct{}

func (b BaseAction) ActionName() string {
	return "empty"
}

func (b BaseAction) ParseArg(info FlagInfo) error {
	return nil
}

func (b BaseAction) BeforeCore(flags *CmdFlag, coreConfig *core.MainConfig) error {
	return nil
}

func (b BaseAction) BeforeLoad(flags *CmdFlag) (string, []string) {
	var confBase []string
	if len(flags.ConfigBase) > 0 {
		confBase = append(confBase, flags.ConfigBase)
	}
	return flags.RunningBase, confBase
}

func (b BaseAction) Exec(core *cryphtron.Core) error {
	return nil
}

func CliExit(err error, forceWait bool) {
	if forceWait {
		CliLog().Println("press any key to continue")
		tmp := make([]byte, 1)
		os.Stdin.ReadAt(tmp, 0)
	}
	if err != nil {
		os.Exit(1)
	}
}
