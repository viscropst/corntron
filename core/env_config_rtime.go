package core

import (
	"cryphtron/internal"
	"path/filepath"
	"strings"
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
		c0.WithWaiting = true
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

func InitRtVars(base envConfig) map[string]string {
	baseDir := filepath.Join(base.coreConfig.CurrentDir, base.coreConfig.RuntimeDir)
	return map[string]string{
		RtIdentifier + "_dir":   baseDir,
		RtIdentifier + "_cache": filepath.Join(baseDir, base.CacheDir),
		RtIdentifier + "_home":  filepath.Join(baseDir, "_home"),
	}
}

func LoadRtEnv(name string, base envConfig) (RtEnvConfig, error) {
	c := RtEnvConfig{}
	c.envConfig = base
	if c.DirName == "" {
		c.DirName = name
	}

	loadPath := filepath.Join(c.coreConfig.CurrentDir, c.coreConfig.RuntimeDir, c.envDirname)
	err := loadConfigRegular(name, &c, loadPath)
	if err != nil {
		return c, err
	}
	for idx := range c.BootstrapExec {
		c.BootstrapExec[idx].Top = &c.ValueScope
	}
	return c, nil
}
