package internal

import (
	"os"
	"strings"
)

var argStr = map[string]string{
	"/#os/": string(os.PathSeparator),
}

func argReplace(s string) string {
	const prefix = "/#"
	const suffix = "/"
	isArgMap := strings.HasPrefix(s, prefix) && strings.HasSuffix(s, suffix)
	if !isArgMap {
		return s
	}
	fcs := strings.SplitN(s, "=", 2)
	if len(fcs) > 1 {
		fcs[0] = fcs[0] + suffix
	}
	if v0, ok := argStr[fcs[0]]; ok {
		fcs[0] = v0
	}
	if len(fcs) > 1 && len(fcs[1]) > 0 {
		if fcs[1] == "jvm/" && fcs[0] == "\\" {
			return "\\\\"
		}
	}
	return fcs[0]
}

const fnQuotting = "()"

var fnMaps = map[string]func(args ...string) string{
	"rp": func(args ...string) string {
		before := args[0]
		tmpArg := strings.SplitN(args[1], "=", 2)
		tmpArg[0] = argReplace(tmpArg[0])
		tmpArg[1] = argReplace(tmpArg[1])
		return strings.ReplaceAll(before, tmpArg[0], tmpArg[1])
	},
	"ospth": func(args ...string) string {
		src := args[0]
		return strings.ReplaceAll(src, "/", argReplace("/#os/"))
	},
}

const funcSeprator = ":"

func (v ValueScope) funcMapping(key string, src map[string]string) string {
	var result string
	keyFn := strings.Split(key, funcSeprator)
	for k, v0 := range src {
		if !strings.HasPrefix(k, keyFn[0]) {
			continue
		}
		innerKeyFn := strings.Split(k, funcSeprator)
		if len(innerKeyFn) < 2 {
			break
		}
		result = v.resolveFn(innerKeyFn, v0)
		if len(result) > 0 {
			break
		}
	}
	return result
}

func (v ValueScope) resolveFn(keyFn []string, result string) string {
	for _, v := range keyFn[1:] {
		idxLeft := strings.IndexRune(v, rune(fnQuotting[0]))
		idxRight := strings.IndexRune(v, rune(fnQuotting[1]))
		hasQuote := idxLeft > 0 && idxRight > 0
		if hasQuote {
			fnName := v[:idxLeft]
			if fn, ok := fnMaps[fnName]; ok {
				result = fn(result, v[idxLeft+1:idxRight])
			}
		}
	}
	return result
}
