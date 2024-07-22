package internal

import "cryphtron/internal/utils"

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
