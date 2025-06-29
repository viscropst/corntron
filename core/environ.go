package core

import (
	"corntron/internal"
	"strings"
)

var environ = internal.GetEnvironMap()

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
