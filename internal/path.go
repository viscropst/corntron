package internal

import (
	"errors"
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
	return filepath.Join(src)
}

func GetExecPath(execStr string, pathList string) (string, error) {
	_ = os.Setenv("PATH", pathList)
	path, err := exec.LookPath(execStr)
	_ = os.Unsetenv("PATH")
	if err != nil {
		errBuilder := strings.Builder{}
		errBuilder.WriteString("exec argument invalid: the command could not found")
		return execStr, errors.New(errBuilder.String())
	} else {
		return path, nil
	}
}
