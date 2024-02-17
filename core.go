package cryphtron

import (
	"cryphtron/core"
	"cryphtron/internal"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func LoadCoreConfigWithRuningBase(runningBase string, altBases ...string) core.MainConfig {
	result := LoadCoreConfig(altBases...)
	if len(runningBase) > 0 && !filepath.IsAbs(runningBase) {
		result.RunningDir = filepath.Join(result.CurrentDir, runningBase)
		return result
	}
	if len(runningBase) > 0 && filepath.IsAbs(runningBase) {
		result.RunningDir = runningBase
		return result
	}
	return result
}

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
			tmpEnv := core.BaseEnv(coreConfig, envDirName)
			tmpEnv.Top = result.ValueScope
			env, envErr := core.LoadRtEnv(configName, tmpEnv)
			if envErr != nil {
				return envErr
			}
			result.RuntimesEnv = append(result.RuntimesEnv, env)
			return nil
		},
		coreConfig.RuntimeDirName, envDirName)
	if err != nil {
		return result, err
	}

	err = coreConfig.FsWalk(
		func(path string, info fs.FileInfo, err error) error {
			if !coreConfig.WithCorn {
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
			tmpEnv := core.BaseEnv(coreConfig, envDirName)
			tmpEnv.Top = result.ValueScope
			env, envErr := core.LoadCornEnv(configName, tmpEnv)
			if envErr != nil {
				return envErr
			}
			result.CornsEnv[configName] = env
			return nil
		},
		coreConfig.CornDirName, envDirName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func ifFolderNotExists(path string) bool {
	stat, _ := os.Stat(path)
	return stat == nil || (stat != nil && !stat.IsDir())
}
