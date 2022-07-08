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
		key := ""
		switch runtime.GOOS {
		case "windows":
			key = strings.ToUpper(pairs[0])
		default:
			key = pairs[0]
		}
		c.Environ[key] = pairs[1]
	}
}

const PathPlaceHolder = "+{PATH}"

func (c *Core) assignWithEnviron(key string) {
	if v, ok := c.Environ[key]; key != "" && ok {
		if key == "PATH" {
			c.Env[key] = PathPlaceHolder
			return
		}
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
		c.assignWithEnviron("APPDATA")
		c.assignWithEnviron("TEMP")
		c.assignWithEnviron("TMP")
		c.assignWithEnviron("WINDIR")
		c.assignWithEnviron("LOCALAPPDATA")
		c.assignWithEnviron("USERPROFILE")
		c.assignWithEnviron("PROGRAMW6432")
		c.assignWithEnviron("PATHEXT")
	case "linux", "freebsd", "openbsd", "macos", "ios", "android":
		c.assignWithEnviron("USER")
		c.assignWithEnviron("PWD")
		c.assignWithEnviron("LANG")
		c.assignWithEnviron("TMPDIR")
	default:
	}

}
