package core

import (
	"cryphtron/internal"
	"path/filepath"
	"strings"
)

func (c *envConfig) runtimesDir() string {
	return c.coreConfig.RuntimesPath()
}

type RtEnvConfig struct {
	envConfig
	MirrorEnv  map[MirrorType]map[string]string `toml:"mirror_env"`
	MirrorExec map[MirrorType][]Command         `toml:"mirror_exec"`
}

const RtIdentifier = "rt"

func (c *RtEnvConfig) UnwrapMirrorsEnv(key MirrorType) map[string]string {
	var result = make(map[string]string)
	for k, v := range c.MirrorEnv[key] {
		result[k] = v
	}
	return result
}

func (c *RtEnvConfig) ExecuteMirrors(mirrorType MirrorType) error {
	mirrorExec, ok := c.MirrorExec[mirrorType]
	if !ok {
		return nil
	}
	c.PrepareScope()
	env := make(map[string]string)
	env["PATH"] = c.Env["PATH"]
	for k, v := range c.Env {
		if k == "PATH" {
			continue
		}
		env[k] = v
	}
	for _, command := range mirrorExec {
		c0 := command.Prepare().
			SetEnv(env)
		if !c0.CanRunning() {
			continue
		}
		c0.WithNoWait = false
		c0.Top = &c.ValueScope
		c0.Env["PATH"] = strings.Replace(
			c0.Env["PATH"], internal.PathPlaceHolder, c.Vars["pth_environ"], 1)
		err := c0.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RtEnvConfig) initRtVars() {
	vars := make(map[string]string)

	if coreConfig := c.coreConfig; coreConfig != nil {
		baseDir := coreConfig.RuntimesPath()
		vars = map[string]string{
			RtIdentifier + "_dir":   coreConfig.RtWithRunningDir(),
			RtIdentifier + "_cache": filepath.Join(baseDir, c.CacheDir),
			RtIdentifier + "_home":  filepath.Join(baseDir, "_home"),
		}
	}
	if c.Top != nil {
		c.Top.AppendVars(vars)
	} else {
		c.AppendVars(vars)
	}

	if len(c.envName) > 0 {
		c.AppendVar(RtIdentifier+"_name", c.envName)
	}

}

func LoadRtEnv(name string, base envConfig) (RtEnvConfig, error) {
	result := RtEnvConfig{}
	result.envConfig = base
	result.ID = RtIdentifier + "_" + name
	result.envName = name
	if result.DirName == "" {
		result.DirName = name
	}
	result.initRtVars()
	loadPath := filepath.Join(
		result.runtimesDir(), result.envDirname)
	err := loadConfigRegular(name, &result, loadPath)
	if err != nil {
		return result, err
	}
	for idx := range result.BootstrapExec {
		result.BootstrapExec[idx].Top = &result.ValueScope
	}
	return result, nil
}
