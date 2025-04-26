//go:build !windows && !appengine

package utils

import "syscall"

func GetNewProcGroupAttr(isBackground, newGroup bool) *syscall.SysProcAttr {
	return nil
}
