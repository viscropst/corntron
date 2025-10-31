package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const PathSeparator = "/"
const PathListSeparator = ";"
const OSPathSeparator = string(os.PathSeparator)
const OSPathListSeparator = string(os.PathListSeparator)

func NormalizePath(src string) string {
	result := filepath.Join(src)
	return strings.TrimSuffix(result, OSPathSeparator)
}

func GetExecPath(execStr string, pathList string) (string, error) {
	_ = os.Setenv("PATH", pathList)
	path, err := exec.LookPath(execStr)
	_ = os.Unsetenv("PATH")
	if err != nil {
		errBuilder := strings.Builder{}
		errBuilder.WriteString("exec argument invalid: the command could not found")
		return execStr, Error(errBuilder.String())
	} else {
		return path, nil
	}
}

func GetSelfPath() string {
	basePath, _ := os.Executable()
	basePath, _ = filepath.EvalSymlinks(basePath)
	return basePath
}

func GetSelfDir() string {
	return filepath.Dir(GetSelfPath())
}

func GetWorkDir(alt ...string) string {
	result, err := os.Getwd()
	if err != nil && len(alt) > 0 {
		return alt[0]
	}
	return result
}

func GetProfileDir() string {
	result, _ := os.UserHomeDir()
	return result
}
