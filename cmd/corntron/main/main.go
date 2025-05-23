package main

import (
	"corntron"
	cmdcontron "corntron/cmd/corntron"
	"corntron/cmd/corntron/actions"
	"corntron/internal/utils/log"
	"net/url"
	"os"
	"strings"
)

type cliFlags struct {
	*cmdcontron.CmdFlag
}

func (f cliFlags) Init() *cliFlags {
	actions := actions.ActionMap()
	f.CmdFlag = cmdcontron.CmdFlag{}.Prepare(actions)
	return &f
}

func main() {

	flags := cliFlags{}.Init()
	errLogger := cmdcontron.CliLog(log.ErrorLevel)
	if !strings.HasSuffix(os.Args[0], "debug") {
		log.CLIOutputLevel = log.InfoLevel
	}
	cmdcontron.CliLog(log.DebugLevel).Println(os.Args, "len:", len(os.Args))
	for i, v := range os.Args {
		cmdcontron.CliLog(log.DebugLevel).Println("arg", i, "was", v, "url value", url.QueryEscape(v), "len", len(v))
	}

	action, err := flags.Parse()
	defer cmdcontron.CliExit(err, err != nil && (!cmdcontron.IsInTerminal()))
	if err != nil {
		errLogger.Println(err)
		return
	}

	runningBase, confBase := action.BeforeLoad(flags.CmdFlag)
	coreConfig := corntron.LoadCoreConfigWithRuningBase(runningBase, confBase...)

	err = action.BeforeCore(&coreConfig)
	if err != nil {
		errLogger.Println("error before load core", err)
		return
	}

	var core corntron.Core
	core, err = corntron.LoadCore(coreConfig)
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
