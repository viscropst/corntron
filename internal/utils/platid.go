package utils

import "runtime"

const platArchSeprator = "-"

var OS = osNoPrefix
var Arch = archNoPrefix
var Platform = platNoPrefix

func osNoPrefix() string {
	return OsId("")
}
func archNoPrefix() string {
	return ArchId("")
}
func platNoPrefix() string {
	return PlatId("")
}

func OsId(prefix string) string {
	goosSuffix := prefix + runtime.GOOS
	return goosSuffix
}

func ArchId(prefix string) string {
	goarchSuffix := prefix + runtime.GOARCH
	return goarchSuffix
}

func PlatId(prefix string) string {
	goplatSuffix := prefix +
		runtime.GOOS +
		platArchSeprator +
		runtime.GOARCH
	return goplatSuffix
}
