package utils

import (
	"errors"
	"io/fs"
	"os"

	"github.com/BurntSushi/toml"
)

func LoadTomlFile(tomlFilename string, value interface{}) *fs.PathError {
	var errFmt = fs.PathError{
		Op:   "LoadTomlFile",
		Path: tomlFilename,
	}
	stat, _ := StatPath(tomlFilename)
	if stat == nil || !stat.Mode().IsRegular() {
		errFmt.Path = tomlFilename
		errFmt.Err = errors.New("could not stat that file")
		return &errFmt
	}

	tomlFile, _ := os.Open(tomlFilename)
	tomlDecoder := toml.NewDecoder(tomlFile)
	//tomlDecoder.DisallowUnknownFields()
	_, err := tomlDecoder.Decode(value)
	if err != nil {
		errFmt.Path = tomlFilename
		errFmt.Err = err
		return &errFmt
	}
	return nil
}
