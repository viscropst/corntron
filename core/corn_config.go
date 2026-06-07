package core

import (
	"corntron/internal"
	"io"
	"path/filepath"
)

const CornConfigExt = ".toml.corn"

func LoadCornConfigFile(tomlPath string, base envConfig) (CornsEnvConfig, error) {
	result := CornsEnvConfig{}
	if base.coreConfig == nil {
		return result, internal.Error("could not load the env without core config")
	}
	loadPath := filepath.Join(
		base.coreConfig.CornDir(), result.envDirname)
	pth, file := filepath.Split(tomlPath)
	if len(pth) == 0 {
		loadPath = base.coreConfig.CurrentWorkDir
	} else if len(pth) > 0 {
		loadPath = pth
	}
	err := loadConfigRegular(file, &result, loadPath)
	if err != nil {
		return result, err
	}
	return InitCornEnv(&result, file, base)
}

func LoadCornConfigReader(name string, reader io.Reader, base envConfig) (CornsEnvConfig, error) {
	result := CornsEnvConfig{}
	if len(name) == 0 {
		return result, internal.Error("could not loading the env without name")
	}
	if reader == nil {
		return result, internal.Error("could not loading the env without reader")
	}
	if base.coreConfig == nil {
		return result, internal.Error("could not loading the env without core config")
	}
	err := internal.LoadTomlReader(reader, &result)
	if err != nil {
		return result, err
	}
	return InitCornEnv(&result, name, base)
}
