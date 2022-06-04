package cryphtron

import (
	"cryphtron/core"
	"cryphtron/internal"
	"fmt"
	"io/fs"
	"strings"
)

func LoadCoreConfig(altBases ...string) core.CoreConfig {
	return core.LoadCoreConfig(altBases...)
}

type Core struct {
	internal.Core
	Config      core.CoreConfig
	Environ     map[string]string
	AppsEnv     map[string]core.AppEnvConfig
	RuntimesEnv []core.RtEnvConfig
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

func (c Core) ComposeAppEnv() *internal.ValueScope {
	c.Prepare()
	for _, config := range c.AppsEnv {
		c.ValueScope.AppendEnv(config.Env)
	}
	return c.ValueScope
}

func (c Core) ProcessRtMirror() error {
	mirrorType := c.Config.UnwrapMirrorType()
	c.Prepare()
	for _, config := range c.RuntimesEnv {
		config.PrepareScope()
		config.AppendEnv(c.Env)
		err := config.ExecuteMirrors(mirrorType)
		if err != nil {
			return fmt.Errorf("mirror[%s]:%s", mirrorType, err.Error())
		}
	}
	return nil
}

func LoadCore(coreConfig core.CoreConfig, altNames ...string) (Core, error) {
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
			env, envErr := core.LoadRtEnv(configName, coreConfig, envDirName)
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
			env, envErr := core.LoadAppEnv(configName, coreConfig, envDirName)
			if envErr != nil {
				return envErr
			}
			env.ID = "app_" + configName
			env.Top = result.ValueScope
			result.AppsEnv[configName] = env
			return nil
		},
		coreConfig.AppDir, envDirName)
	if err != nil {
		return result, err
	}

	return result, nil
}
