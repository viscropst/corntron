//go:build windows
// +build windows

package core

import (
	"os/exec"
	"syscall"
)

const CREATE_NO_WINDOW = 0x08000000

func (c Command) setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	var flag uint32 = syscall.CREATE_UNICODE_ENVIRONMENT |
		syscall.CREATE_NEW_PROCESS_GROUP
	if c.IsBackground {
		flag = flag | CREATE_NO_WINDOW
		cmd.SysProcAttr.HideWindow = true
	}
	cmd.SysProcAttr.CreationFlags = flag
}
