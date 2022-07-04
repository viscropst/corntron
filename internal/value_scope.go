package internal

import (
	"os"
	"strings"
)

type ValueScope struct {
	scopeReady bool
	Top        *ValueScope       `toml:"-"`
	Vars       map[string]string `toml:"vars"`
	Env        map[string]string `toml:"envs"`
}

const valueRefFormat = "#{%s}"

func (v ValueScope) mappingScope(key string) string {
	var result string

	keyFn := strings.Split(key, ":")
	if v0, ok := v.Vars[keyFn[0]]; ok {
		result = v0
	}

	if v0, ok := v.Env[keyFn[0]]; ok {
		result = v0
	}

	if v.Top != nil && result == "" {
		v.Top.PrepareScope()
		result = v.Top.mappingScope(key)
	}

	if len(keyFn) > 1 && !(result == "" || result == key) {
		result = v.resolveFn(keyFn, result)
	}

	if result == "" {
		result = key
	}

	return result
}

func (v ValueScope) expandValue(str string) string {
	idx := 0
	var buf []byte

	for i := 0; i < len(str); i++ {
		if str[i] != valueRefFormat[0] && (i+1) >= len(str) {
			continue
		}
		name := ""
		offset := 0

		if buf == nil {
			buf = make([]byte, 0, 2*len(str))
		}
		buf = append(buf, str[idx:i]...)

		if str[i] == valueRefFormat[0] && str[i+1] == valueRefFormat[1] {
			inner := str[i+1:]
			for j := 1; j < len(inner); j++ {
				if inner[j] == valueRefFormat[4] && j > 1 {
					name = inner[1:j]
					offset = j + 1
				}
				if inner[j] == valueRefFormat[4] && j == 1 {
					offset = 2
				}
				if inner[j] == valueRefFormat[4] {
					break
				}
			}
		}

		if name == "" && offset > 0 {
		} else if name == "" {
			buf = append(buf, str[i])
		} else {
			scopeValue := v.mappingScope(name)
			if scopeValue == name {
				scopeValue = str[i : i+offset+1]
			}
			buf = append(buf, scopeValue...)
		}
		i = i + offset
		idx = i + 1
	}

	if buf == nil {
		return str
	}
	return string(buf) + str[idx:]
}

func (v *ValueScope) PrepareScope() {
	if v.scopeReady {
		return
	}

	if v.Top != nil && !v.Top.scopeReady {
		v.Top.PrepareScope()
	}

	for k, v0 := range v.Vars {
		v.Vars[k] = v.expandValue(v0)
	}

	for k, v0 := range v.Env {
		v.Env[k] = v.expandValue(v0)
	}

	v.scopeReady = true
}

func (v *ValueScope) modifyMap(from, to map[string]string,
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

func (v *ValueScope) AppendEnv(environ map[string]string) *ValueScope {
	v.PrepareScope()
	if len(environ) == 0 {
		return v
	}
	v.Env = v.modifyMap(environ, v.Env, func(k, a, b string) string {
		if a == b {
			return a
		}
		if k == "PATH" {
			a = strings.ReplaceAll(a, ";", string(os.PathListSeparator))
			b = strings.ReplaceAll(b, ";", string(os.PathListSeparator))
			a = strings.ReplaceAll(a, "/", string(os.PathSeparator))
			b = strings.ReplaceAll(b, "/", string(os.PathSeparator))
			var buf []byte
			buf = append(buf, b...)
			buf = append(buf, os.PathListSeparator)
			buf = append(buf, a...)
			return string(buf)
		}
		if a == "" {
			return v.expandValue(b)
		}
		return a
	})
	v.scopeReady = false
	return v
}

func (v *ValueScope) AppendVars(varToAdd map[string]string) *ValueScope {
	v.PrepareScope()
	v.Vars = v.modifyMap(varToAdd, v.Vars)
	v.scopeReady = false
	return v
}

func (v *ValueScope) EnvStrList() []string {
	result := make([]string, 0)
	for k, v0 := range v.Env {
		result = append(result, k+"="+v0)
	}
	return result
}

func (v *ValueScope) Expand(str string) string {
	v.PrepareScope()
	if len(str) > 0 {
		str = v.expandValue(str)
	}
	return str
}
