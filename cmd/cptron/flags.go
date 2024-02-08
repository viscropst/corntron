package cptron

import (
	"errors"
	"flag"
	"os"
	"path"
)

type FlagInfo struct {
	Name    string
	Index   int
	CmdName string
}

type CmdFlag struct {
	host        *flag.FlagSet
	flagLen     int
	argLen      int
	osArgLen    int
	actions     map[string]CmdAction
	NoWaiting   bool
	ConfigBase  string
	EnvDirname  string
	RuntimeBase string
	EditorBase  string
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
	result.host.BoolVar(&result.NoWaiting, "no-wait", false, "executing cryptron without waiting")
	result.host.StringVar(&result.ConfigBase, "cfg-base", "", "/path/to/your/<cryphtron config folder>")
	result.host.StringVar(&result.RuntimeBase, "rt-base", "", "/path/to/your/<runtime profiles folder>")
	result.host.StringVar(&result.EditorBase, "app-base", "", "/path/to/your/<editor profiles folder>")
	result.host.StringVar(&result.EnvDirname, "env-dirname", "", "<folder name of env files to store>")
	return result
}

func (f *CmdFlag) Parse() (CmdAction, error) {
	err := f.host.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	f.flagLen = f.host.NFlag() * 2
	f.argLen = f.host.NArg()
	f.osArgLen = len(os.Args) - 1
	if (f.osArgLen-f.argLen) < 0 || (f.argLen+f.flagLen) == 0 {
		return nil, errors.New("invalid length of args,use '-help' for usage")
	}
	idxArgAct := f.flagLen + 1
	if (f.osArgLen + 1) < (idxArgAct + f.argLen) {
		idxArgAct -= 1
	}
	info := FlagInfo{
		Name:    os.Args[idxArgAct],
		Index:   idxArgAct,
		CmdName: f.host.Name(),
	}
	if action, ok := f.actions[os.Args[idxArgAct]]; ok {
		_ = action.InsertFlags(f)
		return action, action.ParseArg(info)
	}
	return nil, errors.New("invalid action,use '-help' for actions")
}
