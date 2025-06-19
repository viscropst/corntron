package core

import (
	"corntron/internal"
	"path/filepath"
	"strings"
)

type splitString struct {
	SourceStr string    `toml:"src"`
	SplitStr  string    `toml:"split_str"`
	SplitNum  int8      `toml:"split_num"`
	Replaces  [2]string `toml:"replaces"`
}

func (sp *splitString) ToArray() []string {
	if len(sp.Replaces) == 2 {
		if sp.Replaces[0] == "\n" {
			sp.SourceStr = strings.ReplaceAll(sp.SourceStr, "\r", sp.Replaces[1])
		}
		sp.SourceStr = strings.ReplaceAll(
			sp.SourceStr, sp.Replaces[0], sp.Replaces[1])
	}
	if sp.SplitNum != 0 {
		return strings.SplitN(sp.SourceStr, sp.SplitStr, int(sp.SplitNum))
	}
	return strings.Split(sp.SourceStr, sp.SplitStr)
}

type Command struct {
	workDir       string
	withWaiting   bool
	withAttr      bool
	GlobalWaiting bool        `toml:"-"`
	Exec          string      `toml:"exec"`
	PlatStr       string      `toml:"platform"`
	Args          []string    `toml:"args"`
	ArgStr        splitString `toml:"arg_str"`
	WorkDir       string      `toml:"work_dir"`
	ValueScope
	WithEnviron   bool `toml:"with_environ"`
	WithNoWaiting bool `toml:"with_no_waiting"`
	IsBackground  bool `toml:"is_background"`
}

func (c *Command) CanRunning() bool {
	var canRunning = len(c.PlatStr) == 0
	canRunning = canRunning || c.PlatStr == internal.Arch()
	canRunning = canRunning || c.PlatStr == internal.OS()
	canRunning = canRunning || c.PlatStr == internal.Platform()
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

func (c *Command) Prepare(vars ...map[string]string) *Command {
	if !c.WithNoWaiting {
		c.withWaiting = c.GlobalWaiting && internal.IsInTerminal()
	} else {
		c.withWaiting = !c.WithNoWaiting
	}
	if len(vars) > 0 {
		c.AppendVarsByNew(vars[0])
	}
	for idx := range c.Args {
		c.Args[idx] = c.Expand(c.Args[idx])
	}
	c.Exec = c.Expand(c.Exec)
	if len(c.WorkDir) > 0 {
		c.WorkDir = c.Expand(c.WorkDir)
		c.workDir = filepath.FromSlash(c.WorkDir)
	}
	if c.ArgStr.SourceStr != "" {
		c.ArgStr.SourceStr = c.Expand(c.ArgStr.SourceStr)
		c.Args = append(c.Args, c.ArgStr.ToArray()...)
	}
	return c
}

func (c *Command) Execute(vars ...map[string]string) error {
	c.Prepare(vars...)

	if v, ok := Commands[c.Exec]; ok {
		return v(c.Args)
	}

	command := internal.Exec{
		Exec:        c.Exec,
		Args:        c.Args,
		WorkDir:     c.workDir,
		Env:         c.Env,
		WithEnviron: c.WithEnviron,
		WithWaiting: true,
	}

	if c.WithEnviron {
		command.Env = appendMap(c.Env, internal.GetEnvironMap())
	}

	if c.withAttr {
		command.WithWaiting = c.withWaiting
		command.IsBackground = c.IsBackground
	}
	return command.Execute()
}
