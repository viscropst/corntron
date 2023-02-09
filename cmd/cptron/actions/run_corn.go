package actions

import (
	"cryphtron"
	"cryphtron/cmd/cptron"
	"cryphtron/core"
	"cryphtron/internal"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type runCorn struct {
	cptron.BaseAction
	appName       string
	args          []string
	globalWaiting bool
}

func (c *runCorn) ActionName() string {
	return "run-corn"
}

func (c *runCorn) ParseArg(info cptron.FlagInfo) error {
	argCmdIdx := info.Index + 1
	if len(os.Args) <= argCmdIdx && len(os.Args[argCmdIdx]) == 0 {
		return nil
	}
	c.appName = os.Args[argCmdIdx]
	if len(os.Args) > argCmdIdx+1 {
		c.args = os.Args[argCmdIdx+1:]
	}
	return nil
}

func (c *runCorn) BeforeCore(coreConfig *core.MainConfig) error {
	coreConfig.WithApp = true
	return nil
}

func (c *runCorn) InsertFlags(flag *cptron.CmdFlag) error {
	c.globalWaiting = !flag.NoWaiting
	return nil
}

func (c *runCorn) Exec(core *cryphtron.Core) error {
	err := core.ProcessRtBootstrap()
	if err != nil {
		newErr := errors.New("error while bootstrapping")
		return errors.Join(newErr, err)
	}

	err = core.ProcessRtMirror()
	if err != nil {
		newErr := errors.New("error while processing mirror")
		return errors.Join(newErr, err)
	}
	return c.execApp(c.appName, core)
}

func (c *runCorn) prepareApp(name string, coreObj *cryphtron.Core) (*core.CornsEnvConfig, error) {
	var err error
	app, ok := coreObj.CornsEnv[name]
	if !ok {
		return nil, errors.New("could not found the app named " + c.appName)
	}

	if !app.MetaOnly {
		scope := coreObj.ComposeRtEnv()
		app.AppendEnv(scope.Env)
		bootStrapDir := filepath.Join(coreObj.Config.CurrentDir, coreObj.Config.CornDir)
		bootStrapDir = filepath.Join(bootStrapDir, app.DirName)
		stat, _ := os.Stat(bootStrapDir)
		if stat == nil || (stat != nil && !stat.IsDir()) {
			_ = os.Mkdir(bootStrapDir, os.ModeDir|0o666)
			err = app.ExecuteBootstrap()
			if err != nil {
				newErr := errors.New("error while bootstrap app[" + app.ID + "]:")
				return nil, errors.Join(newErr, err)
			}
		}
	}

	err = app.ExecuteConfig()
	if err != nil {
		newErr := errors.New("error while configure app[" + app.ID + "]:")
		return nil, errors.Join(newErr, err)
	}

	for _, depend := range app.DependCorns {
		var cfg *core.CornsEnvConfig
		cfg, err = c.prepareApp(depend, coreObj)
		if err != nil {
			return nil, err
		}
		app.AppendVars(cfg.Vars)
	}
	return &app, nil
}

func (c *runCorn) execApp(name string, core *cryphtron.Core) error {
	app, err := c.prepareApp(name, core)
	if err != nil {
		return err
	}
	scope := core.ComposeAppEnv(app)

	pthVal := scope.Env["PATH"]
	pthVal = strings.Replace(pthVal, internal.PathPlaceHolder, core.Environ["PATH"], 1)
	scope.Env["PATH"] = pthVal
	cmd := &app.Exec
	cmd.Args = append(cmd.Args, c.args...)
	origWait := cmd.WithWaiting
	if c.globalWaiting {
		cmd.WithWaiting = c.globalWaiting
	}
	if origWait {
		cmd.WithWaiting = origWait
	}
	err = cmd.SetEnv(app.Env).Execute(scope.Vars)
	if err != nil {
		newErr := errors.New("error while executing")
		return errors.Join(newErr, err)
	}

	return nil
}
