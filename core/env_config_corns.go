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
	Exec        Command            `toml:"exec"`
	ExecByPlats map[string]Command `toml:"exec-by-plat"`
}

func (c *CornsEnvConfig) initCornVars() {
	coreConfig := c.coreConfig
	vars := make(map[string]string)

	if coreConfig != nil {
		baseDir := coreConfig.CornDir()
		vars = map[string]string{
			CornsIdentifier + "_dir":      coreConfig.CornWithRunningDir(),
			CornsIdentifier + "_cache":    filepath.Join(baseDir, c.CacheDir),
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
		result.Exec = setExec(result, exec)
	}

	if exec, ok := result.ExecByPlats[utils.Platform()]; ok {
		result.Exec = setExec(result, exec)
	}

	result.Exec.Top = &result.ValueScope

	return result, nil
}

func setExec(conf CornsEnvConfig, src Command) Command {
	result := conf.Exec
	result.WithWaiting = src.WithWaiting
	result.WithEnviron = src.WithEnviron
	if len(src.Exec) > 0 && len(src.Args) > 0 {
		return src
	}
	if len(src.Exec) > 0 {
		result.Exec = src.Exec
	}
	if len(src.Args) > 0 {
		result.Args = src.Args
	}
	if arr := result.ArgStr.ToArray(); len(arr) > 0 {
		result.Args = append(result.Args, arr...)
	}
	if len(src.Vars) > 0 {
		result.AppendVarsByNew(src.Vars)
	}
	if len(src.Env) > 0 {
		_ = result.AppendEnvs(src.Env)
	}
	return result
}

func (c CornsEnvConfig) Copy(src ...CornsEnvConfig) CornsEnvConfig {
	result := CornsEnvConfig{}
	if len(src) > 0 {
		tmp := src[0]
		result.envConfig = tmp.envConfig.Copy()
		result.MetaOnly = tmp.MetaOnly
		result.DependCorns = tmp.DependCorns
		result.Exec = tmp.Exec
		result.ExecByPlats = tmp.ExecByPlats
	} else {
		result.envConfig = c.envConfig.Copy()
		result.MetaOnly = c.MetaOnly
		result.DependCorns = c.DependCorns
		result.Exec = c.Exec
		result.ExecByPlats = c.ExecByPlats
	}
	return result
}
