package actions

import (
	"cryphtron"
	"cryphtron/cmd/cptron"
	"cryphtron/core"
	"fmt"
	"os"
	"path/filepath"
)

type execApp struct {
	appName string
}

func (c *execApp) ActionName() string {
	return "exec-app"
}

func (c *execApp) ParseArg(info cptron.FlagInfo) error {
	argCmdIdx := info.Index + 1
	if len(os.Args) <= argCmdIdx && len(os.Args[argCmdIdx]) == 0 {
		return nil
	}
	c.appName = os.Args[argCmdIdx]
	return nil
}

func (c *execApp) BeforeCore(coreConfig *core.MainConfig) error {
	coreConfig.WithApp = true
	return nil
}

func (c *execApp) Exec(core *cryphtron.Core) error {

	return c.execApp(c.appName, core)
}

func (c *execApp) execApp(name string, core *cryphtron.Core) error {
	var err error
	app, ok := core.AppsEnv[name]
	if !ok {
		return fmt.Errorf("could not found the app named %s", c.appName)
	}
	scope := core.ComposeRtEnv()

	app.AppendEnv(scope.Env)

	bootStrapDir := filepath.Join(core.Config.CurrentDir, core.Config.AppDir)
	bootStrapDir = filepath.Join(bootStrapDir, app.DirName)
	stat, _ := os.Stat(bootStrapDir)
	if stat == nil || (stat != nil && !stat.IsDir()) {
		_ = os.Mkdir(bootStrapDir, os.ModeDir|0o666)
		err = app.ExecuteBootstrap()
		if err != nil {
			err = fmt.Errorf("error while bootstrap app["+app.ID+"]:%S", err.Error())
			return err
		}
	}

	err = app.ExecuteConfig()
	if err != nil {
		err = fmt.Errorf("error while configure app["+app.ID+"]:%S", err.Error())
		return err
	}

	for _, depend := range app.DependApps {
		err = c.execApp(depend, core)
		if err != nil {
			return err
		}
	}
	scope = core.ComposeAppEnv()
	cmd := &app.Exec
	err = cmd.SetEnv(scope.Env).Execute(scope.Vars)
	if err != nil {
		err = fmt.Errorf("error while exec %s", err.Error())
		return err
	}

	return nil
}
