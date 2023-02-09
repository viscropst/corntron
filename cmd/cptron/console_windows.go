//go:build windows && !appengine
// +build windows,!appengine

package cptron

import (
	"os"

	"github.com/mattn/go-isatty"
)

func IsInTerminal() bool {
	fd := os.Stdin.Fd()
	isTerm := isatty.IsTerminal(fd)
	isCygwin := isatty.IsCygwinTerminal(fd)
	return isTerm || isCygwin
}
