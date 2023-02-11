package core

import (
	"path/filepath"
)

const CornsIdentifier = "corn"

func (c *envConfig) cornsDir() string {
	return filepath.Join(
		c.coreConfig.CurrentDir,
		c.coreConfig.CornDir)
}

type CornsEnvConfig struct {
	envConfig
	MetaOnly    bool      `toml:"meta-only"`
	DependCorns []string  `toml:"depend-corns"`
	ConfigExec  []Command `toml:"config-exec"`
	Exec        Command   `toml:"exec"`
}

func (c CornsEnvConfig) ExecuteConfig() error {
	c.PrepareScope()
	for _, command := range c.ConfigExec {
		command.WithWaiting = true
		err := command.Prepare().
			SetEnv(c.Env).Execute()
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
		baseDir := c.cornsDir()
		vars = map[string]string{
			CornsIdentifier + "_dir":   baseDir,
			CornsIdentifier + "_cache": filepath.Join(baseDir, c.CacheDir),
			CornsIdentifier + "_home":  filepath.Join(baseDir, "_home"),
		}
	}
	if c.Top != nil {
		c.Top.AppendVars(vars)
	} else {
		c.AppendVars(vars)
	}
	if len(c.envName) > 0 {
		c.Vars[CornsIdentifier+"_name"] = c.envName
	}
}

func LoadCornEnv(name string, base envConfig) (CornsEnvConfig, error) {
	result := CornsEnvConfig{}
	result.envConfig = base
	result.ID = CornsIdentifier + "_" + name
	result.envName = name
	result.initCornVars()

	loadPath := filepath.Join(result.cornsDir(), result.envDirname)
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

	result.Exec.Top = &result.ValueScope

	return result, nil
}
