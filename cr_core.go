package cryphtron

import (
	"cryphtron/internal"
	"io/fs"
	"strings"
)

type Core struct {
	internal.Core
	Config      CoreConfig
	Environ     map[string]string
	EditorsEnv  []EditEnvConfig
	RuntimesEnv []RtEnvConfig
}

func (c Core) ComposeRtEnv() *internal.ValueScope {
	c.Prepare()
	for _, config := range c.RuntimesEnv {
		config.PrepareScope()
		mirrorEnv := config.UnwrapMirrorsEnv(c.Config.UnwrapMirrorType())
		for k, v := range mirrorEnv {
			mirrorEnv[k] = config.Expand(v)
		}
		c.AppendEnv(config.Env).AppendEnv(mirrorEnv)
	}

	return c.ValueScope
}

func (c Core) ComposeEdEnv() *internal.ValueScope {
	c.Prepare()
	for _, config := range c.EditorsEnv {
		c.ValueScope.AppendEnv(config.Env)
	}
	return c.ValueScope
}

func LoadCore(coreConfig CoreConfig, altNames ...string) (Core, error) {
	result := Core{
		Config: coreConfig,
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
			env, envErr := LoadRtEnv(configName, coreConfig, envDirName)
			if envErr != nil {
				return envErr
			}
			env.ID = "rt_" + configName
			env.Top = result.ValueScope
			result.RuntimesEnv = append(result.RuntimesEnv, env)
			return nil
		},
		coreConfig.RuntimeDir, envDirName)
	if err != nil {
		return result, err
	}

	err = coreConfig.FsWalk(
		func(path string, info fs.FileInfo, err error) error {
			if !coreConfig.WithEditor {
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
			env, envErr := LoadEditEnv(configName, coreConfig, envDirName)
			if envErr != nil {
				return envErr
			}
			env.ID = "ed_" + configName
			env.Top = result.ValueScope
			result.EditorsEnv = append(result.EditorsEnv, env)
			return nil
		},
		coreConfig.EditorDir, envDirName)
	if err != nil {
		return result, err
	}

	return result, nil
}
