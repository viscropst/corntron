package internal

func AppendMap[value any](from, to map[string]value,
	proc ...func(k string, a, b value) (string, value)) map[string]value {
	if len(from) == 0 {
		return to
	}
	if to == nil {
		to = make(map[string]value)
	}
	result := to
	for key, val := range from {
		tmpKey, tmpVal := key, val
		tmpOrig, ok := to[key]

		if len(proc) > 0 && ok {
			tmpKey, tmpVal = proc[0](tmpKey, tmpOrig, tmpVal)
		} else if len(proc) > 0 && !ok {
			tmpKey, tmpVal = proc[0](tmpKey, tmpVal, tmpVal)
		} else {
			tmpKey, tmpVal = proc[0](tmpKey, tmpOrig, tmpOrig)
		}

		result[tmpKey] = tmpVal
	}
	return result
}

func ModifyMapByMap[v any](from, to map[string]v,
	beforeAdd ...func(k string, a, b v) v) map[string]v {
	if len(from) == 0 {
		return to
	}
	if to == nil {
		to = make(map[string]v)
	}
	for k, v1 := range from {
		v0, ok := to[k]

		if ok && len(beforeAdd) > 0 {
			to[k] = beforeAdd[0](k, v0, v1)
		} else if !ok {
			to[k] = v1
		} else {
			to[k] = v0
		}
	}
	return to
}

func ModifyMapByPair[v any](to map[string]v, key string, value v, beforeAdd ...func(k string, a, b v) v) map[string]v {
	if to == nil {
		to = make(map[string]v)
	}
	v0, ok := to[key]
	if ok && len(beforeAdd) > 0 {
		to[key] = beforeAdd[0](key, v0, value)
	} else if !ok {
		to[key] = value
	} else {
		to[key] = v0
	}
	return to
}
