package cryphtron

import (
	"cryphtron/internal"
	"path"
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
	}
	return MirrorTypDefault
}

const (
	MirrorTypDefault            = MirrorTypNone
	MirrorTypCN      MirrorType = "cn"
	MirrorTypNone    MirrorType = "none"
)

type envConfig struct {
	coreConfig *CoreConfig
	internal.ValueScope
	envDirname string
	ID         string
	CacheDir   string
}

func (c envConfig) setCore(coreConfig CoreConfig) envConfig {
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
		c.envDirname = altCacheDirname[0]
	}
	if c.CacheDir == "" {
		c.CacheDir = "_cache"
	}
}

type RtEnvConfig struct {
	envConfig
	MirrorEnv map[MirrorType]map[string]string
}

func (c *RtEnvConfig) UnwrapMirrorsEnv(key MirrorType) map[string]string {
	var result = make(map[string]string)
	for k, v := range c.MirrorEnv[key] {
		result[k] = v
	}
	return result
}

func LoadRtEnv(name string, coreConfig CoreConfig, altEnvDirname ...string) (RtEnvConfig, error) {
	c := RtEnvConfig{}
	c.setEnvDirname(altEnvDirname...)
	c.setCacheDirname()
	c.AppendVars(map[string]string{
		"rt_dir":   path.Join(coreConfig.CurrentDir, coreConfig.RuntimeDir),
		"rt_cache": path.Join(coreConfig.CurrentDir, coreConfig.RuntimeDir, c.CacheDir),
	})
	loadPath := path.Join(coreConfig.CurrentDir, coreConfig.RuntimeDir, c.envDirname)
	err := loadConfigRegular(name, &c, loadPath)
	if err != nil {
		c.setCore(coreConfig)
		return c, err
	}
	c.envConfig = c.setCore(coreConfig)
	return c, nil
}

type EditEnvConfig struct {
	envConfig
	EditorExec []Command
}

func (c EditEnvConfig) ExecuteAll() error {
	c.PrepareScope()
	for _, command := range c.EditorExec {
		err := command.Prepare().
			SetEnv(c.Env).Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadEditEnv(name string, coreConfig CoreConfig, altEnvDirname ...string) (EditEnvConfig, error) {
	result := EditEnvConfig{}
	result.setEnvDirname(altEnvDirname...)
	result.setCacheDirname()
	result.AppendVars(map[string]string{
		"ed_dir":   path.Join(coreConfig.CurrentDir, coreConfig.EditorDir),
		"ed_cache": path.Join(coreConfig.CurrentDir, coreConfig.EditorDir, result.CacheDir),
	})
	loadPath := path.Join(coreConfig.CurrentDir, coreConfig.EditorDir, result.envDirname)
	err := loadConfigRegular(name, &result, loadPath)
	if err != nil {
		result.setCore(coreConfig)
		return result, err
	}

	for idx := range result.EditorExec {
		result.EditorExec[idx].Top = &result.ValueScope
	}
	result.envConfig = result.setCore(coreConfig)
	return result, nil
}
