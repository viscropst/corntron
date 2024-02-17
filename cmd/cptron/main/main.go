package main

import (
	"cryphtron"
	"cryphtron/cmd/cptron"
	"cryphtron/cmd/cptron/actions"
	"os"

	"github.com/rs/zerolog"
)

type cliFlags struct {
	*cptron.CmdFlag
}

func (f cliFlags) Init() *cliFlags {
	actions := actions.ActionMap()
	f.CmdFlag = cptron.CmdFlag{}.Prepare(actions)
	return &f
}

func main() {

	flags := cliFlags{}.Init()
	errLogger := cptron.CliLog(zerolog.ErrorLevel)

	action, err := flags.Parse()
	defer cptron.CliExit(err, err != nil && !flags.NoWaiting)
	if err != nil {
		errLogger.Println(err)
		return
	}

	var confBase []string
	if len(flags.ConfigBase) > 0 {
		confBase = append(confBase, flags.ConfigBase)
	}
	coreConfig := cryphtron.LoadCoreConfigWithRuningBase(flags.RunningBase, confBase...)

	err = action.BeforeCore(&coreConfig)
	if err != nil {
		errLogger.Println("error before load core", err)
		return
	}

	var core cryphtron.Core
	core, err = cryphtron.LoadCore(coreConfig)
	if err != nil {
		errLogger.Println("error while load core", err)
		return
	}

	err = action.Exec(&core)
	if err != nil {
		errLogger.Println(err.Error())
		return
	}

	os.Exit(0)

}
