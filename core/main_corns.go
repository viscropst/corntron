package core

import (
	"corntron/internal"
	"path/filepath"
)

func (c *Core) prepareCornConfig(corn CornsEnvConfig) (*CornsEnvConfig, error) {
	var err error
	currentCornDir := c.Config.CornWithRunningDir()
	if internal.IfFolderNotExists(currentCornDir) {
		_ = internal.Mkdir(currentCornDir)
	}
	scope := c.ComposeRtEnv()
	tmpCorn := corn.Copy()
	tmpCorn.RePrepareScope()
	tmpCorn.EnvPath = corn.EnvPath.AppendList(scope.EnvPath)
	tmpCorn.AppendEnvs(scope.Env)
	if !corn.MetaOnly {
		bootstrapDir := filepath.Join(currentCornDir, corn.DirName)
		tmpCorn.Vars["pth_environ"] = c.Environ["PATH"]
		if internal.IfFolderNotExists(bootstrapDir) {
			_ = internal.Mkdir(bootstrapDir)
			tmpCorn.AppendEnvs(c.Env)
			err = tmpCorn.ExecuteBootstrap()
			if err != nil {
				_ = internal.Remove(bootstrapDir)
				newErr := internal.Error(
					"error while bootstrap ",
					CornsIdentifier, "[", corn.ID, "]:", err.Error())
				return nil, newErr
			}
		}
	}

	err = tmpCorn.ExecuteConfig()
	if err != nil {
		newErr := internal.Error("error while configure ",
			CornsIdentifier, "[", corn.ID, "]:", err.Error())
		return nil, newErr
	}

	for _, depend := range corn.DependCorns {
		var cfg *CornsEnvConfig
		cfg, err = c.prepareCorn(depend)
		if err != nil {
			return nil, err
		}
		tmpCorn.AppendVars(cfg.Vars)
	}
	return &tmpCorn, nil
}

func (c *Core) prepareCorn(name string) (*CornsEnvConfig, error) {
	corn, ok := c.CornsEnv[name]
	if !ok {
		return nil, internal.Error("could not found the ",
			CornsIdentifier, " named ", name)
	}
	return c.prepareCornConfig(corn)
}

func (c *Core) ExecCornConfig(cornConfig CornsEnvConfig, isWaiting bool, args ...string) error {
	corn, err := c.prepareCornConfig(cornConfig)
	if err != nil {
		return err
	}
	scope := c.ComposeCornEnv(corn)
	cmd := &corn.Exec
	if !cmd.CanRunning() {
		return internal.Error(
			"Cannot running this corn(named:",
			cornConfig.ID, ") on current platform")
	}
	cmd.Args = append(cmd.Args, args...)
	cmd.GlobalWaiting = isWaiting
	return c.execCmd(cmd, scope)
}

func (c *Core) ExecCorn(name string, isWaiting bool, args ...string) error {
	corn, ok := c.CornsEnv[name]
	if !ok {
		return internal.Error("could not found the ",
			CornsIdentifier, " named ", name)
	}
	return c.ExecCornConfig(corn, isWaiting, args...)
}

func (c Core) ComposeCornEnv(corn *CornsEnvConfig) *ValueScope {
	c.Prepare()
	if corn == nil {
		return c.ValueScope
	}
	for _, depends := range corn.DependCorns {
		config := c.CornsEnv[depends]
		config.AppendVars(corn.Vars)
		config.RePrepareScope()
		corn.AppendEnvs(config.Env)
		corn.AppendVars(config.Vars)
	}
	return &corn.ValueScope
}
