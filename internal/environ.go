package internal

import (
	"os"
	"strings"
)

type Environ map[string]string

var environ Environ

var environPath string

func GetEnvironMap() map[string]string {
	if environ == nil {
		environ = makeEnvironMap()
	}
	return environ
}

func GetEnvironPath() string {
	if len(environPath) == 0 {
		_ = makeEnvironMap()
	}
	return environPath
}

func makeEnvironMap() map[string]string {
	tmp := make(map[string]string)
	for _, s := range os.Environ() {
		pairs := strings.SplitN(s, "=", 2)
		if pairs[1] == "" {
			continue
		}
		key := ""
		switch osNoPrefix() {
		case "windows":
			key = strings.ToUpper(pairs[0])
		default:
			key = pairs[0]
		}
		if key == "PATH" {
			environPath = pairs[1]
			continue
		}
		tmp[key] = pairs[1]
	}
	return tmp
}

func AssertWithEnviron(args ...string) bool {
	if len(args) == 0 {
		return false
	}
	key := args[0]
	v, ok := GetEnvironMap()[key]
	if len(args) == 1 {
		return ok
	} else {
		value := args[1]
		return ok && v == value
	}
}

func FillEnviron(profileDir ...string) Environ {
	var result Environ = make(map[string]string)
	result.PrepareEnvsByEnviron(profileDir...)
	return result
}

func (c Environ) assignWithEnviron(key string) {
	if v, ok := GetEnvironMap()[key]; key != "" && ok {
		if key == "PATH" {
			return
		}
		c[key] = v
	}
}

func (c Environ) EnvStrList() []string {
	result := make([]string, 0)
	for k, v0 := range c {
		result = append(result, k+"="+v0)
	}
	return result
}
