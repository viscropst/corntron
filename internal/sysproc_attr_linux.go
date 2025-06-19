//go:build linux && !windows && !appengine
// +build linux,!windows,!appengine

package internal

import (
	"syscall"
)

func GetNewProcGroupAttr(isBackground, newGroup bool) *syscall.SysProcAttr {
	result := syscall.SysProcAttr{}
	result.Setpgid = newGroup
	return &result
}
