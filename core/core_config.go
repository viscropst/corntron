package core

import (
	"corntron/internal"
	"corntron/internal/log"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
)

type MainConfig struct {
	CurrentDir         string            `toml:"base_dir,omitempty"`
	CurrentPlatformDir string            `toml:"-"`
	CurrentWorkDir     string            `toml:"-"`
	RuntimeDirName     string            `toml:"runtime_dirname,omitempty"`
	CornDirName        string            `toml:"corn_dirname,omitempty"`
	MirrorType         string            `toml:"mirror_type,omitempty"`
	MirrorTypes        []string          `toml:"mirror_types,omitempty"`
	mirrorTypes        []MirrorType      `toml:"-"`
	WithCorn           bool              `toml:"with_corn"`
	PlatformDirs       map[string]string `toml:"platform_dir,omitempty"`
	ProfileDir         string            `toml:"profile_dir,omitempty"`
}

func (c MainConfig) RuntimeDir() string {
	return filepath.Join(c.CurrentDir, c.RuntimeDirName)
}

func (c MainConfig) CornDir() string {
	return filepath.Join(c.CurrentDir, c.CornDirName)
}

func (c MainConfig) RtWithRunningDir() string {
	return filepath.Join(c.CurrentPlatformDir, c.RuntimeDirName)
}

func (c MainConfig) CornWithRunningDir() string {
	return filepath.Join(c.CurrentPlatformDir, c.CornDirName)
}

func (c MainConfig) FsWalk(walkFunc filepath.WalkFunc, DirNames ...string) error {
	rootDir := c.CurrentDir
	if len(DirNames) > 0 {
		tmp := filepath.Join(DirNames...)
		rootDir = filepath.Join(rootDir, tmp)
	}
	return filepath.Walk(rootDir, walkFunc)
}

func (c MainConfig) FsWalkDir(walkFunc fs.WalkDirFunc, DirNames ...string) error {
	rootDir := filepath.Base(c.CurrentDir)
	if len(DirNames) > 0 {
		tmp := filepath.Join(DirNames...)
		rootDir = filepath.Join(rootDir, tmp)
	}
	return filepath.WalkDir(rootDir, walkFunc)
}

func (c MainConfig) UnwrapMirrorType() MirrorType {
	return MirrorType(c.MirrorType).ConvertWithTypes(c.mirrorTypes...)
}

const execDirWithoutLink = "${dp0}"
const profileAsUserProfile = "${userprofile}"

func (c MainConfig) IsUserProfile() bool {
	return (c.ProfileDir == profileAsUserProfile)
}

const profileAsCurrentDir = "${currentdir}"

const platID = "${platid}"

var defaultCoreConfig = MainConfig{
	CurrentDir:         execDirWithoutLink,
	CurrentPlatformDir: platID,
	RuntimeDirName:     "runtimes",
	CornDirName:        "corns",
	ProfileDir:         profileAsUserProfile,
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
		basePath = internal.GetSelfPath()
	}
	if len(basePath) == 0 {
		return errors.New("could not load workdir")
	}
	tomlFilename := filepath.Join(basePath, config+".toml")
	if strings.HasSuffix(config, CornConfigExt) {
		tomlFilename = filepath.Join(basePath, config)
	}
	pathErr := internal.LoadTomlFile(tomlFilename, value)
	if pathErr != nil {
		errFmt.Err = pathErr.Err
		return &errFmt
	}
	return nil
}

func LoadCoreConfig(altBases ...string) MainConfig {
	basePath := internal.GetSelfPath()
	if len(basePath) == 0 && len(altBases) == 0 {
		LogCLI(log.FatalLevel).Println("Could not load workdir,use altBase instead")
	}
	basePath = filepath.Dir(basePath)

	if len(altBases) > 0 {
		basePath = altBases[0]
	}

	var result = defaultCoreConfig
	err := loadConfigRegular("core", &result, basePath)
	if err != nil {
		LogCLI(log.WarnLevel).Println(err.Error())
	}

	if result.CurrentDir == execDirWithoutLink {
		result.CurrentDir = basePath
	}

	result.CurrentWorkDir = internal.GetWorkDir(basePath)

	result.ProfileDir = prepareProfileDir(
		result, basePath)

	result.CurrentPlatformDir = prepareRunningDir(
		result, basePath)

	result.mirrorTypes = prepareMirrorTypes(result)

	return result
}

func prepareProfileDir(src MainConfig, basePath string) string {
	result := src.ProfileDir
	if result == profileAsUserProfile {
		result = internal.GetProfileDir()
	}
	if result == "" {
		result = profileAsCurrentDir
	}
	if result == profileAsCurrentDir {
		return filepath.Join(basePath, "profile")
	}
	if !filepath.IsAbs(result) {
		return filepath.Join(basePath, result)
	}
	return result
}

func prepareRunningDir(src MainConfig, basePath string) string {
	result := src.CurrentPlatformDir
	if result == execDirWithoutLink {
		return basePath
	}
	if path, ok := src.PlatformDirs[platID]; ok {
		result = path + result
	}
	if path, ok := src.PlatformDirs[internal.OS()]; ok {
		result = path
	}
	if path, ok := src.PlatformDirs[internal.Platform()]; ok {
		result = path
	}
	result = strings.ReplaceAll(result,
		platID, internal.Platform())

	if !filepath.IsAbs(result) {
		result = filepath.Join(basePath, result)
	}

	return result
}

func prepareMirrorTypes(src MainConfig) []MirrorType {
	return MirrorTypesFromSlice(src.MirrorTypes)
}
