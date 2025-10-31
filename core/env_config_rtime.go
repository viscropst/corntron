package core

import (
	"corntron/internal"
	"path/filepath"
)

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

func (c RtEnvConfig) Copy(src ...RtEnvConfig) RtEnvConfig {
	result := RtEnvConfig{}
	if len(src) > 0 {
		tmp := src[0]
		result.envConfig = tmp.envConfig.Copy()
		result.MirrorEnv = tmp.MirrorEnv
		result.MirrorExec = tmp.MirrorExec
	} else {
		result.envConfig = c.envConfig.Copy()
		result.MirrorEnv = c.MirrorEnv
		result.MirrorExec = c.MirrorExec
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
	for k, v := range c.Env {
		if k == "PATH" {
			continue
		}
		env[k] = v
	}
	for _, command := range mirrorExec {
		c0 := command.SetEnv(env)
		c0.EnvPath = c.EnvPath
		if !c0.CanRunning() {
			continue
		}
		err := c.executeCommand(command)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RtEnvConfig) initRtVars() {
	vars := make(map[string]string)

	if coreConfig := c.coreConfig; coreConfig != nil {
		baseDir := coreConfig.RuntimeDir()
		vars = map[string]string{
			RtIdentifier + "_dir":   coreConfig.RtWithRunningDir(),
			RtIdentifier + "_cache": filepath.Join(baseDir, c.CacheDir),
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
	if base.coreConfig == nil {
		return result, internal.Error("could not loading the env wihout core config")
	}
	result.envConfig = base
	result.ID = RtIdentifier + "_" + name
	result.envName = name
	if result.DirName == "" {
		result.DirName = name
	}
	result.initRtVars()
	loadPath := filepath.Join(
		result.coreConfig.RuntimeDir(), result.envDirname)
	err := loadConfigRegular(name, &result, loadPath)
	if err != nil {
		return result, err
	}
	for idx := range result.BootstrapExec {
		result.BootstrapExec[idx].Top = &result.ValueScope
	}
	return result, nil
}
