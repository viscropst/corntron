package core

import (
	"corntron/internal"
	"errors"
	"path/filepath"
	"strings"
)

const CornsIdentifier = "corn"

type CornsEnvConfig struct {
	envConfig
	MetaOnly    bool               `toml:"meta_only"`
	DependCorns []string           `toml:"depend_corns"`
	Exec        Command            `toml:"exec"`
	ExecByPlats map[string]Command `toml:"exec_by_plat"`
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

func LoadCornEnv(name string, base envConfig, altPath ...string) (CornsEnvConfig, error) {
	result := CornsEnvConfig{}
	if base.coreConfig == nil {
		return result, errors.New("could not loading the env wihout core config")
	}
	result.envConfig = base
	result.ID = CornsIdentifier + "_" + name
	result.envName = name
	if strings.HasSuffix(name, CornConfigExt) {
		result.envName = name[:len(name)-len(CornConfigExt)]
	}
	result.initCornVars()

	loadPath := filepath.Join(
		result.coreConfig.CornDir(), result.envDirname)
	if len(altPath) > 0 {
		loadPath = altPath[0]
	}
	err := loadConfigRegular(name, &result, loadPath)
	if err != nil {
		return result, err
	}

	if result.DirName == "" {
		result.DirName = result.envName
	}

	for idx := range result.BootstrapExec {
		result.BootstrapExec[idx].Top = &result.ValueScope
	}

	for idx := range result.ConfigExec {
		result.ConfigExec[idx].Top = &result.ValueScope
	}

	if exec, ok := result.ExecByPlats[internal.OS()]; ok {
		result.Exec = setExec(result, exec)
	}

	if exec, ok := result.ExecByPlats[internal.Platform()]; ok {
		result.Exec = setExec(result, exec)
	}

	result.Exec.Top = &result.ValueScope

	return result, nil
}

func setExec(conf CornsEnvConfig, src Command) Command {
	result := conf.Exec
	if src.WithNoWaiting {
		result.WithNoWaiting = src.WithNoWaiting
	}
	if src.WithEnviron {
		result.WithEnviron = src.WithEnviron
	}
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
