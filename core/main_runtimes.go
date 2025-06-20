package core

import (
	"corntron/internal"
	"corntron/internal/log"
	"fmt"
	"path/filepath"
)

func (c Core) ComposeRtEnv() *ValueScope {
	c.Prepare()
	pthValue := PathListBuilder()
	for _, config := range c.RuntimesEnv {
		config.RePrepareScope()
		mirrorEnv := config.UnwrapMirrorsEnv(c.Config.UnwrapMirrorType())
		for k, v := range mirrorEnv {
			mirrorEnv[k] = config.Expand(v)
		}
		pthValue = pthValue.Append(config.Env["PATH"])
		c.AppendEnvs(config.Env).AppendEnvs(mirrorEnv)
	}
	c.EnvPath = pthValue
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
		config.Vars["pth_environ"] = c.EnvironPath.String()
		err := config.ExecuteMirrors(mirrorType)
		if err != nil {
			err = fmt.Errorf(RtIdentifier+" mirror[%s]:%s", mirrorType, err.Error())
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
		config.Vars["pth_environ"] = c.EnvironPath.String()
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
	if internal.IfFolderNotExists(currentRtDir) {
		_ = internal.Mkdir(currentRtDir)
	}

	for _, runtime := range c.RuntimesEnv {
		config := runtime.Copy()
		bootstrapDir := filepath.Join(currentRtDir,
			config.DirName)
		if !internal.IfFolderNotExists(bootstrapDir) {
			continue
		}
		_ = internal.Mkdir(bootstrapDir)
		config.PrepareScope()
		config.AppendEnvs(c.Env)
		config.Vars["pth_environ"] = c.EnvironPath.String()
		err := config.ExecuteBootstrap()
		if err != nil {
			_ = internal.Remove(bootstrapDir)
			err = fmt.Errorf(RtIdentifier+" bootstrsp[%s]:%s",
				config.DirName, err.Error())
			LogCLI(log.ErrorLevel).Println(err)
			if ifResume {
				return err
			}
		}
	}
	return nil
}
