package core

import (
	"errors"
	"path/filepath"
)

const CornConfigExt = ".toml.corn"

func LoadCornConfig(tomlPath string, base envConfig) (CornsEnvConfig, error) {
	result := CornsEnvConfig{}
	if base.coreConfig == nil {
		return result, errors.New("could not load the env without core config")
	}

	pth, file := filepath.Split(tomlPath)
	if len(pth) == 0 {
		return LoadCornEnv(file, base, base.coreConfig.CurrentWorkDir)
	}
	return LoadCornEnv(file, base, pth)
}
