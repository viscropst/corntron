package cryphtron

import (
	"cryphtron/core"
	"cryphtron/internal"
	"io/fs"
	"strings"
)

func LoadCoreConfig(altBases ...string) core.MainConfig {
	return core.LoadCoreConfig(altBases...)
}

type Core struct {
	internal.Core
	Config      core.MainConfig
	CornsEnv    map[string]core.CornsEnvConfig
	RuntimesEnv []core.RtEnvConfig
}

func LoadCore(coreConfig core.MainConfig, altNames ...string) (Core, error) {
	result := Core{
		Config:   coreConfig,
		CornsEnv: make(map[string]core.CornsEnvConfig),
	}

	result.ValueScope = &internal.ValueScope{
		Env:  make(map[string]string),
		Vars: make(map[string]string),
	}

	envDirName := "_env"
	if len(altNames) > 0 {
		envDirName = altNames[0]
	}

	baseEnv := core.BaseEnv(coreConfig, envDirName)

	baseEnv.AppendVar("base_dir", coreConfig.CurrentDir)
	baseEnv.AppendVar(core.CornsIdentifier+"_dirname", coreConfig.CornDir)
	baseEnv.AppendVar(core.RtIdentifier+"_dirname", coreConfig.RuntimeDir)

	result.Prepare()
	err := coreConfig.FsWalk(
		func(path string, info fs.FileInfo, err error) error {
			if info == nil {
				return nil
			}
			if !info.Mode().IsRegular() ||
				(!info.IsDir() && !strings.HasSuffix(info.Name(), ".toml")) {
				return nil
			}

			configName := strings.TrimSuffix(info.Name(), ".toml")
			tmpEnv := baseEnv
			tmpEnv.Top = result.ValueScope
			env, envErr := core.LoadRtEnv(configName, tmpEnv)
			if envErr != nil {
				return envErr
			}
			result.RuntimesEnv = append(result.RuntimesEnv, env)
			return nil
		},
		coreConfig.RuntimeDir, envDirName)
	if err != nil {
		return result, err
	}

	err = coreConfig.FsWalk(
		func(path string, info fs.FileInfo, err error) error {
			if !coreConfig.WithApp {
				return nil
			}
			if info == nil {
				return nil
			}
			if !info.Mode().IsRegular() ||
				(!info.IsDir() && !strings.HasSuffix(info.Name(), ".toml")) {
				return nil
			}
			configName := strings.TrimSuffix(info.Name(), ".toml")
			tmpEnv := baseEnv
			tmpEnv.Top = result.ValueScope
			env, envErr := core.LoadCornEnv(configName, tmpEnv)
			if envErr != nil {
				return envErr
			}
			result.CornsEnv[configName] = env
			return nil
		},
		coreConfig.CornDir, envDirName)
	if err != nil {
		return result, err
	}

	return result, nil
}
