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
	corn, ok := c.CornsEnv[name]
	if !ok {
		return nil,
			errors.New("could not found the " +
				core.CornsIdentifier + " named " + name)
	}

	currentCornDir := c.Config.CornWithRunningDir()
	if ifFolderNotExists(currentCornDir) {
		_ = os.MkdirAll(currentCornDir, os.ModeDir|0o666)
	}

	if !corn.MetaOnly {
		scope := c.ComposeRtEnv()
		corn.AppendEnv(scope.Env)
		bootstrapDir := filepath.Join(currentCornDir, corn.DirName)
		corn.Vars["pth_environ"] = c.Environ["PATH"]
		if ifFolderNotExists(bootstrapDir) {
			_ = os.Mkdir(bootstrapDir, os.ModeDir|0o666)
			corn.AppendEnv(c.Env)
			err = corn.ExecuteBootstrap()
			if err != nil {
				_ = os.RemoveAll(bootstrapDir)
				newErr := errors.New(
					"error while bootstrap " +
						core.CornsIdentifier + "[" + corn.ID + "]:")
				return nil, errors.Join(newErr, err)
			}
		}
	}

	err = corn.ExecuteConfig()
	if err != nil {
		newErr := errors.New("error while configure " +
			core.CornsIdentifier + "[" + corn.ID + "]:")
		return nil, errors.Join(newErr, err)
	}

	for _, depend := range corn.DependCorns {
		var cfg *core.CornsEnvConfig
		cfg, err = c.prepareCorn(depend)
		if err != nil {
			return nil, err
		}
		corn.AppendVars(cfg.Vars)
	}
	return &corn, nil
}

func (c *Core) execCorn(name string, isWaiting bool, args ...string) error {
	corn, err := c.prepareCorn(name)
	if err != nil {
		return err
	}
	scope := c.ComposeCornEnv(corn)
	cmd := &corn.Exec
	if !cmd.CanRunning() {
		return errors.New(
			"Cannot running this corn(named:" +
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
