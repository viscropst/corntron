package cptron

import (
	"cryphtron"
	"cryphtron/core"
	"fmt"
	"os"
)

type CmdAction interface {
	ActionName() string
	ParseArg(info FlagInfo) error
	BeforeCore(coreConfig *core.MainConfig) error
	Exec(core *cryphtron.Core) error
}

func GracefulExit(err error) {
	if !IsInTerminal() {
		CliLog().Println("press any key to exit")
		fmt.Scanln()
	}
	if err != nil {
		os.Exit(1)
	}
}
