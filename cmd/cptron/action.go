package cptron

import (
	"cryphtron"
	"cryphtron/core"
	"os"
)

type CmdAction interface {
	ActionName() string
	InsertFlags(flag *CmdFlag) error
	ParseArg(info FlagInfo) error
	BeforeCore(coreConfig *core.MainConfig) error
	Exec(core *cryphtron.Core) error
}

type BaseAction struct{}

func (b BaseAction) ActionName() string {
	return "empty"
}

func (b BaseAction) InsertFlags(flag *CmdFlag) error {
	return nil
}

func (b BaseAction) ParseArg(info FlagInfo) error {
	return nil
}

func (b BaseAction) BeforeCore(coreConfig *core.MainConfig) error {
	return nil
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
