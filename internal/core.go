package internal

type Core struct {
	*ValueScope
	Environ    map[string]string
	ProfileDir string
}

func (c *Core) Prepare() {
	if c.Environ != nil {
		return
	}
	c.fillEnviron()
	c.prepareEnvsByEnviron()
}
