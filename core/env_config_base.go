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
	envDirname       string
	envName          string
	ID               string    `toml:"-"`
	CacheDir         string    `toml:"cache_dir"`
	DirName          string    `toml:"dir_name"`
	BootstrapExec    []Command `toml:"bootstrap_exec"`
	IsCommonPlatform bool      `toml:"is_common_platform"`
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
	bootstraps := make([]Command, 0)
	for _, v := range c.BootstrapExec {
		if !v.CanRunning() {
			continue
		}
		bootstraps = append(bootstraps, v)
	}
	for _, command := range bootstraps {
		err := c.executeCommand(command)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *envConfig) executeCommand(command Command) error {
	cmd := command.Prepare().
		SetEnv(c.Env)
	cmd.Env["PATH"] = strings.Replace(
		cmd.Env["PATH"],
		internal.PathPlaceHolder,
		c.Vars["pth_environ"], 1)

	cmd.WithWaiting = true
	return cmd.Execute()
}

func BaseEnv(coreConfig MainConfig, altEnvDirname ...string) envConfig {
	tmp := envConfig{}
	tmp.setCacheDirname()
	tmp.setEnvDirname(altEnvDirname...)
	tmp = tmp.setCore(coreConfig)
	tmp.AppendVar("base_dir", coreConfig.CurrentDir)
	tmp.AppendVar("base_platform_dir", coreConfig.CurrentPlatformDir)
	if !coreConfig.IsUserProfile() {
		tmp.AppendVar("profile", coreConfig.ProfileDir)
	}
	tmp.AppendVar(CornsIdentifier+"_dirname", coreConfig.CornDirName)
	tmp.AppendVar(CornsIdentifier+"_home", filepath.Join(coreConfig.CornDir(), "_home"))
	tmp.AppendVar(RtIdentifier+"_dirname", coreConfig.RuntimeDirName)
	tmp.AppendVar(RtIdentifier+"_home", filepath.Join(coreConfig.RuntimeDir(), "_home"))
	return tmp
}
