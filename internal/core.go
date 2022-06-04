package internal

import (
	"os"
	"runtime"
	"strings"
)

type Core struct {
	*ValueScope
	Environ map[string]string
}

func (c *Core) fillEnviron() {
	if c.Environ == nil {
		c.Environ = make(map[string]string)
	}
	for _, s := range os.Environ() {
		pairs := strings.SplitN(s, "=", 2)
		if pairs[1] == "" {
			continue
		}
		c.Environ[pairs[0]] = pairs[1]
	}
}

func (c *Core) assignWithEnviron(key string) {
	if v, ok := c.Environ[key]; key != "" && ok {
		c.Env[key] = v
	}
}

func (c *Core) Prepare() {
	if c.Environ != nil {
		return
	}
	c.fillEnviron()
	c.assignWithEnviron("PATH")

	switch runtime.GOOS {
	case "windows":
		c.assignWithEnviron("USERNAME")
	case "linux", "freebsd", "openbsd", "macos", "ios", "android":
		c.assignWithEnviron("USER")
		c.assignWithEnviron("PWD")
		c.assignWithEnviron("LANG")
	default:
	}

}
