//go:build windows && !appengine
// +build windows,!appengine

package utils

import (
	"os"

	"github.com/inconshreveable/mousetrap"
	"github.com/mattn/go-isatty"
)

func IsInTerminal() bool {
	fd := os.Stdin.Fd()
	if mousetrap.StartedByExplorer() {
		return false
	}
	isTerm := isatty.IsTerminal(fd)
	isCygwin := isatty.IsCygwinTerminal(fd)
	return isTerm || isCygwin
}
