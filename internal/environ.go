package internal

import (
	"os"
	"runtime"
	"strings"
)

var environ map[string]string

func hasEnvironSelector(str string) bool {
	return strings.HasSuffix(str, selectorPrefix+"environ")
}

func environMapping(key string) string {
	var result string
	canMapping := hasEnvironSelector(key)
	if !canMapping {
		return result
	}
	k := strings.TrimSuffix(key, selectorPrefix+"environ")
	if v0, ok := environ[k]; ok {
		result = v0
	}
	return result
}

func (c *Core) fillEnviron() {
	if environ == nil {
		environ = make(map[string]string)
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
		environ[key] = pairs[1]
	}
	c.Environ = environ
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
