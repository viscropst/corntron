package main

import (
	"corntron"
	cmdcorntron "corntron/cmd/corntron"
	"corntron/cmd/corntron/actions"
	"corntron/core"
	"corntron/internal/log"
	"errors"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type cliFlags struct {
	*cmdcorntron.CmdFlag
	actions     map[string]cmdcorntron.CmdAction
	NoWaiting   bool
	EnvDirname  string
	RuntimeBase string
	EditorBase  string
}

func (f cliFlags) Init() *cliFlags {
	f.actions = make(map[string]cmdcorntron.CmdAction)
	for k, v := range actions.ActionMap() {
		if k == "run-selfcontained-config" {
			continue
		}
		f.actions[k] = v
	}
	f.CmdFlag = cmdcorntron.CmdFlag{}.Prepare()
	f.Host.BoolVar(&f.NoWaiting, "no-waiting", false, "executing cryptron without waiting")
	f.Host.StringVar(&f.RuntimeBase, "rt-base", "", "/path/to/your/<runtime profiles folder>")
	f.Host.StringVar(&f.EditorBase, "corn-base", "", "/path/to/your/<corns profiles folder>")
	f.Host.StringVar(&f.EnvDirname, "env-dirname", "", "<folder name of env files to store>")
	f.Host.StringVar(&f.MirrorType, "mirror-type", "", "mirror type, default is without mirror")
	f.Host.StringVar(&f.ConfigBase, "cfg-base", "", "/path/to/your/<corntron config folder>")
	f.Host.StringVar(&f.RunningBase, "running-base", "", "/path/to/your/<corntron running folder>")
	f.Host.Usage = func() {
		cmdcorntron.CliLog().Println(path.Base(f.Host.Name()) + " " + cmdcorntron.Version() + " [options] <actions> [args]")
		actKeys := make([]string, 0)
		for k := range f.actions {
			actKeys = append(actKeys, k)
		}
		cmdcorntron.CliLog().Printf("actions was: %v \n", actKeys)
		cmdcorntron.CliLog().Println("options has:")
		f.Host.PrintDefaults()
		cmdcorntron.CliExit(nil, !cmdcorntron.IsInTerminal())
	}
	return &f
}

func (f *cliFlags) Parse() (cmdcorntron.CmdAction, error) {
	info, err := f.CmdFlag.Parse()
	if err != nil {
		return nil, err
	}
	idxArgAct := info.Index
	if fileArg := info.Args[idxArgAct]; strings.HasSuffix(fileArg, core.CornConfigExt) {
		action := f.actions["run-"+core.CornsIdentifier+"-config"]
		info.Index = idxArgAct - 1
		return action, action.ParseArg(*info)
	}
	if action, ok := f.actions[info.Args[idxArgAct]]; ok {
		return action, action.ParseArg(*info)
	}
	return nil, errors.New("invalid action,use '-help' for actions")
}

func main() {
	flags := cliFlags{}.Init()
	_, selfFile := filepath.Split(os.Args[0])
	selfName := strings.TrimSuffix(selfFile, filepath.Ext(selfFile))
	cmdcorntron.CliLog(log.InfoLevel).Println(selfName, "version:", cmdcorntron.Version())
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

	action, err := flags.Parse()
	if err != nil {
		cmdcorntron.ErrorLog(err)
		return
	}

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
