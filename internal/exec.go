package internal

import (
	"corntron/internal/log"
	"os"
	"os/exec"
	"path/filepath"
)

type Exec struct {
	cmd          *exec.Cmd
	Exec         string
	Args         []string
	WorkDir      string
	Env          Environ
	WithWaiting  bool
	IsBackground bool
	WithEnviron  bool
}

func (c *Exec) prepareCmd() (*exec.Cmd, error) {
	cmd := exec.Cmd{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    c.WorkDir,
	}
	cmd.Path = NormalizePath(c.Exec)
	if pth, _ := filepath.Split(c.Exec); pth == "" {
		var err0 error
		cmd.Path, err0 = GetExecPath(cmd.Path, c.Env["PATH"])
		if err0 != nil {
			return nil, err0
		}
	}

	cmd.Args = append(cmd.Args, cmd.Path)
	cmd.Args = append(cmd.Args, c.Args...)
	cmd.Env = c.Env.EnvStrList()
	return &cmd, nil
}

func (c *Exec) Execute() error {
	var err error
	c.cmd, err = c.prepareCmd()
	if err != nil {
		return err
	}

	LogCLI(log.InfoLevel).Println("executing command", c.cmd.String())

	attr := GetNewProcGroupAttr(c.IsBackground, !c.WithWaiting)
	if attr != nil {
		c.cmd.SysProcAttr = attr
	}

	if c.WithWaiting {
		return c.cmd.Run()
	}

	err = c.cmd.Start()
	if err != nil {
		return err
	}
	return c.cmd.Process.Release()
}
