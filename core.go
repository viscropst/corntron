package cryphtron

import (
	"cryphtron/core"
	"cryphtron/internal"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

func (c Core) ComposeRtEnv() *internal.ValueScope {
	c.Prepare()
	for _, config := range c.RuntimesEnv {
		config.PrepareScope()
		mirrorEnv := config.UnwrapMirrorsEnv(c.Config.UnwrapMirrorType())
		for k, v := range mirrorEnv {
			mirrorEnv[k] = config.Expand(v)
		}
		pthArr := strings.SplitN(config.Env["PATH"], string(os.PathListSeparator), 2)
		if len(pthArr) > 1 && pthArr[0] == internal.PathPlaceHolder {
			idxAfter := strings.Index(internal.PathPlaceHolder, pthArr[0])
			config.Env["PATH"] = config.Env["PATH"][idxAfter+len(pthArr[0])+1:]
		}
		c.AppendEnv(config.Env).AppendEnv(mirrorEnv)
	}

	return c.ValueScope
}

func (c Core) ComposeAppEnv(corn *core.CornsEnvConfig) *internal.ValueScope {
	c.Prepare()
	if corn == nil {
		return c.ValueScope
	}
	for _, depends := range corn.DependCorns {
		config := c.CornsEnv[depends]
		config.AppendVars(corn.Vars)
		config.PrepareScope()
		corn.AppendEnv(config.Env)
		corn.AppendVars(config.Vars)
	}
	return &corn.ValueScope
}

func (c Core) ProcessRtMirror() error {
	mirrorType := c.Config.UnwrapMirrorType()
	c.Prepare()
	for _, config := range c.RuntimesEnv {
		config.PrepareScope()
		config.AppendEnv(c.Env)
		config.Vars["pth_environ"] = c.Environ["PATH"]
		err := config.ExecuteMirrors(mirrorType)
		if err != nil {
			return fmt.Errorf("mirror[%s]:%s", mirrorType, err.Error())
		}
	}
	return nil
}

func (c Core) ProcessRtBootstrap() error {
	c.Prepare()
	currentRtDir := filepath.Join(
		c.Config.CurrentDir, c.Config.RuntimeDir)
	for _, config := range c.RuntimesEnv {
		bootstrapDir := filepath.Join(currentRtDir,
			config.DirName)
		stat, _ := os.Stat(bootstrapDir)
		if stat != nil || (stat != nil && stat.IsDir()) {

			continue
		}
		_ = os.Mkdir(bootstrapDir, os.ModeDir|0o666)
		config.PrepareScope()
		config.AppendEnv(c.Env)
		config.Vars["pth_environ"] = c.Environ["PATH"]
		err := config.ExecuteBootstrap()
		if err != nil {
			_ = os.RemoveAll(bootstrapDir)
			return fmt.Errorf("bootstrsp[%s]:%s",
				config.DirName, err.Error())
		}
	}
	return nil
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

	result.Prepare()
	result.AppendVars(core.InitRtVars(baseEnv))
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
			env, envErr := core.LoadRtEnv(configName, baseEnv)
			if envErr != nil {
				return envErr
			}
			env.ID = core.RtIdentifier + "_" + configName
			env.Top = result.ValueScope
			result.RuntimesEnv = append(result.RuntimesEnv, env)
			return nil
		},
		coreConfig.RuntimeDir, envDirName)
	if err != nil {
		return result, err
	}

	result.AppendVars(core.InitCornVars(baseEnv))
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
			env, envErr := core.LoadCornEnv(configName, baseEnv)
			if envErr != nil {
				return envErr
			}
			env.ID = core.CornsIdentifier + "_" + configName
			env.Top = result.ValueScope
			env.Vars[core.CornsNameAttr] = configName
			result.CornsEnv[configName] = env
			return nil
		},
		coreConfig.CornDir, envDirName)
	if err != nil {
		return result, err
	}

	return result, nil
}
