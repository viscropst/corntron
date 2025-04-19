package cptron

import (
	"cryphtron/core"
	"cryphtron/internal/utils/log"
	"errors"
	"flag"
	"os"
	"path"
	"strings"
)

type FlagInfo struct {
	Name     string
	Index    int
	CmdName  string
	TotalLen int
	Args     []string
}

type CmdFlag struct {
	host             *flag.FlagSet
	flagLen          int
	argLen           int
	osArgLen         int
	actions          map[string]CmdAction
	args             []string
	NoWaiting        bool
	ConfigBase       string
	EnvDirname       string
	RuntimeBase      string
	EditorBase       string
	RunningBase      string
	PassToCornConfig bool
}

func (f CmdFlag) Prepare(actions map[string]CmdAction) *CmdFlag {
	result := &CmdFlag{}
	if len(actions) == 0 && actions != nil {
		return nil
	}
	result.host = flag.CommandLine
	result.actions = actions
	result.host.Usage = func() {
		CliLog().Println(path.Base(result.host.Name()) + " [options] <actions> [args]")
		actKeys := make([]string, 0)
		for k := range result.actions {
			actKeys = append(actKeys, k)
		}
		CliLog().Printf("actions was: %v \n", actKeys)
		CliLog().Println("options has:")
		result.host.PrintDefaults()
		CliExit(nil, !IsInTerminal() || !result.NoWaiting)
	}
	result.host.BoolVar(&result.NoWaiting, "no-waiting", false, "executing cryptron without waiting")
	result.host.StringVar(&result.ConfigBase, "cfg-base", "", "/path/to/your/<cryphtron config folder>")
	result.host.StringVar(&result.RunningBase, "running-base", "", "/path/to/your/<cryphtron running folder>")
	result.host.StringVar(&result.RuntimeBase, "rt-base", "", "/path/to/your/<runtime profiles folder>")
	result.host.StringVar(&result.EditorBase, "corn-base", "", "/path/to/your/<corns profiles folder>")
	result.host.StringVar(&result.EnvDirname, "env-dirname", "", "<folder name of env files to store>")
	return result
}

func (f *CmdFlag) Parse() (CmdAction, error) {
	startIdx := 1
	f.args = os.Args
	if len(f.args) > 2 {
		if len(f.args[1]) < 2 {
			startIdx += 1
		}
		if sp := strings.Split(strings.TrimSpace(f.args[startIdx]), " "); len(sp) > 1 {
			f.args = make([]string, 1)
			f.args[0] = os.Args[0]
			f.args = append(f.args, sp...)
			f.args = append(f.args, os.Args[startIdx+1:]...)
		}
	}
	err := f.host.Parse(f.args[startIdx:])
	if err != nil && !f.PassToCornConfig {
		return nil, err
	}
	f.flagLen = f.host.NFlag() * 2
	f.argLen = f.host.NArg()
	f.osArgLen = len(f.args) - 1
	if (f.osArgLen-f.argLen) < 0 || (f.argLen+f.flagLen) == 0 {
		return nil, errors.New("invalid length of args,use '-help' for usage")
	}
	idxArgAct := startIdx - 1
	idxArgAct = idxArgAct + f.flagLen + 1
	if (f.osArgLen + 1) < (idxArgAct + f.argLen) {
		idxArgAct -= 1
	}
	info := FlagInfo{
		Name:     os.Args[idxArgAct],
		Index:    idxArgAct,
		CmdName:  f.host.Name(),
		TotalLen: f.osArgLen,
		Args:     f.args,
	}
	CliLog(log.DebugLevel).Println("startIndex", startIdx)
	if fileArg := f.args[idxArgAct]; strings.HasSuffix(fileArg, core.CornConfigExt) {
		action := f.actions["run-"+core.CornsIdentifier+"-config"]
		info.Index = idxArgAct - 1
		return action, action.ParseArg(info)
	}
	if action, ok := f.actions[f.args[idxArgAct]]; ok {
		return action, action.ParseArg(info)
	}
	return nil, errors.New("invalid action,use '-help' for actions")
}
