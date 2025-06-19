package internal

import (
	"os"
	"runtime"
	"strings"
)

type Environ map[string]string

var environ Environ

func MakeEnvironMap() map[string]string {
	if environ == nil {
		environ = makeEnvironMap()
	}
	return environ
}
func makeEnvironMap() map[string]string {
	tmp := make(map[string]string)
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
		tmp[key] = pairs[1]
	}
	return tmp
}
func AppendToPath(value string) string {
	if AssertWithEnviron("PATH", value) {
		return value
	}
	pathValue := environ["PATH"]
	return AppendToPathList(pathValue, value)
}

func AssertWithEnviron(args ...string) bool {
	if len(args) == 0 {
		return false
	}
	key := args[0]
	v, ok := environ[key]
	if len(args) == 1 {
		return ok
	} else {
		value := args[1]
		return ok && v == value
	}
}

func AppendToPathList(src, value string) string {
	if len(value) == 0 {
		return src
	}
	if src == value {
		return src
	}
	if strings.Contains(value, PathListSeparator) {
		tmp := ""
		for _, v := range strings.Split(value, PathListSeparator) {
			tmp = tmp + OSPathListSeparator + NormalizePath(v)
		}
		if len(tmp) > 1 {
			value = tmp[1:]
		} else {
			value = tmp
		}
	} else {
		value = NormalizePath(value)
	}
	if len(src) == 0 {
		return value
	}
	if strings.Contains(src, value) {
		return src
	}
	if strings.Contains(value, src) {
		return value
	}
	pthList := strings.Split(src, string(OSPathListSeparator))
	for _, a := range pthList {
		if a == value {
			return src
		}
	}
	return src + OSPathListSeparator + value
}

func FillEnviron(profileDir ...string) Environ {
	if environ == nil {
		environ = make(map[string]string)
	}
	environ = MakeEnvironMap()
	var result Environ = make(map[string]string)
	result.PrepareEnvsByEnviron(profileDir...)
	return result
}

func (c Environ) assignWithEnviron(key string) {
	if v, ok := environ[key]; key != "" && ok {
		if key == "PATH" {
			return
		}
		c[key] = v
	}
}
