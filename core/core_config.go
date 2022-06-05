package core

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type MainConfig struct {
	CurrentDir string
	RuntimeDir string
	AppDir     string
	MirrorType string
	WithApp    bool
}

func (c MainConfig) FsWalk(walkFunc filepath.WalkFunc, DirNames ...string) error {
	rootDir := c.CurrentDir
	if len(DirNames) > 0 {
		tmp := path.Join(DirNames...)
		rootDir = path.Join(rootDir, tmp)
	}
	return filepath.Walk(rootDir, walkFunc)
}

func (c MainConfig) FsWalkDir(walkFunc fs.WalkDirFunc, DirNames ...string) error {
	rootDir := path.Base(c.CurrentDir)
	if len(DirNames) > 0 {
		tmp := path.Join(DirNames...)
		rootDir = path.Join(rootDir, tmp)
	}
	return filepath.WalkDir(rootDir, walkFunc)
}

func (c MainConfig) UnwrapMirrorType() MirrorType {
	return MirrorType(c.MirrorType).Convert()
}

const execDirWithoutLink = "${dp0}"

var defaultCoreConfig = MainConfig{
	CurrentDir: execDirWithoutLink,
	RuntimeDir: "runtimes",
	AppDir:     "apps",
}

func loadConfigRegular(config string, value interface{}, altBases ...string) error {
	errFmt := fs.PathError{
		Op:   "loadConfig",
		Path: config,
	}

	if len(config) == 0 {
		errFmt.Err = fmt.Errorf("could not load config by empty name")
		return &errFmt
	}

	basePath := ""
	if len(altBases) > 0 {
		basePath = altBases[0]
	} else {
		basePath, _ = os.Getwd()
	}
	if len(basePath) == 0 {
		return fmt.Errorf("could not load workdir")
	}

	tomlFilename := path.Join(basePath, config+".toml")
	stat, _ := os.Stat(tomlFilename)
	if stat == nil || !stat.Mode().IsRegular() {
		errFmt.Path = tomlFilename
		errFmt.Err = fmt.Errorf("could not stat that file")
		return &errFmt
	}

	tomlFile, _ := os.Open(tomlFilename)
	tomlDecoder := toml.NewDecoder(tomlFile)
	//tomlDecoder.DisallowUnknownFields()
	err := tomlDecoder.Decode(value)
	if err != nil {
		errFmt.Path = tomlFilename
		errFmt.Err = err
		return &errFmt
	}
	return nil
}

func LoadCoreConfig(altBases ...string) MainConfig {
	basePath, _ := os.Getwd()
	if len(basePath) == 0 && len(altBases) == 0 {
		panic("Could not load workdir,use altBase instead")
	}

	if len(altBases) > 0 {
		basePath = altBases[0]
	}

	var result = defaultCoreConfig
	err := loadConfigRegular("core", &result, basePath)
	if err != nil {
		fmt.Println("WARN:>", err.Error())
	}

	if result.CurrentDir == execDirWithoutLink {
		result.CurrentDir = basePath
	}

	if len(result.AppDir) == 0 {
		result.AppDir = defaultCoreConfig.AppDir
	}

	return result
}
