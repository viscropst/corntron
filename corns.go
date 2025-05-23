package corntron

import (
	"corntron/core"
	"corntron/internal"
	"errors"
)

func (c Core) ComposeCornEnv(corn *core.CornsEnvConfig) *internal.ValueScope {
	c.Prepare()
	if corn == nil {
		return c.ValueScope
	}
	for _, depends := range corn.DependCorns {
		config := c.CornsEnv[depends]
		config.AppendVars(corn.Vars)
		config.RePrepareScope()
		corn.AppendEnvs(config.Env)
		corn.AppendVars(config.Vars)
	}
	return &corn.ValueScope
}

func (c *Core) prepareCorn(name string) (*core.CornsEnvConfig, error) {
	corn, ok := c.CornsEnv[name]
	if !ok {
		return nil, errors.New("could not found the " +
			core.CornsIdentifier + " named " + name)
	}
	return c.prepareCornConfig(corn)
}

func (c *Core) execCorn(name string, isWaiting bool, args ...string) error {
	corn, ok := c.CornsEnv[name]
	if !ok {
		return errors.New("could not found the " +
			core.CornsIdentifier + " named " + name)
	}
	return c.execCornConfig(corn, isWaiting, args...)
}

func (c *Core) ExecCorn(name string, isWaiting bool, args ...string) error {
	return c.execCorn(name, isWaiting, args...)
}
