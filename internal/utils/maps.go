package utils

func ModifyMap(from, to map[string]string,
	beforeAdd ...func(k, a, b string) string) map[string]string {
	if len(from) == 0 {
		return to
	}
	if to == nil {
		to = make(map[string]string)
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
