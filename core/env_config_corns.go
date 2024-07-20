package core

import (
	"cryphtron/internal/utils"
	"errors"
	"path/filepath"
)

const CornsIdentifier = "corn"

type CornsEnvConfig struct {
	envConfig
	MetaOnly    bool               `toml:"meta-only"`
	DependCorns []string           `toml:"depend-corns"`
	ConfigExec  []Command          `toml:"config-exec"`
	Exec        Command            `toml:"exec"`
	ExecByPlats map[string]Command `toml:"exec-by-plat"`
}

func (c CornsEnvConfig) ExecuteConfig() error {
	c.PrepareScope()
	for _, command := range c.ConfigExec {
		c0 := command.Prepare().
			SetEnv(c.Env)
		if !c0.CanRunning() {
			continue
		}
		c0.WithNoWait = false
		c0.Top = &c.ValueScope

		c0.AppendEnv(map[string]string{
			"PATH": c.Vars["pth_environ"],
		})

		err := c0.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CornsEnvConfig) initCornVars() {
	coreConfig := c.coreConfig
	vars := make(map[string]string)

	if coreConfig != nil {
		baseDir := coreConfig.CornDir()
		vars = map[string]string{
			CornsIdentifier + "_dir":      coreConfig.CornWithRunningDir(),
			CornsIdentifier + "_cache":    filepath.Join(baseDir, c.CacheDir),
			CornsIdentifier + "_home":     filepath.Join(baseDir, "_home"),
			CornsIdentifier + "_dir_envs": filepath.Join(baseDir, c.envDirname),
		}
		if c.IsCommonPlatform {
			vars[CornsIdentifier+"_dir"] = baseDir
		}
	}
	if c.Top != nil {
		c.Top.AppendVars(vars)
	} else {
		c.AppendVars(vars)
	}

	if len(c.envName) > 0 {
		c.AppendVar(CornsIdentifier+"_name", c.envName)
	}
}

func LoadCornEnv(name string, base envConfig) (CornsEnvConfig, error) {
	result := CornsEnvConfig{}
	if base.coreConfig == nil {
		return result, errors.New("could not loading the env wihout core config")
	}
	result.envConfig = base
	result.ID = CornsIdentifier + "_" + name
	result.envName = name
	result.initCornVars()

	loadPath := filepath.Join(
		result.coreConfig.CornDir(), result.envDirname)
	err := loadConfigRegular(name, &result, loadPath)
	if err != nil {
		return result, err
	}

	if result.DirName == "" {
		result.DirName = name
	}

	for idx := range result.BootstrapExec {
		result.BootstrapExec[idx].Top = &result.ValueScope
	}

	for idx := range result.ConfigExec {
		result.ConfigExec[idx].Top = &result.ValueScope
	}

	if exec, ok := result.ExecByPlats[utils.OS()]; ok {
		result.Exec = exec
	}

	if exec, ok := result.ExecByPlats[utils.Platform()]; ok {
		result.Exec = exec
	}

	result.Exec.Top = &result.ValueScope

	return result, nil
}
