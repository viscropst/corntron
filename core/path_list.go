package core

import (
	"corntron/internal"
	"strings"
)

type PathList map[string]int

func PathListBuilder(src ...string) PathList {
	result := make(PathList)
	if len(src) == 0 {
		return result
	}
	if len(src) < 2 {
		return result.Append(src[0])
	}
	for _, v := range src {
		result = result.Append(v)
	}
	return result
}

func (l PathList) ToUpper() PathList {
	result := PathListBuilder()
	for k, v := range l {
		tmp := internal.NormalizePath(k)
		tmp = strings.ToUpper(tmp)
		result[tmp] = v
	}
	return result
}

func (l PathList) Append(src string) PathList {
	if len(src) == 0 {
		return l
	}
	tmpSrc := internal.NormalizePath(src)
	pthList := strings.Split(tmpSrc, internal.OSPathListSeparator)
	if len(pthList) < 2 && strings.Contains(tmpSrc, internal.PathListSeparator) {
		pthList = strings.Split(tmpSrc, internal.PathListSeparator)
	}
	if len(pthList) < 2 {
		l[tmpSrc] = 1
		return l
	}
	dstMap := l
	if internal.OS() == "windows" {
		dstMap = dstMap.ToUpper()
	}
	for _, v := range pthList {
		if len(v) == 0 {
			continue
		}
		tmp := internal.NormalizePath(v)
		if internal.OS() == "windows" {
			tmp = strings.ToUpper(tmp)
		}
		if _, ok := dstMap[tmp]; !ok {
			l[v] = 1
		} else {
			continue
		}
	}
	return l
}

func (l PathList) AppendList(src PathList) PathList {
	if len(src) == 0 {
		return l
	}
	dstMap := l
	if internal.OS() == "windows" {
		dstMap = dstMap.ToUpper()
	}
	for k := range src {
		if len(k) == 0 {
			continue
		}
		tmp := internal.NormalizePath(k)
		if internal.OS() == "windows" {
			tmp = strings.ToUpper(tmp)
		}
		if _, ok := dstMap[tmp]; !ok {
			l[k] = 1
		} else {
			continue
		}
	}
	return l
}

func (l PathList) String() string {
	result := ""
	for k := range l {
		result = result + internal.OSPathListSeparator + k
	}
	return strings.TrimPrefix(result, internal.OSPathListSeparator)
}

var environPathList PathList

func EnvironPathList() PathList {
	if environPathList == nil {
		environPathList = PathListBuilder(internal.GetEnvironPath())
	}
	return environPathList
}
