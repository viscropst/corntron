//go:build linux && !windows && !appengine
// +build linux,!windows,!appengine

package utils

import "syscall"

func GetNewProcGroupAttr(isBackground, newGroup bool) *syscall.SysProcAttr {
	result := syscall.SysProcAttr{}
	if newGroup {
		result.Setpgid = true
		result.Foreground = true
	}
	return &result
}
