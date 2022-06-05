package main

import (
	"cryphtron"
	"cryphtron/cmd/cptron"
	"cryphtron/cmd/cptron/actions"
	"fmt"
	"os"
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
	action, err := flags.Parse()
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
		return
	}

	var confBase []string
	if len(flags.ConfigBase) > 0 {
		confBase = append(confBase, flags.ConfigBase)
	}
	coreConfig := cryphtron.LoadCoreConfig(confBase...)

	err = action.BeforeCore(&coreConfig)
	if err != nil {
		fmt.Println("error before load core", err)
		return
	}

	var core cryphtron.Core
	core, err = cryphtron.LoadCore(coreConfig)
	if err != nil {
		fmt.Println("error while load core", err)
		return
	}

	err = action.Exec(&core)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
