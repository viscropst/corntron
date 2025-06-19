package unarchive

import (
	"archive/zip"
	"os"
	"sort"
)

func ZipReader(src *os.File) (*zip.Reader, error) {
	stat, err := src.Stat()
	if err != nil {
		return nil, err
	}
	return zip.NewReader(src, stat.Size())
}

func FilterZipFiles(src []*zip.File, includes ...string) []*zip.File {
	sort.Slice(src, func(i, j int) bool {
		return src[i].Name < src[j].Name
	})
	if len(includes) == 0 || len(src) == 0 {
		return src
	}
	result := make([]*zip.File, 0, getMinLen(src, includes))
	for _, v := range src {
		if v == nil {
			continue
		}
		if IsInInclude(v, includes...) {
			result = append(result, v)
		}
	}
	return result
}

func getMinLen(src []*zip.File, includes []string) int {
	result := 0
	if len(src) < len(includes) {
		result = len(src)
	} else if len(src) > len(includes) {
		result = len(includes)
	} else {
		result = len(src)
	}
	return result
}
