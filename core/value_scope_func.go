package core

import (
	"corntron/internal"
	"corntron/internal/log"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

var argStr = map[string]string{
	"/#os/": internal.OSPathSeparator,
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
		src = strings.ReplaceAll(src, internal.PathSeparator, argReplace("/#os/"))
		src = strings.ReplaceAll(src, internal.PathListSeparator, internal.OSPathListSeparator)
		return strings.Trim(src, internal.OSPathListSeparator)
	},
	"split-elem": func(args ...string) string {
		origin := args[0]
		funcArgs := strings.Split(args[1], ",")
		if len(funcArgs) == 0 && funcArgs[0] == "" {
			internal.LogCLI(log.PanicLevel).Println("empty string while doing split-elem")
		}
		if len(funcArgs) < 2 {
			internal.LogCLI(log.PanicLevel).Println("not enough args while doing split-elem")
		}
		splitStr := funcArgs[0]
		elemNum, err := strconv.Atoi(funcArgs[1])
		if err != nil {
			internal.LogCLI(log.PanicLevel).Println("error while doing split-elem:", err)
			return origin
		}
		splitList := strings.Split(origin, splitStr)
		if len(splitList) < elemNum {
			internal.LogCLI(log.PanicLevel).Println("out of range while doing split-elem")
			return origin
		}
		return splitList[elemNum-1]
	},
	"webreq": func(args ...string) string {
		origin := args[0]
		funcArgs := strings.Split(args[1], ",")
		if len(funcArgs) == 0 && funcArgs[0] == "" {
			internal.LogCLI(log.PanicLevel).Println("empty url while doing web request")
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
		result, err := internal.HttpRequestString(url.String(), funcArgs[2:]...)
		if err != nil {
			internal.LogCLI(log.ErrorLevel).Println("error while doing web request:", err)
			return origin
		}
		return strings.TrimSpace(result)
	},
	"gh-rel-ver": func(args ...string) string {
		origin := args[0]
		funcArgs := strings.Split(args[1], ",")
		if len(funcArgs) == 0 && funcArgs[0] == "" {
			internal.LogCLI(log.PanicLevel).Println("empty owner and project while doing gh-latest-rel")
		}
		project := strings.Split(funcArgs[0], "/")
		if len(project) < 2 {
			internal.LogCLI(log.PanicLevel).Println("unknown fromat of owner and project while doing gh-latest-rel")
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
		result, err := internal.HttpRequestBytes(url.String(), "GET")
		if err != nil {
			internal.LogCLI(log.ErrorLevel).Println("error while doing gh-latest-rel:", err)
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
	"gl-rel-ver": func(args ...string) string {
		origin := args[0]
		funcArgs := strings.Split(args[1], ",")
		if len(funcArgs) == 0 && funcArgs[0] == "" {
			internal.LogCLI(log.PanicLevel).Println("empty owner and project while doing gl-latest-rel")
		}
		project := strings.Split(funcArgs[0], "/")
		if len(project) < 2 {
			internal.LogCLI(log.PanicLevel).Println("unknown fromat of owner and project while doing gl-latest-rel")
		}
		apiPath := "/projects/" + project[0] + "%2F" + project[1] + "/releases"
		tagName := "latest"
		if len(funcArgs) > 1 {
			tagName = funcArgs[1]
		}
		if tagName == "latest" {
			apiPath = apiPath + "/permalink/latest"
		} else {
			apiPath = apiPath + "/tags/" + tagName
		}
		domain := "gitlab.com"
		if len(funcArgs) > 2 {
			domain = funcArgs[2]
		}
		apiUrl := domain + "/api/v4" + apiPath
		url, err := url.Parse("https://" + apiUrl)
		if err != nil {
			return origin
		}
		result, err := internal.HttpRequestBytes(url.String(), "GET")
		if err != nil {
			internal.LogCLI(log.ErrorLevel).Println("error while doing gl-latest-rel:", err)
			return origin
		}
		var glRelease struct {
			TagName string `json:"tag_name"`
		}
		err = json.Unmarshal(result, &glRelease)
		if err != nil {
			return origin
		}
		return strings.TrimSpace(glRelease.TagName)
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
