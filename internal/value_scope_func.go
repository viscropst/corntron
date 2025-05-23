package internal

import (
	"corntron/internal/utils"
	"corntron/internal/utils/log"
	"encoding/json"
	"net/url"
	"strings"
)

var argStr = map[string]string{
	"/#os/": utils.OSPathSeparator,
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
		src = strings.ReplaceAll(src, utils.PathSeparator, argReplace("/#os/"))
		src = strings.ReplaceAll(src, utils.PathListSeparator, utils.OSPathListSeparator)
		return strings.Trim(src, utils.OSPathListSeparator)
	},
	"webreq": func(args ...string) string {
		origin := args[0]
		funcArgs := strings.Split(args[1], ",")
		if len(funcArgs) == 0 && funcArgs[0] == "" {
			utils.LogCLI(log.PanicLevel).Println("empty url while doing web request")
		}
		if len(funcArgs) < 2 {
			funcArgs = append(funcArgs, "http")
		}
		if len(funcArgs) < 3 {
			funcArgs = append(funcArgs, "GET")
		}
		url, err := url.Parse(funcArgs[1] + "://" + funcArgs[0])
		if err != nil {
			return origin
		}
		result, err := utils.HttpRequestString(url.String(), funcArgs[2:]...)
		if err != nil {
			utils.LogCLI(log.ErrorLevel).Println("error while doing web request:", err)
			return origin
		}
		return strings.TrimSpace(result)
	},
	"gh-rel-ver": func(args ...string) string {
		origin := args[0]
		funcArgs := strings.Split(args[1], ",")
		if len(funcArgs) == 0 && funcArgs[0] == "" {
			utils.LogCLI(log.PanicLevel).Println("empty owner and project while doing gh-latest-rel")
		}
		project := strings.Split(funcArgs[0], "/")
		if len(project) < 2 {
			utils.LogCLI(log.PanicLevel).Println("unknown fromat of owner and project while doing gh-latest-rel")
		}
		apiPath := "/" + project[0] + "/" + project[1] + "/releases"
		tagName := "latest"
		if len(funcArgs) > 1 {
			tagName = funcArgs[1]
		}
		if tagName == "latest" {
			apiPath = apiPath + "/latest"
		} else {
			apiPath = apiPath + "/tags/" + tagName
		}
		domain := "github.com"
		if len(funcArgs) > 2 {
			domain = funcArgs[2]
		}
		apiUrl := "api." + domain + "/repos" + apiPath
		url, err := url.Parse("https://" + apiUrl)
		if err != nil {
			return origin
		}
		result, err := utils.HttpRequestBytes(url.String(), "GET")
		if err != nil {
			utils.LogCLI(log.ErrorLevel).Println("error while doing gh-latest-rel:", err)
			return origin
		}
		var ghRelease struct {
			TagName string `json:"tag_name"`
		}
		err = json.Unmarshal(result, &ghRelease)
		if err != nil {
			return origin
		}
		return strings.TrimSpace(ghRelease.TagName)
	},
}

const funcSeprator = ":"

func (v ValueScope) funcMapping(key string, src map[string]string) string {
	keyFn := strings.Split(key, funcSeprator)
	resultKey, resultValue := platMappingWithKey(keyFn[0], key, src)
	innerKeyFn := strings.Split(resultKey, funcSeprator)
	if len(innerKeyFn) > 1 {
		resultValue = v.resolveFn(innerKeyFn, resultValue)
	}
	return resultValue
}

func (v ValueScope) resolveFn(keyFn []string, result string) string {
	for _, fv := range keyFn[1:] {
		idxLeft := strings.IndexRune(fv, rune(fnQuotting[0]))
		idxRight := strings.IndexRune(fv, rune(fnQuotting[1]))
		hasQuote := idxLeft > 0 && idxRight > 0
		if hasQuote {
			fnName := fv[:idxLeft]
			fnValue := v.expandValue(fv[idxLeft+1 : idxRight])
			if fn, ok := fnMaps[fnName]; ok {
				result = fn(result, fnValue)
			}
		}
	}
	return result
}
