package actions

import (
	"cryphtron"
	"cryphtron/cmd/cptron"
	"cryphtron/core"
	"cryphtron/internal"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type execApp struct {
	appName string
	args    []string
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
	if len(os.Args) > argCmdIdx+1 {
		c.args = os.Args[argCmdIdx+1:]
	}
	return nil
}

func (c *execApp) BeforeCore(coreConfig *core.MainConfig) error {
	coreConfig.WithApp = true
	return nil
}

func (c *execApp) Exec(core *cryphtron.Core) error {
	err := core.ProcessRtMirror()
	if err != nil {
		err = fmt.Errorf("error while processing mirror %s", err.Error())
		return err
	}
	return c.execApp(c.appName, core)
}

func (c *execApp) prepareApp(name string, coreObj *cryphtron.Core) (*core.AppEnvConfig, error) {
	var err error
	app, ok := coreObj.AppsEnv[name]
	if !ok {
		return nil, fmt.Errorf("could not found the app named %s", c.appName)
	}

	if !app.MetaOnly {
		scope := coreObj.ComposeRtEnv()
		app.AppendEnv(scope.Env)
		bootStrapDir := filepath.Join(coreObj.Config.CurrentDir, coreObj.Config.AppDir)
		bootStrapDir = filepath.Join(bootStrapDir, app.DirName)
		stat, _ := os.Stat(bootStrapDir)
		if stat == nil || (stat != nil && !stat.IsDir()) {
			_ = os.Mkdir(bootStrapDir, os.ModeDir|0o666)
			err = app.ExecuteBootstrap()
			if err != nil {
				err = fmt.Errorf("error while bootstrap app["+app.ID+"]:%S", err.Error())
				return nil, err
			}
		}
	}

	err = app.ExecuteConfig()
	if err != nil {
		err = fmt.Errorf("error while configure app["+app.ID+"]:%S", err.Error())
		return nil, err
	}

	for _, depend := range app.DependApps {
		var cfg *core.AppEnvConfig
		cfg, err = c.prepareApp(depend, coreObj)
		if err != nil {
			return nil, err
		}
		app.AppendVars(cfg.Vars)
	}
	return &app, nil
}

func (c *execApp) execApp(name string, core *cryphtron.Core) error {
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
	err = cmd.SetEnv(app.Env).ExecuteNoWait(scope.Vars)
	if err != nil {
		err = fmt.Errorf("error while exec %s", err.Error())
		return err
	}

	return nil
}
