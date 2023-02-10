package core

import (
	"path/filepath"
)

const CornsIdentifier = "corn"
const CornsNameAttr = CornsIdentifier + "_name"

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

func InitCornVars(base envConfig) map[string]string {
	coreConfig := base.coreConfig
	baseDir := filepath.Join(coreConfig.CurrentDir, coreConfig.CornDir)
	return map[string]string{
		CornsIdentifier + "_dir":   baseDir,
		CornsIdentifier + "_cache": filepath.Join(baseDir, base.CacheDir),
	}
}

func LoadCornEnv(name string, base envConfig) (CornsEnvConfig, error) {
	result := CornsEnvConfig{}
	result.envConfig = base

	baseDir := filepath.Join(base.coreConfig.CurrentDir, base.coreConfig.CornDir)
	loadPath := filepath.Join(baseDir, result.envDirname)
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
