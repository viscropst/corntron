package internal

import (
	"cryphtron/internal/utils"
	"strings"
)

const selectorPrefix = "+"

func platMapping[v any](key string, altKey string, src map[string]v) v {
	var result v
	if v0, ok := src[key]; ok {
		result = v0
	}

	var tmpKey = key + utils.OsID(selectorPrefix)
	if v0, ok := src[tmpKey]; ok && altKey != tmpKey {
		result = v0
	}
	tmpKey = key + utils.ArchID(selectorPrefix)
	if v0, ok := src[tmpKey]; ok && altKey != tmpKey {
		result = v0
	}
	tmpKey = key + utils.PlatID(selectorPrefix)
	if v0, ok := src[tmpKey]; ok && altKey != tmpKey {
		result = v0
	}
	return result
}

func platMappingWithKey[v any](key string, altKey string, src map[string]v) (string, v) {
	var result = key
	for k := range src {
		if !strings.HasPrefix(k, key) {
			continue
		}

		splitFunc := strings.Split(k, funcSeprator)
		splitSelector := strings.Split(splitFunc[0], selectorPrefix)
		if len(splitSelector) < 2 {
			result = k
			continue
		}
		if len(splitSelector) >= 2 && splitSelector[1] == utils.OS() {
			result = k
			break
		}
		if len(splitSelector) >= 2 && splitSelector[1] == utils.Arch() {
			result = k
			break
		}
		if len(splitSelector) >= 2 && splitSelector[1] == utils.Platform() {
			result = k
			break
		}
	}
	return result, platMapping(key, altKey, src)
}
