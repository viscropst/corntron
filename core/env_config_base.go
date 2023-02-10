package core

import (
	"cryphtron/internal"
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

		cmd.WithWaiting = true
		err := cmd.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func BaseEnv(coreConfig MainConfig, altEnvDirname ...string) envConfig {
	tmp := envConfig{}
	tmp.setCacheDirname()
	tmp.setEnvDirname(altEnvDirname...)
	tmp = tmp.setCore(coreConfig)
	return tmp
}
