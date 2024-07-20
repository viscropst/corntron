package internal

import "cryphtron/internal/utils"

const selectorPrefix = "+"

func platMapping[v any](key string, src map[string]v) v {
	var result v
	if v0, ok := src[key]; ok {
		result = v0
	}

	if v0, ok := src[key+utils.OsID(selectorPrefix)]; ok {
		result = v0
	}

	if v0, ok := src[key+utils.ArchID(selectorPrefix)]; ok {
		result = v0
	}

	if v0, ok := src[key+utils.PlatID(selectorPrefix)]; ok {
		result = v0
	}
	return result
}
