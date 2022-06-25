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

func (v ValueScope) resolveFn(keyFn []string, result string) string {
	for _, v := range keyFn[1:] {
		idxLeft := strings.IndexRune(v, '(')
		idxRight := strings.IndexRune(v, ')')
		hasQuote := idxLeft > 0 && idxRight > 0
		if hasQuote {
			fnName := v[:idxLeft]
			switch fnName {
			case "rp":
				args := strings.SplitN(v[idxLeft+1:idxRight], "=", 2)
				args[0] = argReplace(args[0])
				args[1] = argReplace(args[1])
				result = strings.ReplaceAll(result, args[0], args[1])
			}
		}
	}
	return result
}
