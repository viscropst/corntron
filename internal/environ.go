package internal

import (
	"os"
	"runtime"
	"strings"
)

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

func (c *Core) assertWithEnviron(args ...string) bool {
	if len(args) == 0 {
		return false
	}
	key := args[0]
	v, ok := c.Environ[key]
	if len(args) == 1 {
		return ok
	} else {
		value := args[1]
		return ok && v == value
	}
}
