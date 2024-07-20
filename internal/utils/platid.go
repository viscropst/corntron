package utils

import "runtime"

const platArchSeprator = "-"

var OS = osNoPrefix
var Arch = archNoPrefix
var Platform = platNoPrefix

func osNoPrefix() string {
	return OsID("")
}
func archNoPrefix() string {
	return ArchID("")
}
func platNoPrefix() string {
	return PlatID("")
}

func OsID(prefix string) string {
	goosSuffix := prefix + runtime.GOOS
	return goosSuffix
}

func ArchID(prefix string) string {
	goarchSuffix := prefix + runtime.GOARCH
	return goarchSuffix
}

func PlatID(prefix string) string {
	goplatSuffix := prefix +
		runtime.GOOS +
		platArchSeprator +
		runtime.GOARCH
	return goplatSuffix
}
