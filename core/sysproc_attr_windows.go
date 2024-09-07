//go:build windows
// +build windows

package core

import (
	"os/exec"
	"syscall"
)

func (c Command) setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_UNICODE_ENVIRONMENT |
			syscall.CREATE_NEW_PROCESS_GROUP,
		HideWindow: true,
	}
}
