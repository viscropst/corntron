package cryphtron

import (
	"cryphtron/internal"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func (c Core) ProcessRtMirror() error {
	mirrorType := c.Config.UnwrapMirrorType()
	c.Prepare()
	for _, config := range c.RuntimesEnv {
		config.PrepareScope()
		config.AppendEnv(c.Env)
		config.Vars["pth_environ"] = c.Environ["PATH"]
		err := config.ExecuteMirrors(mirrorType)
		if err != nil {
			return fmt.Errorf("mirror[%s]:%s", mirrorType, err.Error())
		}
	}
	return nil
}

func (c Core) ProcessRtBootstrap() error {
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
			return fmt.Errorf("bootstrsp[%s]:%s",
				config.DirName, err.Error())
		}
	}
	return nil
}
