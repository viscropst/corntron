//go:build windows
// +build windows

package utils

import (
	"syscall"
)

const CREATE_NO_WINDOW = 0x08000000

func GetNewProcGroupAttr(isBackground, newGroup bool) *syscall.SysProcAttr {
	result := syscall.SysProcAttr{}
	var flag uint32 = syscall.CREATE_UNICODE_ENVIRONMENT
	if newGroup || IsInTerminal() {
		flag = flag | syscall.CREATE_NEW_PROCESS_GROUP
	}
	if isBackground {
		flag = flag | CREATE_NO_WINDOW
		result.HideWindow = true
	}
	result.CreationFlags = flag
	return &result
}
