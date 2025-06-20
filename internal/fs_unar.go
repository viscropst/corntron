package internal

import (
	"errors"
	"os"
)

type UnarLogicPath uint8

const (
	Default UnarLogicPath = iota
	UnZip
	UnTar
)

type UnarchiveFlag struct {
	LogicType  UnarLogicPath
	OutputPath string
	SourceFile string
	BaseDir    string
}

func Unarchive(srcFile *os.File, flags UnarchiveFlag, includes ...string) error {
	switch flags.LogicType {
	case UnTar:
		return UnTarFromFileWithBaseDir(
			srcFile, flags.OutputPath, flags.BaseDir, includes...)
	case UnZip:
		return UnZipFromFileWithBaseDir(
			srcFile, flags.OutputPath, flags.BaseDir, includes...)
	default:
		return errors.New("unknown command")
	}
}
