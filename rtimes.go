package cryphtron

import (
	"cryphtron/core"
	"cryphtron/internal"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

func (c Core) ComposeRtEnv() *internal.ValueScope {
	c.Prepare()
	for _, config := range c.RuntimesEnv {
		config.PrepareScope()
		mirrorEnv := config.UnwrapMirrorsEnv(c.Config.UnwrapMirrorType())
		for k, v := range mirrorEnv {
			mirrorEnv[k] = config.Expand(v)
		}
		pthArr := strings.SplitN(config.Env["PATH"], string(os.PathListSeparator), 2)
		if len(pthArr) > 1 && pthArr[0] == internal.PathPlaceHolder {
			idxAfter := strings.Index(internal.PathPlaceHolder, pthArr[0])
			config.Env["PATH"] = config.Env["PATH"][idxAfter+len(pthArr[0])+1:]
		}
		c.AppendEnv(config.Env).AppendEnv(mirrorEnv)
	}

	return c.ValueScope
}

// ProcessRtBootstrap will execute bootstrap for each runtime.
// If ifResume is true, it will skip if runtime mirror cannot execute.
func (c Core) ProcessRtMirror(ifResume bool) error {
	mirrorType := c.Config.UnwrapMirrorType()
	c.Prepare()
	for _, config := range c.RuntimesEnv {
		config.PrepareScope()
		config.AppendEnv(c.Env)
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
// If ifResume is true, it will skip if runtime bootstrap cannot execute.
func (c Core) ProcessRtBootstrap(ifResume bool) error {
	c.Prepare()
	currentRtDir := c.Config.RtWithRunningDir()
	if ifFolderNotExists(currentRtDir) {
		_ = os.MkdirAll(currentRtDir, os.ModeDir|0o666)
	}

	for _, config := range c.RuntimesEnv {
		bootstrapDir := filepath.Join(currentRtDir,
			config.DirName)
		if !ifFolderNotExists(bootstrapDir) {
			continue
		}
		_ = os.Mkdir(bootstrapDir, os.ModeDir|0o666)
		config.PrepareScope()
		config.AppendEnv(c.Env)
		config.Vars["pth_environ"] = c.Environ["PATH"]
		err := config.ExecuteBootstrap()
		if err != nil {
			_ = os.RemoveAll(bootstrapDir)
			err = fmt.Errorf(core.RtIdentifier+" bootstrsp[%s]:%s",
				config.DirName, err.Error())
			core.LogCLI(zerolog.ErrorLevel).Println(err)
			if ifResume {
				return err
			}
		}
	}
	return nil
}
