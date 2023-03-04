package cryphtron

import (
	"cryphtron/core"
	"cryphtron/internal"
	"errors"
	"os"
	"path/filepath"
)

func (c Core) ComposeCornEnv(corn *core.CornsEnvConfig) *internal.ValueScope {
	c.Prepare()
	if corn == nil {
		return c.ValueScope
	}
	for _, depends := range corn.DependCorns {
		config := c.CornsEnv[depends]
		config.AppendVars(corn.Vars)
		config.PrepareScope()
		corn.AppendEnv(config.Env)
		corn.AppendVars(config.Vars)
	}
	return &corn.ValueScope
}

func (c *Core) prepareCorn(name string) (*core.CornsEnvConfig, error) {
	var err error
	app, ok := c.CornsEnv[name]
	if !ok {
		return nil,
			errors.New("could not found the " +
				core.CornsIdentifier + " named " + name)
	}

	if !app.MetaOnly {
		scope := c.ComposeRtEnv()
		app.AppendEnv(scope.Env)
		bootstrapDir := c.Config.CornsPath()
		bootstrapDir = filepath.Join(bootstrapDir, app.DirName)
		stat, _ := os.Stat(bootstrapDir)
		app.Vars["pth_environ"] = c.Environ["PATH"]
		if stat == nil || (stat != nil && !stat.IsDir()) {
			_ = os.Mkdir(bootstrapDir, os.ModeDir|0o666)
			app.AppendEnv(c.Env)
			err = app.ExecuteBootstrap()
			if err != nil {
				_ = os.RemoveAll(bootstrapDir)
				newErr := errors.New(
					"error while bootstrap " +
						core.CornsIdentifier + "[" + app.ID + "]:")
				return nil, errors.Join(newErr, err)
			}
		}
	}

	err = app.ExecuteConfig()
	if err != nil {
		newErr := errors.New("error while configure " +
			core.CornsIdentifier + "[" + app.ID + "]:")
		return nil, errors.Join(newErr, err)
	}

	for _, depend := range app.DependCorns {
		var cfg *core.CornsEnvConfig
		cfg, err = c.prepareCorn(depend)
		if err != nil {
			return nil, err
		}
		app.AppendVars(cfg.Vars)
	}
	return &app, nil
}

func (c *Core) execCorn(name string, isWaiting bool, args ...string) error {
	app, err := c.prepareCorn(name)
	if err != nil {
		return err
	}
	scope := c.ComposeCornEnv(app)
	cmd := &app.Exec
	if !cmd.CanRunning() {
		return errors.New(
			"Cannot running this app(named:" +
				name + ") on current platform")
	}
	cmd.Args = append(cmd.Args, args...)
	origWait := cmd.WithNoWait
	if isWaiting {
		cmd.WithNoWait = !isWaiting
	}
	if origWait == isWaiting {
		cmd.WithNoWait = origWait
	}
	return c.execCmd(cmd, scope)
}

func (c *Core) ExecCorn(name string, isWaiting bool, args ...string) error {
	return c.execCorn(name, isWaiting, args...)
}
