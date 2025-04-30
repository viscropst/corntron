package utils

import (
	"os"
	"runtime"
	"strings"
)

var environ map[string]string

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
	pathValue := environ["PATH"]
	// pathValue = "+{PATH}"
	return AppendToPathList(pathValue, value)
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
