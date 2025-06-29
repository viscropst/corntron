package core

import (
	"path/filepath"
)

type MirrorType string

func (m MirrorType) String() string {
	if len(m) == 0 {
		return string(MirrorTypNone)
	}
	return string(m)
}

func (m MirrorType) Convert() MirrorType {
	if len(m) == 0 {
		return MirrorTypDefault
	}
	switch m {
	case MirrorTypNone:
		return m
	case MirrorTypCN:
		return m
	default:
		return MirrorTypCustomized
	}
}

func (m MirrorType) ConvertWithTypes(types ...MirrorType) MirrorType {
	typ := m.Convert()
	if typ != MirrorTypCustomized {
		return typ
	}
	for _, v := range types {
		if v == m {
			return v
		}
	}
	return MirrorTypDefault
}

func MirrorTypesFromSlice(types []string) []MirrorType {
	res := make([]MirrorType, 0)
	for _, v := range types {
		res = append(res, MirrorType(v))
	}
	return res
}

const (
	MirrorTypDefault               = MirrorTypNone
	MirrorTypCN         MirrorType = "cn"
	MirrorTypNone       MirrorType = "none"
	MirrorTypCustomized MirrorType = "customized"
)

type envConfig struct {
	coreConfig *MainConfig
	ValueScope
	envDirname       string
	envName          string
	ID               string    `toml:"-"`
	CacheDir         string    `toml:"cache_dir"`
	DirName          string    `toml:"dir_name"`
	BootstrapExec    []Command `toml:"bootstrap_exec"`
	IsCommonPlatform bool      `toml:"is_common_platform"`
	ConfigExec       []Command `toml:"config_exec"`
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

func (c *envConfig) ExecuteConfig() error {
	c.PrepareScope()
	for _, command := range c.ConfigExec {
		if !command.CanRunning() {
			continue
		}
		command.ValueScope = c.ValueScope
		err := c.executeCommand(command)
		if err != nil {
			return err
		}
	}
	return nil
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
	if c.Top != nil {
		command.Top = c.Top
	}
	if c.Env == nil {
		c.Env = make(map[string]string)
	}
	c.Env["PATH"] = c.EnvPath.String()
	cmd := command.
		Prepare(c.Vars).
		SetEnv(c.Env)
	cmd.withAttr = false

	return cmd.Execute()
}

func (c envConfig) Copy(src ...envConfig) envConfig {
	var result envConfig
	if len(src) > 0 {
		tmp := src[0]
		result.coreConfig = tmp.coreConfig
		result.ValueScope = tmp.ValueScope
		result.envDirname = tmp.envDirname
		result.envName = tmp.envName
		result.ID = tmp.ID
		result.CacheDir = tmp.CacheDir
		result.DirName = tmp.DirName
		result.IsCommonPlatform = tmp.IsCommonPlatform
		result.BootstrapExec = tmp.BootstrapExec
		result.ConfigExec = tmp.ConfigExec
	} else {
		result.coreConfig = c.coreConfig
		result.ValueScope = c.ValueScope
		result.envDirname = c.envDirname
		result.envName = c.envName
		result.ID = c.ID
		result.CacheDir = c.CacheDir
		result.DirName = c.DirName
		result.IsCommonPlatform = c.IsCommonPlatform
		result.BootstrapExec = c.BootstrapExec
		result.ConfigExec = c.ConfigExec
	}
	return result
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
