package main

import (
	"cryphtron"
	"cryphtron/cmd/cptron"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type cliFlags struct {
	*cptron.CmdFlag
	Execute  string
	ExecArgs []string
}

func (f cliFlags) Init() *cliFlags {
	actions := map[string]cptron.CmdActionFn{
		"exec": f.parseExec,
	}
	f.CmdFlag = cptron.CmdFlag{}.Prepare(actions)
	return &f
}

func (f *cliFlags) parseExec(info cptron.FlagInfo) error {
	argCmdIdx := info.Index + 1
	if len(os.Args) > argCmdIdx && len(os.Args[argCmdIdx]) > 0 {
		f.Execute = os.Args[argCmdIdx]
	} else {
		switch runtime.GOOS {
		case "windows":
			f.Execute = os.Getenv("COMSPEC")
		case "linux", "macos", "openbsd", "freebsd", "ios", "android":
			f.Execute = os.Getenv("SHELL")
			if len(f.Execute) == 0 {
				f.Execute = "/bin/sh"
			}
		}
		fmt.Println("warn:no command to exec,will use default shell or cmd")
	}

	if path, err := exec.LookPath(f.Execute); err != nil {
		return fmt.Errorf(
			"exec argument invalid:usage %s %s <command>",
			info.CmdName,
			info.Name)
	} else {
		f.Execute = path
	}

	if len(os.Args) > argCmdIdx+1 {
		f.ExecArgs = os.Args[argCmdIdx+1:]
	}

	return nil
}

func main() {

	flags := cliFlags{}.Init()
	err := flags.Parse()
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
		return
	}

	var confBase []string
	if len(flags.ConfigBase) > 0 {
		confBase = append(confBase, flags.ConfigBase)
	}
	coreConfig := cryphtron.LoadCoreConfig(confBase...)

	var core cryphtron.Core
	core, err = cryphtron.LoadCore(coreConfig)
	if err != nil {
		fmt.Println("error while load core", err)
		return
	}

	scope := core.ComposeRtEnv()
	for _, config := range core.EditorsEnv {
		config.Env = scope.Env
		err = config.ExecuteAll()
		if err != nil {
			fmt.Println("error while exec editor["+config.ID+"]:", err)
			return
		}
	}

	scope = core.ComposeEdEnv()
	cmd := cryphtron.Command{
		Exec: flags.Execute,
		Args: flags.ExecArgs,
	}
	err = cmd.SetEnv(scope.Env).Execute(scope.Vars)
	if err != nil {
		fmt.Println("error while exec", err)
		return
	}

}
