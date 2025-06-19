package internal

import (
	"os"
	"runtime"
	"strings"
)

type Environ map[string]string

var environ Environ

func GetEnvironMap() map[string]string {
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
	pathValue := GetEnvironMap()["PATH"]
	return AppendToPathList(pathValue, value)
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

func NormalizePathList(src string) string {
	_, result := NormalizePathListWithMap(src)
	return result
}

func NormalizePathListWithMap(src string) (map[string]int, string) {
	result := ""
	if len(src) == 0 {
		return map[string]int{src: 1}, src
	}
	pthList := strings.Split(src, OSPathListSeparator)
	if len(pthList) == 0 && strings.Contains(src, PathListSeparator) {
		pthList = strings.Split(src, PathListSeparator)
	}
	if len(pthList) < 2 {
		return map[string]int{src: 1}, src
	}
	pthCountMap := make(map[string]int)
	for _, v := range pthList {
		if len(v) == 0 {
			continue
		}
		tmp := NormalizePath(v)
		if osNoPrefix() == "windows" {
			tmp = strings.ToUpper(tmp)
		}
		if _, ok := pthCountMap[tmp]; !ok {
			pthCountMap[tmp] = 1
			result = result + OSPathListSeparator + NormalizePath(v)
		} else {
			continue
		}
	}
	return pthCountMap, strings.TrimPrefix(result, OSPathListSeparator)
}

func AppendToPathList(src, value string) string {
	if len(value) == 0 {
		return src
	}
	if src == value {
		return src
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
	srcListMap, result := NormalizePathListWithMap(src)
	valueList := strings.Split(NormalizePathList(value), OSPathListSeparator)
	if strings.Contains(value, PathListSeparator) {
		for _, v := range valueList {
			if len(v) == 0 {
				continue
			}
			tmp := NormalizePath(v)
			if osNoPrefix() == "windows" {
				tmp = strings.ToUpper(tmp)
			}
			if _, ok := srcListMap[tmp]; ok {
				continue
			} else {
				result = result + OSPathListSeparator + v
			}
		}
	} else {
		result = result + OSPathListSeparator + NormalizePath(value)
	}
	return strings.TrimPrefix(result, OSPathListSeparator)
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
