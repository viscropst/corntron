package core

import (
	"cryphtron/internal"
	"cryphtron/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

type splitString struct {
	SourceStr string    `toml:"src"`
	SplitStr  string    `toml:"split-str"`
	SplitNum  int8      `toml:"split-num"`
	Replaces  [2]string `toml:"replaces"`
}

func (sp *splitString) ToArray() []string {
	if len(sp.Replaces) == 2 {
		sp.SourceStr = strings.ReplaceAll(
			sp.SourceStr, sp.Replaces[0], sp.Replaces[1])
		if sp.Replaces[0] == "\n" && utils.OS() == "windows" {
			sp.SourceStr = strings.ReplaceAll(sp.SourceStr, "\r", "")
		}
	}
	if sp.SplitNum != 0 {
		return strings.SplitN(sp.SourceStr, sp.SplitStr, int(sp.SplitNum))
	}
	return strings.Split(sp.SourceStr, sp.SplitStr)
}

type Command struct {
	cmd     exec.Cmd
	Exec    string      `toml:"exec"`
	PlatStr string      `toml:"platform"`
	Args    []string    `toml:"args"`
	ArgStr  splitString `toml:"arg-str"`
	WorkDir string      `toml:"work-dir"`
	internal.ValueScope
	WithEnviron bool `toml:"with-environ"`
	WithNoWait  bool `toml:"with-no-wait"`
}

func (c *Command) CanRunning() bool {
	var canRunning = len(c.PlatStr) == 0
	canRunning = canRunning || c.PlatStr == utils.Arch()
	canRunning = canRunning || c.PlatStr == utils.OS()
	canRunning = canRunning || c.PlatStr == utils.Platform()
	return canRunning
}

func (c *Command) SetEnv(environ map[string]string) *Command {
	if environ == nil {
		return c
	}
	tmpEnv := c.Env
	c.Env = environ
	for k, v := range tmpEnv {
		_, ok := c.Env[k]
		if ok && len(v) > 0 {
			c.Env[k] = v
		} else if !ok && len(v) > 0 {
			c.Env[k] = v
		}
	}
	return c
}

func (c *Command) appendVars(src map[string]string) {
	filter := func(k, a, b string) string {
		if a == b {
			return c.Expand(a)
		}
		if a == "" {
			return c.Expand(b)
		} else {
			return c.Expand(a)
		}
	}
	c.Vars = utils.ModifyMap(src, c.Vars, filter)
}

func (c *Command) Prepare(vars ...map[string]string) *Command {
	c.cmd = exec.Cmd{
		Stderr: os.Stderr,
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
	}
	if len(vars) > 0 {
		c.appendVars(vars[0])
	}
	for idx := range c.Args {
		c.Args[idx] = c.Expand(c.Args[idx])
	}
	c.Exec = c.Expand(c.Exec)
	if len(c.WorkDir) > 0 {
		c.WorkDir = c.Expand(c.WorkDir)
		c.cmd.Dir = filepath.FromSlash(c.WorkDir)
	}
	if c.ArgStr.SourceStr != "" {
		c.ArgStr.SourceStr = c.Expand(c.ArgStr.SourceStr)
		c.Args = append(c.Args, c.ArgStr.ToArray()...)
	}
	return c
}

func (c *Command) prepareCmd() (*exec.Cmd, error) {
	cmd := exec.Cmd{
		Stdin:  c.cmd.Stdin,
		Stdout: c.cmd.Stdout,
		Stderr: c.cmd.Stderr,
		Dir:    c.cmd.Dir,
	}
	cmd.Path = c.Exec
	if filepath.Base(c.Exec) == c.Exec {
		var err0 error
		_ = os.Setenv("PATH", c.Env["PATH"])
		cmd.Path, err0 = exec.LookPath(c.Exec)
		_ = os.Unsetenv("PATH")
		if err0 != nil {
			return nil, err0
		}
	}

	cmd.Args = append(cmd.Args, cmd.Path)
	cmd.Args = append(cmd.Args, c.Args...)

	cmd.Env = c.EnvStrList()
	if c.WithEnviron {
		cmd.Env = append(cmd.Env, os.Environ()...)
	}
	return &cmd, nil
}

func (c *Command) Execute(vars ...map[string]string) error {
	c.Prepare(vars...)

	if v, ok := internal.Commands[c.Exec]; ok {
		return v(c.Args)
	}

	cmd, err := c.prepareCmd()
	if err != nil {
		return err
	}

	LogCLI(zerolog.InfoLevel).Println("executing command", cmd.String())
	err = cmd.Start()
	if err != nil {
		return err
	}

	if !c.WithNoWait {
		return cmd.Wait()
	}

	return cmd.Process.Release()
}
