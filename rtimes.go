package cryphtron

import (
	"cryphtron/core"
	"cryphtron/internal"
	"cryphtron/internal/utils"
	"cryphtron/internal/utils/log"
	"fmt"
	"os"
	"path/filepath"
)

func (c Core) ComposeRtEnv() *internal.ValueScope {
	c.Prepare()
	pthValue := ""
	for _, config := range c.RuntimesEnv {
		config.RePrepareScope()
		mirrorEnv := config.UnwrapMirrorsEnv(c.Config.UnwrapMirrorType())
		for k, v := range mirrorEnv {
			mirrorEnv[k] = config.Expand(v)
		}
		pthValue = utils.AppendToPathList(pthValue, config.Env["PATH"])
		c.AppendEnvs(config.Env).AppendEnvs(mirrorEnv)
	}
	c.Env["PATH"] = pthValue
	return c.ValueScope
}

// ProcessRtBootstrap will execute bootstrap for each runtime.
// If ifResume is true, it will skip if runtime mirror cannot execute.
func (c Core) ProcessRtMirror(ifResume bool) error {
	mirrorType := c.Config.UnwrapMirrorType()
	c.Prepare()
	for _, runtime := range c.RuntimesEnv {
		config := runtime.Copy()
		config.PrepareScope()
		config.AppendEnvs(c.Env)
		config.Vars["pth_environ"] = c.Environ["PATH"]
		err := config.ExecuteMirrors(mirrorType)
		if err != nil {
			err = fmt.Errorf(core.RtIdentifier+" mirror[%s]:%s", mirrorType, err.Error())
			return err
		}
	}
	return nil
}

// ProcessRtBootstrap will execute bootstrap for each runtime.
// If ifResume is true, it will skip if runtime config cannot execute.
func (c Core) ProcessRtConfig(ifResume bool) error {
	for _, runtime := range c.RuntimesEnv {
		config := runtime.Copy()
		config.PrepareScope()
		config.AppendEnvs(c.Env)
		config.Vars["pth_environ"] = c.Environ["PATH"]
		err := config.ExecuteConfig()
		if err != nil {
			err = fmt.Errorf("%s config: %s", config.ID, err.Error())
			return err
		}
	}
	return nil
}

// ProcessRtBootstrap will execute bootstrap for each runtime.
// If ifResume is true, it will skip if runtime bootstrap cannot execute.
func (c Core) ProcessRtBootstrap(ifResume bool) error {
	c.Prepare()
	currentRtDir := c.Config.RtWithRunningDir()
	if ifFolderNotExists(currentRtDir) {
		_ = os.MkdirAll(currentRtDir, os.ModeDir|0o666)
	}

	for _, runtime := range c.RuntimesEnv {
		config := runtime.Copy()
		bootstrapDir := filepath.Join(currentRtDir,
			config.DirName)
		if !ifFolderNotExists(bootstrapDir) {
			continue
		}
		_ = os.Mkdir(bootstrapDir, os.ModeDir|0o666)
		config.PrepareScope()
		config.AppendEnvs(c.Env)
		config.Vars["pth_environ"] = c.Environ["PATH"]
		err := config.ExecuteBootstrap()
		if err != nil {
			_ = os.RemoveAll(bootstrapDir)
			err = fmt.Errorf(core.RtIdentifier+" bootstrsp[%s]:%s",
				config.DirName, err.Error())
			core.LogCLI(log.ErrorLevel).Println(err)
			if ifResume {
				return err
			}
		}
	}
	return nil
}
