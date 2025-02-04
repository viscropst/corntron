package internal

import (
	"cryphtron/internal/utils"
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
	environ = utils.MakeEnvironMap()
	c.Environ = environ
}

const PathPlaceHolder = "+{PATH}"

func (c *Core) assignWithEnviron(key string) {
	if v, ok := c.Environ[key]; key != "" && ok {
		if key == "PATH" {
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
