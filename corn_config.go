package cryphtron

import (
	"cryphtron/core"
	"errors"
	"os"
	"path/filepath"
)

func (c *Core) prepareCornConfig(corn core.CornsEnvConfig) (*core.CornsEnvConfig, error) {
	var err error
	currentCornDir := c.Config.CornWithRunningDir()
	if ifFolderNotExists(currentCornDir) {
		_ = os.MkdirAll(currentCornDir, os.ModeDir|0o666)
	}
	scope := c.ComposeRtEnv()
	corn.AppendEnv(scope.Env)
	if !corn.MetaOnly {
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

func (c *Core) execCornConfig(cornConfig core.CornsEnvConfig, isWaiting bool, args ...string) error {
	corn, err := c.prepareCornConfig(cornConfig)
	if err != nil {
		return err
	}
	scope := c.ComposeCornEnv(corn)
	cmd := &corn.Exec
	if !cmd.CanRunning() {
		return errors.New(
			"Cannot running this corn(named:" +
				cornConfig.ID + ") on current platform")
	}
	cmd.Args = append(cmd.Args, args...)
	if !cmd.WithWaiting {
		cmd.WithWaiting = isWaiting
	}
	return c.execCmd(cmd, scope)
}

func (c *Core) ExecCornConfig(cornConfig core.CornsEnvConfig, isWaiting bool, args ...string) error {
	return c.execCornConfig(cornConfig, isWaiting, args...)
}
