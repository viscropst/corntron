//go:build !windows && !appengine

package core

import (
	"os/exec"
)

func (c Command) setSysProcAttr(cmd *exec.Cmd) {
}
