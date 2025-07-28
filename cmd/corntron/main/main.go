package main

import (
	"corntron"
	cmdcontron "corntron/cmd/corntron"
	"corntron/cmd/corntron/actions"
	"corntron/internal/log"
	"net/url"
	"os"
	"path/filepath"
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
	_, selfFile := filepath.Split(os.Args[0])
	selfName := strings.TrimSuffix(selfFile, filepath.Ext(selfFile))
	if !strings.HasSuffix(selfName, "debug") {
		log.CLIOutputLevel = log.InfoLevel
	}
	if strings.HasPrefix(selfName, "__debug") && log.CLIOutputLevel != log.DebugLevel {
		log.CLIOutputLevel = log.DebugLevel
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
	if len(flags.MirrorType) > 0 {
		coreConfig.MirrorType = flags.MirrorType
	}
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
