package core

import (
	"corntron/internal"
	"path/filepath"
)

const CornConfigExt = ".toml.corn"

func LoadCornConfig(tomlPath string, base envConfig) (CornsEnvConfig, error) {
	result := CornsEnvConfig{}
	if base.coreConfig == nil {
		return result, internal.Error("could not load the env without core config")
	}

	pth, file := filepath.Split(tomlPath)
	if len(pth) == 0 {
		return LoadCornEnv(file, base, base.coreConfig.CurrentWorkDir)
	}
	return LoadCornEnv(file, base, pth)
}
