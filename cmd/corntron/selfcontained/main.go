package main

import (
	"corntron"
	cmdcorntron "corntron/cmd/corntron"
	"corntron/cmd/corntron/actions"
	"corntron/internal/log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type cliFlags struct {
	*cmdcorntron.CmdFlag
}

func (f cliFlags) Init() *cliFlags {
	actionsMap := make(map[string]cmdcorntron.CmdAction)
	for k, v := range actions.ActionMap() {
		if k != "run-cmd" {
			continue
		}
		actionsMap[k] = v
	}
	f.CmdFlag = f.CmdFlag.Prepare()
	return &f
}

func main() {
	flags := cliFlags{}.Init()
	_, selfFile := filepath.Split(os.Args[0])
	selfName := strings.TrimSuffix(selfFile, filepath.Ext(selfFile))
	cmdcorntron.CliLog(log.InfoLevel).Println(selfName, "version: ", cmdcorntron.Version())
	if !strings.HasSuffix(selfName, "debug") {
		log.CLIOutputLevel = log.InfoLevel
	}
	if strings.HasPrefix(selfName, "__debug") && log.CLIOutputLevel != log.DebugLevel {
		log.CLIOutputLevel = log.DebugLevel
	}
	cmdcorntron.CliLog(log.DebugLevel).Println(os.Args, "len:", len(os.Args))
	for i, v := range os.Args {
		cmdcorntron.CliLog(log.DebugLevel).Println("arg", i, "was", v, "url value", url.QueryEscape(v), "len", len(v))
	}

	flagInfo, err := flags.Parse()
	if err != nil {
		cmdcorntron.ErrorLog(err)
		return
	}
	actionIdentifier := "run-selfcontained-config"
	if arg := flagInfo.Args[flagInfo.Index]; arg == "run-cmd" {
		actionIdentifier = "run-cmd"
	}
	action := actions.ActionMap()[actionIdentifier]
	runningBase, confBase := action.BeforeLoad(flags.CmdFlag)
	coreConfig := corntron.LoadCoreConfigWithRuningBase(runningBase, confBase...)
	if len(flags.MirrorType) > 0 {
		coreConfig.MirrorType = flags.MirrorType
	}
	err = action.BeforeCore(&coreConfig)
	if err != nil {
		cmdcorntron.ErrorLog(err, "error before load core")
		return
	}

	var core corntron.Core
	core, err = corntron.LoadCore(coreConfig)
	if err != nil {
		cmdcorntron.ErrorLog(err, "error while load core")
		return
	}

	err = action.Exec(&core)
	if err != nil {
		cmdcorntron.ErrorLog(err)
		return
	}
}
