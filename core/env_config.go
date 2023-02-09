package core

import (
	"cryphtron/internal"
	"path/filepath"
	"strings"
)

type MirrorType string

func (m MirrorType) String() string {
	if len(m) == 0 {
		return string(MirrorTypDefault)
	}
	return string(m)
}

func (m MirrorType) Convert() MirrorType {
	switch m {
	case MirrorTypNone:
		fallthrough
	case MirrorTypCN:
		return m
	default:
		return MirrorTypDefault
	}
}

const (
	MirrorTypDefault            = MirrorTypNone
	MirrorTypCN      MirrorType = "cn"
	MirrorTypNone    MirrorType = "none"
)

type envConfig struct {
	coreConfig *MainConfig
	internal.ValueScope
	envDirname    string
	ID            string    `toml:"-"`
	CacheDir      string    `toml:"cache_dir"`
	DirName       string    `toml:"dir_name"`
	BootstrapExec []Command `toml:"bootstrap_exec"`
}

func (c envConfig) setCore(coreConfig MainConfig) envConfig {
	c.coreConfig = &coreConfig
	return c
}

func (c *envConfig) setEnvDirname(altEnvDirname ...string) {
	if len(altEnvDirname) > 0 {
		c.envDirname = altEnvDirname[0]
	}
	if c.envDirname == "" {
		c.envDirname = "_env"
	}
}

func (c *envConfig) setCacheDirname(altCacheDirname ...string) {
	if len(altCacheDirname) > 0 {
		c.CacheDir = altCacheDirname[0]
	}
	if c.CacheDir == "" {
		c.CacheDir = "_cache"
	}
}

func (c *envConfig) ExecuteBootstrap() error {
	if len(c.BootstrapExec) == 0 {
		return nil
	}
	c.PrepareScope()
	for _, command := range c.BootstrapExec {
		cmd := command.Prepare().
			SetEnv(c.Env)
		cmd.Env["PATH"] = strings.Replace(
			cmd.Env["PATH"],
			internal.PathPlaceHolder,
			c.Vars["pth_environ"], 1)

		err := cmd.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

type RtEnvConfig struct {
	envConfig
	MirrorEnv  map[MirrorType]map[string]string `toml:"mirror_env"`
	MirrorExec map[MirrorType][]Command         `toml:"mirror_exec"`
}

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

func LoadRtEnv(name string, coreConfig MainConfig, altEnvDirname ...string) (RtEnvConfig, error) {
	c := RtEnvConfig{}
	c.setEnvDirname(altEnvDirname...)
	c.setCacheDirname()
	if c.DirName == "" {
		c.DirName = name
	}
	rtDir := filepath.Join(coreConfig.CurrentDir, coreConfig.RuntimeDir)
	c.AppendVars(map[string]string{
		"rt_dir":   rtDir,
		"rt_cache": filepath.Join(rtDir, c.CacheDir),
		"rt_home":  filepath.Join(rtDir, "_home"),
	})
	loadPath := filepath.Join(coreConfig.CurrentDir, coreConfig.RuntimeDir, c.envDirname)
	err := loadConfigRegular(name, &c, loadPath)
	if err != nil {
		c.setCore(coreConfig)
		return c, err
	}
	for idx := range c.BootstrapExec {
		c.BootstrapExec[idx].Top = &c.ValueScope
	}
	c.envConfig = c.setCore(coreConfig)
	return c, nil
}

type AppEnvConfig struct {
	envConfig
	MetaOnly   bool      `toml:"meta-only"`
	DependApps []string  `toml:"depend-app"`
	ConfigExec []Command `toml:"config-exec"`
	Exec       Command   `toml:"exec"`
}

func (c AppEnvConfig) ExecuteConfig() error {
	c.PrepareScope()
	for _, command := range c.ConfigExec {
		command.WithWaiting = true
		err := command.Prepare().
			SetEnv(c.Env).Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadAppEnv(name string, coreConfig MainConfig, altEnvDirname ...string) (AppEnvConfig, error) {
	result := AppEnvConfig{}
	result.setEnvDirname(altEnvDirname...)

	loadPath := filepath.Join(coreConfig.CurrentDir, coreConfig.AppDir, result.envDirname)
	err := loadConfigRegular(name, &result, loadPath)
	if err != nil {
		result.setCore(coreConfig)
		return result, err
	}

	result.setCacheDirname()
	result.AppendVars(map[string]string{
		"app_dir":   filepath.Join(coreConfig.CurrentDir, coreConfig.AppDir),
		"app_cache": filepath.Join(coreConfig.CurrentDir, coreConfig.AppDir, result.CacheDir),
	})
	if result.DirName == "" {
		result.DirName = name
	}

	for idx := range result.BootstrapExec {
		result.BootstrapExec[idx].Top = &result.ValueScope
	}

	for idx := range result.ConfigExec {
		result.ConfigExec[idx].Top = &result.ValueScope
	}

	result.Exec.Top = &result.ValueScope

	result.envConfig = result.setCore(coreConfig)
	return result, nil
}
