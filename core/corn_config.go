package core

import (
	"cryphtron/internal/utils"
	"errors"
	"path/filepath"
)

const CornConfigExt = ".toml.corn"

func LoadCornConfig(tomlPath string, base envConfig) (CornsEnvConfig, error) {
	result := CornsEnvConfig{}
	if base.coreConfig == nil {
		return result, errors.New("could not loading the env wihout core config")
	}

	loadPath := tomlPath
	pth, file := filepath.Split(tomlPath)
	if len(pth) == 0 {
		loadPath = filepath.Join(base.coreConfig.CurrentWorkDir, file)
	}
	name := file[0 : len(file)-len(CornConfigExt)]
	result.envConfig = base
	result.ID = CornsIdentifier + "_" + name
	result.envName = name
	result.initCornVars()

	err := utils.LoadTomlFile(loadPath, &result)
	if err != nil {
		return result, err
	}

	if result.DirName == "" {
		result.DirName = name
	}

	for idx := range result.BootstrapExec {
		result.BootstrapExec[idx].Top = &result.ValueScope
	}

	for idx := range result.ConfigExec {
		result.ConfigExec[idx].Top = &result.ValueScope
	}

	if exec, ok := result.ExecByPlats[utils.OS()]; ok {
		result.Exec = setExec(result, exec)
	}

	if exec, ok := result.ExecByPlats[utils.Platform()]; ok {
		result.Exec = setExec(result, exec)
	}

	result.Exec.Top = &result.ValueScope

	return result, nil
}
