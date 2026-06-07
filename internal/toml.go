package internal

import (
	"io"
	"io/fs"
	"os"

	"github.com/BurntSushi/toml"
)

func LoadTomlFile(tomlFilename string, value any) *fs.PathError {
	var errFmt = fs.PathError{
		Op:   "LoadTomlFile",
		Path: tomlFilename,
	}
	stat, _ := StatPath(tomlFilename)
	if stat == nil || !stat.Mode().IsRegular() {
		errFmt.Path = tomlFilename
		errFmt.Err = Error("could not stat that file")
		return &errFmt
	}

	tomlFile, _ := os.Open(tomlFilename)
	defer tomlFile.Close()
	err := LoadTomlReader(tomlFile, value)
	if err != nil {
		errFmt.Path = tomlFilename
		errFmt.Err = err
		return &errFmt
	}
	return nil
}

func LoadTomlReader(reader io.Reader, value any) error {
	tomlDecoder := toml.NewDecoder(reader)
	//tomlDecoder.DisallowUnknownFields()
	_, err := tomlDecoder.Decode(value)
	if err != nil {
		return Error("failed to decode from reader: ", err.Error())
	}
	return nil
}
