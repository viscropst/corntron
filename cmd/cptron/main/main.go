package main

import (
	"cryphtron"
	"cryphtron/cmd/cptron"
	"cryphtron/cmd/cptron/actions"
	"cryphtron/internal/utils/log"
	"net/url"
	"os"
	"strings"
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
	errLogger := cptron.CliLog(log.ErrorLevel)
	if !strings.HasSuffix(os.Args[0], "debug") {
		log.CLIOutputLevel = log.InfoLevel
	}
	cptron.CliLog(log.DebugLevel).Println(os.Args, "len:", len(os.Args))
	for i, v := range os.Args {
		cptron.CliLog(log.DebugLevel).Println("arg", i, "was", v, "url value", url.QueryEscape(v), "len", len(v))
	}

	action, err := flags.Parse()
	defer cptron.CliExit(err, err != nil && (!cptron.IsInTerminal()))
	if err != nil {
		errLogger.Println(err)
		return
	}

	runningBase, confBase := action.BeforeLoad(flags.CmdFlag)
	coreConfig := cryphtron.LoadCoreConfigWithRuningBase(runningBase, confBase...)

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
