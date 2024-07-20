package core

import (
	"cryphtron/internal/utils"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
)

type MainConfig struct {
	CurrentDir           string            `toml:"base_dir,omitempty"`
	RunningDir           string            `toml:"-"`
	RuntimeDirName       string            `toml:"runtime_dirname,omitempty"`
	CornDirName          string            `toml:"corn_dirname,omitempty"`
	MirrorType           string            `toml:"mirror_type,omitempty"`
	WithCorn             bool              `toml:"with_corn"`
	RunningDirByPlatfrom map[string]string `toml:"running_dir,omitempty"`
	ProfileDir           string            `toml:"profile_dir,omitempty"`
}

func (c MainConfig) RuntimeDir() string {
	return filepath.Join(c.CurrentDir, c.RuntimeDirName)
}

func (c MainConfig) CornDir() string {
	return filepath.Join(c.CurrentDir, c.CornDirName)
}

func (c MainConfig) RtWithRunningDir() string {
	return filepath.Join(c.RunningDir, c.RuntimeDirName)
}

func (c MainConfig) CornWithRunningDir() string {
	return filepath.Join(c.RunningDir, c.CornDirName)
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
const profileAsUserProfile = "${userprofile}"

func (c MainConfig) IsUserProfile() bool {
	return (c.ProfileDir == profileAsUserProfile)
}

const profileAsCurrentDir = "${currentdir}"

const platId = "${platid}"

var defaultCoreConfig = MainConfig{
	CurrentDir:     execDirWithoutLink,
	RunningDir:     execDirWithoutLink,
	RuntimeDirName: "runtimes",
	CornDirName:    "corns",
	ProfileDir:     profileAsUserProfile,
}

func loadConfigRegular(config string, value interface{}, altBases ...string) error {
	errFmt := fs.PathError{
		Op:   "loadConfig",
		Path: config,
	}

	if len(config) == 0 {
		errFmt.Err = errors.New("could not load config by empty name")
		return &errFmt
	}

	basePath := ""
	if len(altBases) > 0 {
		basePath = altBases[0]
	} else {
		basePath, _ = os.Executable()
		basePath, _ = filepath.EvalSymlinks(basePath)
		basePath = filepath.Dir(basePath)
	}
	if len(basePath) == 0 {
		return errors.New("could not load workdir")
	}

	tomlFilename := path.Join(basePath, config+".toml")
	stat, _ := os.Stat(tomlFilename)
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

func LoadCoreConfig(altBases ...string) MainConfig {
	basePath, _ := os.Executable()
	basePath, _ = filepath.EvalSymlinks(basePath)
	if len(basePath) == 0 && len(altBases) == 0 {
		LogCLI(zerolog.FatalLevel).Println("Could not load workdir,use altBase instead")
	}
	basePath = filepath.Dir(basePath)

	if len(altBases) > 0 {
		basePath = altBases[0]
	}

	var result = defaultCoreConfig
	err := loadConfigRegular("core", &result, basePath)
	if err != nil {
		LogCLI(zerolog.WarnLevel).Println(err.Error())
	}

	if result.CurrentDir == execDirWithoutLink {
		result.CurrentDir = basePath
	}

	if result.RunningDir == execDirWithoutLink {
		result.RunningDir = basePath
	}

	if result.ProfileDir == profileAsCurrentDir {
		result.ProfileDir = filepath.Join(basePath, "profile")
	}

	if result.ProfileDir != profileAsUserProfile &&
		filepath.IsLocal(result.ProfileDir) {
		result.ProfileDir = filepath.Join(basePath, result.ProfileDir)
	}

	result.RunningDir = platId
	if path, ok := result.RunningDirByPlatfrom[platId]; ok {
		result.RunningDir = path + result.RunningDir
	}
	if path, ok := result.RunningDirByPlatfrom[utils.OS()]; ok {
		result.RunningDir = path
	}
	if path, ok := result.RunningDirByPlatfrom[utils.Platform()]; ok {
		result.RunningDir = path
	}
	result.RunningDir = strings.ReplaceAll(result.RunningDir,
		platId, utils.Platform())

	if !filepath.IsAbs(result.RunningDir) {
		result.RunningDir = filepath.Join(basePath, result.RunningDir)
	}

	return result
}
