package cryphtron

import (
	"cryphtron/core"
	"cryphtron/internal/utils"
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
	tmpCorn := corn.Copy()
	tmpCorn.RePrepareScope()
	tmpCorn.ModifyEnv("PATH",
		utils.AppendToPathList(corn.Env["PATH"], scope.Env["PATH"]))
	tmpCorn.AppendEnvs(scope.Env)
	if !corn.MetaOnly {
		bootstrapDir := filepath.Join(currentCornDir, corn.DirName)
		tmpCorn.Vars["pth_environ"] = c.Environ["PATH"]
		if ifFolderNotExists(bootstrapDir) {
			_ = os.Mkdir(bootstrapDir, os.ModeDir|0o666)
			tmpCorn.AppendEnvs(c.Env)
			err = tmpCorn.ExecuteBootstrap()
			if err != nil {
				_ = os.RemoveAll(bootstrapDir)
				newErr := errors.New(
					"error while bootstrap " +
						core.CornsIdentifier + "[" + corn.ID + "]:")
				return nil, errors.Join(newErr, err)
			}
		}
	}

	err = tmpCorn.ExecuteConfig()
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
		tmpCorn.AppendVars(cfg.Vars)
	}
	return &tmpCorn, nil
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
	cmd.GlobalWaiting = isWaiting
	return c.execCmd(cmd, scope)
}

func (c *Core) ExecCornConfig(cornConfig core.CornsEnvConfig, isWaiting bool, args ...string) error {
	return c.execCornConfig(cornConfig, isWaiting, args...)
}
